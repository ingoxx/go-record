package main

import (
	"encoding/json"
	"fmt"
	"github.com/ingoxx/go-record/http/wx/redis"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Group 一个群聊包含多个客户端连接 + 消息历史
type Group struct {
	Clients  map[*websocket.Conn]bool
	Messages []Message
	Lock     sync.Mutex
}

type Message struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
	Time    string `json:"time"`

	Type      string `json:"type"`       // normal / count
	UserCount int    `json:"user_count"` // 当前群人数
}

type Resp struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	w    http.ResponseWriter
}

func (r Resp) message(rd Resp) ([]byte, error) {
	b, err := json.Marshal(rd)
	if err != nil {
		return b, err
	}
	return b, nil
}

func (r Resp) h(rd Resp) {
	message, err := r.message(rd)
	if err != nil {
		log.Printf("[ERROR] fail to respone, error '%v'", err)
		return
	}
	if _, err := r.w.Write(message); err != nil {
		log.Printf("[ERROR] fail to respone, error '%v'", err)
		return
	}
}

var (
	// 所有群: groupID => Group
	groups   = make(map[string]*Group)
	groupsMu sync.Mutex

	// 全局广播
	broadcast = make(chan Message)

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	log.Println("version: v1.0.1")

	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/get-online", handleOnline)
	http.HandleFunc("/add-square", handleAddBasketballSquare)
	http.HandleFunc("/show-square", handleShowBasketballSquare)

	// 启动广播处理器
	go handleBroadcast()

	log.Println("Server started on :11806")
	log.Fatal(http.ListenAndServe(":11806", nil))
}

func handleShowBasketballSquare(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodGet {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}
	lng := r.FormValue("lng")
	lat := r.FormValue("lat")
	if lng == "" || lat == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	ol, err := redis.NewRM().GetAllData()
	if err != nil {
		rp.h(Resp{
			Msg:  err.Error(),
			Code: 1003,
			Data: ol,
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: ol,
	})
}

func handleAddBasketballSquare(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodPost {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}

}

func handleOnline(w http.ResponseWriter, r *http.Request) {
	var rp = Resp{w: w}
	if r.Method != http.MethodGet {
		rp.h(Resp{
			Msg:  "invalid request",
			Code: 1001,
			Data: "0",
		})
		return
	}

	gid := r.FormValue("gid")
	if gid == "" {
		rp.h(Resp{
			Msg:  "invalid parameter",
			Code: 1002,
			Data: "0",
		})
		return
	}

	ol, err := redis.NewRM().Get(gid)
	if err != nil {
		rp.h(Resp{
			Msg:  fmt.Sprintf("'%s' not found", gid),
			Code: 1003,
			Data: "0",
		})
		return
	}

	rp.h(Resp{
		Msg:  "ok",
		Code: 1000,
		Data: ol,
	})

}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket 升级失败:", err)
		return
	}
	defer ws.Close()

	// 先读取第一条消息，拿到 groupID
	_, msgBytes, err := ws.ReadMessage()
	if err != nil {
		log.Println("初始化连接读取失败:", err)
		return
	}

	var initMsg Message
	if err := json.Unmarshal(msgBytes, &initMsg); err != nil {
		log.Println("解析初始化消息失败:", err)
		return
	}
	groupID := initMsg.GroupID

	// 把连接放入对应群
	groupsMu.Lock()
	if _, ok := groups[groupID]; !ok {
		groups[groupID] = &Group{
			Clients:  make(map[*websocket.Conn]bool),
			Messages: []Message{},
		}
	}
	group := groups[groupID]
	groupsMu.Unlock()

	// 给群组添加ws
	group.Lock.Lock()
	group.Clients[ws] = true
	userCount := len(group.Clients)
	if err := redis.NewRM().Set(groupID, userCount); err != nil {
		log.Printf("[ERROR] 写入redis失败, 错误信息：%v", err)
	}
	group.Lock.Unlock()

	log.Printf("用户 %s 加入群 %s，当前人数: %d", initMsg.UserID, groupID, userCount)

	// 广播新的群人数
	broadcast <- Message{
		GroupID:   groupID,
		Type:      "count",
		UserCount: userCount,
	}

	// 先把历史消息发给新连接（可选）
	group.Lock.Lock()
	for _, oldMsg := range group.Messages {
		if err := ws.WriteJSON(oldMsg); err != nil {
			log.Println("发送历史消息失败:", err)
		}
	}
	group.Lock.Unlock()

	// 持续读取新消息
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("客户端 %s, 当前人数: %d,  断开: %v", initMsg.UserID, initMsg.UserCount, err)

			// 从群里删除这个连接
			group.Lock.Lock()
			delete(group.Clients, ws)
			userCount = len(group.Clients)
			if err := redis.NewRM().Set(msg.GroupID, userCount); err != nil {
				log.Printf("[ERROR] 写入redis失败, 错误信息：%v", err)
			}
			group.Lock.Unlock()

			// 广播新的群人数
			broadcast <- Message{
				GroupID:   groupID,
				Type:      "count",
				UserCount: userCount,
			}

			break
		}

		// 普通消息
		msg.Type = "normal"

		// 保存到历史
		group.Lock.Lock()
		group.Messages = append(group.Messages, msg)
		group.Lock.Unlock()

		// 广播消息
		broadcast <- msg
	}
}

func handleBroadcast() {
	for {
		msg := <-broadcast
		groupID := msg.GroupID

		groupsMu.Lock()
		group, ok := groups[groupID]
		if err := redis.NewRM().Set(msg.GroupID, msg.UserCount); err != nil {
			log.Println("[ERROR] fail to save user count.")
		}
		groupsMu.Unlock()
		if !ok {
			continue
		}

		group.Lock.Lock()
		for client := range group.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("广播失败，删除连接:", err)
				client.Close()
				delete(group.Clients, client)
			}
		}
		group.Lock.Unlock()
	}
}
