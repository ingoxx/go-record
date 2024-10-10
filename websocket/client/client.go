package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Resp struct {
	Result string `json:"result"`
}

func main() {
	socketUrl := "wss://api.anythingai.online/bd?openid=ogR3E62jXXJMbVcImRqMA1gTSegM"
	dial := websocket.Dialer{
		HandshakeTimeout: 300 * time.Second,
	}

	conn, _, err := dial.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()

	var data = map[string]interface{}{
		"context": []map[string]string{
			{
				"role":      "user",
				"assistant": "golang的map,slice的初始化？",
			},
		},
	}
	b, _ := json.Marshal(&data)

	err = conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		log.Println("Error during writing to websocket:", err)
		return
	}

	var res Resp
	for {

		//接收
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		json.Unmarshal(msg, &res)

		log.Printf("Received: %s\n", res)
	}

}
