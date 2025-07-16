package main

import (
	"encoding/json"
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
	http.HandleFunc("/ws", handleConnections)

	// 启动广播处理器
	go handleBroadcast()

	log.Println("Server started on :11806")
	log.Fatal(http.ListenAndServe(":11806", nil))
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
	group.Lock.Unlock()

	log.Printf("用户 %s 加入群 %s", initMsg.UserID, groupID)

	// 先把历史消息发给新连接
	group.Lock.Lock()
	for _, oldMsg := range group.Messages {
		if err := ws.WriteJSON(oldMsg); err != nil {
			log.Println("发送历史消息失败:", err)
		}
	}
	group.Lock.Unlock()

	// 广播用户加入通知（可选）
	broadcast <- initMsg

	// 持续读取新消息
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("客户端 %s 断开: %v", initMsg.UserID, err)
			group.Lock.Lock()
			delete(group.Clients, ws)
			group.Lock.Unlock()
			break
		}

		// 把消息保存到历史
		group.Lock.Lock()
		group.Messages = append(group.Messages, msg)
		group.Lock.Unlock()

		// 放入广播队列
		broadcast <- msg
	}
}

func handleBroadcast() {
	for {
		msg := <-broadcast
		groupID := msg.GroupID

		groupsMu.Lock()
		group, ok := groups[groupID]
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
