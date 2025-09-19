package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// 消息结构
type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
	Time    string `json:"time"`
	Role    string `json:"role"` // "publisher" 或 "user"
}

// 用户结构
type User struct {
	ID   string
	Role string
	Conn *websocket.Conn
}

var (
	mu             sync.RWMutex
	clients        = make(map[string]*User)           // 在线用户
	roomHistory    = make([]Message, 0)               // 所有消息统一存房间历史
	publisherUsers = make(map[string]map[string]bool) // 发布者 -> 曾聊过的用户
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// 异步发送消息给指定用户
func sendMessage(to string, msg Message) {
	mu.RLock()
	user, ok := clients[to]
	mu.RUnlock()
	if !ok {
		// 接收方不在线，消息已存历史，等待上线
		return
	}

	go func() {
		if err := user.Conn.WriteJSON(msg); err != nil {
			fmt.Println("发送失败:", err)
		}
	}()
}

// WebSocket 处理
func handleWS(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user")
	role := r.URL.Query().Get("role")
	if userID == "" || role == "" {
		http.Error(w, "missing user or role", http.StatusBadRequest)
		return
	}

	conn, _ := upgrader.Upgrade(w, r, nil)

	mu.Lock()
	clients[userID] = &User{ID: userID, Role: role, Conn: conn}

	// 上线时推送自己参与的历史消息
	for _, msg := range roomHistory {
		if role == "publisher" || msg.From == userID || msg.To == userID {
			go func(m Message) { conn.WriteJSON(m) }(msg)
		}
	}

	mu.Unlock()
	fmt.Println(userID, "connected as", role)

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		msg.Time = time.Now().Format("2006-01-02 15:04:05")
		msg.Role = role

		mu.Lock()
		// 存入房间历史
		roomHistory = append(roomHistory, msg)

		// 记录发布者的用户列表
		if role == "publisher" {
			if publisherUsers[msg.From] == nil {
				publisherUsers[msg.From] = make(map[string]bool)
			}
			publisherUsers[msg.From][msg.To] = true
		} else if toUser, ok := clients[msg.To]; ok && toUser.Role == "publisher" {
			if publisherUsers[msg.To] == nil {
				publisherUsers[msg.To] = make(map[string]bool)
			}
			publisherUsers[msg.To][msg.From] = true
		}
		mu.Unlock()

		// 发送消息给目标用户
		sendMessage(msg.To, msg)
	}

	mu.Lock()
	delete(clients, userID)
	mu.Unlock()
	fmt.Println(userID, "disconnected")
}

// 获取发布者的用户列表
func getPublisherUsers(w http.ResponseWriter, r *http.Request) {
	publisherID := r.URL.Query().Get("publisher")
	mu.RLock()
	users := []string{}
	for uid := range publisherUsers[publisherID] {
		users = append(users, uid)
	}
	mu.RUnlock()
	json.NewEncoder(w).Encode(users)
}

func main() {
	http.HandleFunc("/ws", handleWS)
	http.HandleFunc("/publisher-users", getPublisherUsers)
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
