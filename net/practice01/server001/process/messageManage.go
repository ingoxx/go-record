package process

import (
	"fmt"

	"github.com/ingoxx/Golang-practise/net/practice01/model"
)

var (
	Mm *MessageManage
)

func init() {
	Mm = &MessageManage{
		MessageList: make(map[int][]model.ChatUserToUserMessage, 1024),
	}
}

type MessageManage struct {
	MessageList map[int][]model.ChatUserToUserMessage
}

// 保存离线用户的消息
func (m *MessageManage) AddMessage(sm model.ChatUserToUserMessage) {
	m.MessageList[sm.UserId] = append(m.MessageList[sm.UserId], sm)
	fmt.Println("chp1 MessageList =", m.MessageList)
}
