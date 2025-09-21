package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// 内存存储结构
type MemoryStorage struct {
	Conversations map[uint]*Conversation
	Messages      map[uint][]*Message // conversationID -> messages
	Users         map[uint]*User
	Bounties      map[uint]*Bounty
	mutex         sync.RWMutex
}

// 用户模型
type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// 悬赏任务模型
type Bounty struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PublisherID uint      `json:"publisher_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// 会话模型
type Conversation struct {
	ID        uint      `json:"id"`
	BountyID  uint      `json:"bounty_id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 消息模型
type Message struct {
	ID             uint      `json:"id"`
	ConversationID uint      `json:"conversation_id"`
	SenderID       uint      `json:"sender_id"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

// WebSocket相关结构
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	UserID         uint
	ConversationID uint
	Conn           *websocket.Conn
	Send           chan []byte
}

type Hub struct {
	Conversations map[uint]map[*Client]bool
	Register      chan *Client
	Unregister    chan *Client
	Broadcast     chan MessageEvent
	Mutex         sync.RWMutex
}

type MessageEvent struct {
	ConversationID uint   `json:"conversation_id"`
	SenderID       uint   `json:"sender_id"`
	Content        string `json:"content"`
}

// 全局变量
var (
	hub          *Hub
	storage      *MemoryStorage
	nextConvID   uint = 1
	nextMsgID    uint = 1
	nextUserID   uint = 1
	nextBountyID uint = 1
)

// 初始化函数
func init() {
	// 初始化内存存储
	storage = &MemoryStorage{
		Conversations: make(map[uint]*Conversation),
		Messages:      make(map[uint][]*Message),
		Users:         make(map[uint]*User),
		Bounties:      make(map[uint]*Bounty),
	}

	// 初始化一些测试数据
	initTestData()

	// 初始化WebSocket Hub
	hub = &Hub{
		Conversations: make(map[uint]map[*Client]bool),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Broadcast:     make(chan MessageEvent),
	}
	go hub.Run()
}

// 初始化测试数据
func initTestData() {
	// 创建测试用户
	user1 := &User{ID: nextUserID, Name: "发布者"}
	nextUserID++
	storage.Users[user1.ID] = user1

	user2 := &User{ID: nextUserID, Name: "接单用户1"}
	nextUserID++
	storage.Users[user2.ID] = user2

	user3 := &User{ID: nextUserID, Name: "接单用户2"}
	nextUserID++
	storage.Users[user3.ID] = user3

	// 创建测试悬赏任务
	bounty := &Bounty{
		ID:          nextBountyID,
		Title:       "测试悬赏任务",
		Description: "这是一个测试用的悬赏任务",
		PublisherID: user1.ID,
		CreatedAt:   time.Now(),
	}
	nextBountyID++
	storage.Bounties[bounty.ID] = bounty
}

// Hub运行函数
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			if h.Conversations[client.ConversationID] == nil {
				h.Conversations[client.ConversationID] = make(map[*Client]bool)
			}
			h.Conversations[client.ConversationID][client] = true
			h.Mutex.Unlock()

		case client := <-h.Unregister:
			h.Mutex.Lock()
			if _, ok := h.Conversations[client.ConversationID][client]; ok {
				delete(h.Conversations[client.ConversationID], client)
				close(client.Send)
			}
			h.Mutex.Unlock()

		case event := <-h.Broadcast:
			h.Mutex.RLock()
			clients := h.Conversations[event.ConversationID]
			for client := range clients {
				msg, _ := json.Marshal(event)
				select {
				case client.Send <- msg:
				default:
					close(client.Send)
					delete(clients, client)
				}
			}
			h.Mutex.RUnlock()
		}
	}
}

func main() {
	r := mux.NewRouter()

	// WebSocket路由
	r.HandleFunc("/ws", handleWebSocket)

	// API路由
	r.HandleFunc("/conversations/{bountyId}", getConversations).Methods("GET")
	r.HandleFunc("/messages/{conversationId}", getMessages).Methods("GET")
	r.HandleFunc("/conversation", createConversation).Methods("POST")
	r.HandleFunc("/bounties", getBounties).Methods("GET")
	r.HandleFunc("/users/{userId}", getUser).Methods("GET")

	// 静态文件服务（可选，用于测试）
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// WebSocket处理
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// 从查询参数获取用户ID和会话ID
	userID := r.URL.Query().Get("user_id")
	conversationID := r.URL.Query().Get("conversation_id")

	if userID == "" || conversationID == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("参数错误"))
		conn.Close()
		return
	}

	var uid, cid uint
	fmt.Sscanf(userID, "%d", &uid)
	fmt.Sscanf(conversationID, "%d", &cid)

	client := &Client{
		UserID:         uid,
		ConversationID: cid,
		Conn:           conn,
		Send:           make(chan []byte, 256),
	}

	hub.Register <- client

	// 启动读写goroutine
	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msgEvent MessageEvent
		if err := json.Unmarshal(message, &msgEvent); err != nil {
			continue
		}

		// 保存消息到内存
		storage.mutex.Lock()
		msg := &Message{
			ID:             nextMsgID,
			ConversationID: c.ConversationID,
			SenderID:       c.UserID,
			Content:        msgEvent.Content,
			CreatedAt:      time.Now(),
		}
		nextMsgID++

		storage.Messages[c.ConversationID] = append(storage.Messages[c.ConversationID], msg)

		// 更新会话时间
		if conv, exists := storage.Conversations[c.ConversationID]; exists {
			conv.UpdatedAt = time.Now()
		}
		storage.mutex.Unlock()

		// 广播消息
		hub.Broadcast <- MessageEvent{
			ConversationID: c.ConversationID,
			SenderID:       c.UserID,
			Content:        msgEvent.Content,
		}
	}
}

func (c *Client) writePump() {
	defer c.Conn.Close()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

// HTTP API处理
func getConversations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bountyID := vars["bountyId"]

	if bountyID == "" {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}

	var bid uint
	fmt.Sscanf(bountyID, "%d", &bid)

	storage.mutex.RLock()
	defer storage.mutex.RUnlock()

	// 查找特定悬赏的所有会话
	var conversations []*Conversation
	for _, conv := range storage.Conversations {
		if conv.BountyID == bid {
			conversations = append(conversations, conv)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationID := vars["conversationId"]

	if conversationID == "" {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}

	var cid uint
	fmt.Sscanf(conversationID, "%d", &cid)

	storage.mutex.RLock()
	messages, exists := storage.Messages[cid]
	storage.mutex.RUnlock()

	if !exists {
		http.Error(w, "会话不存在", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func createConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		BountyID uint `json:"bounty_id"`
		UserID   uint `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}

	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	// 检查是否已存在会话
	for _, conv := range storage.Conversations {
		if conv.BountyID == data.BountyID && conv.UserID == data.UserID {
			http.Error(w, "会话已存在", http.StatusConflict)
			return
		}
	}

	// 创建新会话
	conversation := &Conversation{
		ID:        nextConvID,
		BountyID:  data.BountyID,
		UserID:    data.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	nextConvID++

	storage.Conversations[conversation.ID] = conversation
	storage.Messages[conversation.ID] = []*Message{}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}

func getBounties(w http.ResponseWriter, r *http.Request) {
	storage.mutex.RLock()
	defer storage.mutex.RUnlock()

	var bounties []*Bounty
	for _, bounty := range storage.Bounties {
		bounties = append(bounties, bounty)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bounties)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	if userID == "" {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}

	var uid uint
	fmt.Sscanf(userID, "%d", &uid)

	storage.mutex.RLock()
	user, exists := storage.Users[uid]
	storage.mutex.RUnlock()

	if !exists {
		http.Error(w, "用户不存在", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
