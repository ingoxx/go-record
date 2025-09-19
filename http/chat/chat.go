package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type PublishMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

var (
	clients     = make(map[string]*websocket.Conn)  // 在线用户
	chatHistory = make(map[string][]PublishMessage) // key: "a|b"
	mu          sync.RWMutex
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// 生成唯一key
func chatKey(a, b string) string {
	if a < b {
		return a + "|" + b
	}
	return b + "|" + a
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user")
	if userID == "" {
		http.Error(w, "missing user", http.StatusBadRequest)
		return
	}
	conn, _ := upgrader.Upgrade(w, r, nil)

	mu.Lock()
	clients[userID] = conn

	// 🔥 登录时推送该用户参与过的所有历史消息
	for key, history := range chatHistory {
		if strings.Contains(key, userID) {
			for _, m := range history {
				conn.WriteJSON(m)
			}
		}
	}
	mu.Unlock()

	// 监听消息
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var msg PublishMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		msg.Time = time.Now().Format("2006-01-02 15:04:05")

		key := chatKey(msg.From, msg.To)

		mu.Lock()
		chatHistory[key] = append(chatHistory[key], msg)

		// 转发给接收方
		if toConn, ok := clients[msg.To]; ok {
			toConn.WriteJSON(msg)
		}
		mu.Unlock()
	}

	mu.Lock()
	delete(clients, userID)
	mu.Unlock()
	fmt.Println(userID, "disconnected")
}

func main() {
	http.HandleFunc("/ws", handleWS)
	fmt.Println("Server on :8080")
	http.ListenAndServe(":8080", nil)
}
