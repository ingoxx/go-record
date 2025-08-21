package process

import (
	"fmt"

	UserMessage "github.com/ingoxx/Golang-practise/net/practice01/model"
)

var (
	Cr *ChatRecord
)

func init() {
	Cr = &ChatRecord{
		MessageRecord: []UserMessage.ChatUserToUserMessage{},
	}
}

type ChatRecord struct {
	MessageRecord []UserMessage.ChatUserToUserMessage
}

func (c *ChatRecord) ShowChatRecord() {
	for _, v1 := range c.MessageRecord {
		fmt.Printf("\t发送者:%d, 内容:%v\n", v1.Sender, v1.Content)
	}
}
