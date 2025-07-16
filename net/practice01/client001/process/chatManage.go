package process

import (
	"encoding/json"
	"fmt"

	UserMessage "github.com/ingoxx/Golang-practise/net/practice01/model"
)

// 接收群消息
func RecvGroupMsg(mes *UserMessage.Message) {

	cr := UserMessage.ChatMessageResult{}

	err := json.Unmarshal([]byte(mes.Data), &cr)
	if err != nil {
		fmt.Println("群消息反序列化失败,err=", err)
		return
	}

	fmt.Printf("用户:%d,消息:%v\n", cr.UserId, cr.Content)
	fmt.Println()
}

// 接收一对一消息
func RecvUserMsg(mes *UserMessage.Message) {

	cr := UserMessage.ChatUserToUserMessageResult{}

	err := json.Unmarshal([]byte(mes.Data), &cr)
	if err != nil {
		fmt.Println("群消息反序列化失败,err=", err)
		return
	}

	fmt.Printf("用户:%d,私聊的消息:%v\n", cr.UserId, cr.Content)
	fmt.Println()
}

// 接收服务端推送的离线聊天记录
func RecvChatRecord(mes *UserMessage.Message) {

	cr := UserMessage.ChatRrcordMessageResult{}

	err := json.Unmarshal([]byte(mes.Data), &cr)
	if err != nil {
		fmt.Println("群消息反序列化失败,err=", err)
		return
	}

	//把服务端推送的聊天记录保存在本地
	Cr.MessageRecord = cr.ChatRecord
}
