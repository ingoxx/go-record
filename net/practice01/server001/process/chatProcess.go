package process

import (
	"encoding/json"
	"fmt"
	"net"

	Common "github.com/ingoxx/Golang-practise/net/practice01/common"
	"github.com/ingoxx/Golang-practise/net/practice01/model"
)

var (
	Chp *ChatProcessor
)

func init() {
	Chp = &ChatProcessor{}
}

type ChatProcessor struct {
}

func (c *ChatProcessor) SendMsgToAllUsers(mes *model.Message) {
	//遍历客户端维护的一个map,一个在线人数的列表Conn,然后获取每一个Conn转发消息
	cr := model.ChatMessage{}
	err := json.Unmarshal([]byte(mes.Data), &cr)
	if err != nil {
		fmt.Println("转发消息反序列化失败 err=", err)
		return
	}

	for id, up := range UserMgr.OnLineUser {
		if id == cr.UserId { //服务端转发的时候不推送到自己
			continue
		}
		c.SendMsgAll(cr, up.Conn)
	}
}

func (c *ChatProcessor) SendMsgAll(ctr model.ChatMessage, conn net.Conn) {

	mes := model.Message{}
	cr := model.ChatMessageResult{}

	mes.Type = model.ChatMessageResultType
	cr.UserId = ctr.UserId
	cr.Content = ctr.Content

	data, err := json.Marshal(&cr)
	if err != nil {
		fmt.Println("发消息序列化ChatMessageResult失败, err=", err)
		return
	}

	mes.Data = string(data)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("发消息序列化Message失败, err=", err)
		return
	}

	t := &Common.TransData{
		Conn: conn,
	}
	err = t.SendMessage(data)
	if err != nil {
		fmt.Println("发消息失败, err= ", err)
		return
	}
}

func (c *ChatProcessor) SendMsgToOneOnLingUser(mes *model.Message) {
	//私聊
	cr := model.ChatUserToUserMessage{}
	err := json.Unmarshal([]byte(mes.Data), &cr)
	if err != nil {
		fmt.Println("转发消息反序列化失败 err=", err)
		return
	}
	sender := cr.Sender
	up, k := UserMgr.OnLineUser[cr.UserId]
	if !k {
		//保存发给离线用户的消息
		sen, senk := UserMgr.OnLineUser[sender]
		if senk {
			c.SendMsgToOffLingUser(cr, sen.Conn)
		}
	} else {
		//发给在线用户
		c.SendMsgToOnLingUser(cr, up.Conn)
	}
}

func (c *ChatProcessor) SendMsgToOnLingUser(ctr model.ChatUserToUserMessage, conn net.Conn) {

	mes := model.Message{}
	cr := model.ChatUserToUserMessageResult{}

	mes.Type = model.ChatUserToUserMessageResultType
	cr.UserId = ctr.UserId
	cr.Content = ctr.Content

	data, err := json.Marshal(&cr)
	if err != nil {
		fmt.Println("发消息序列化ChatUserToUserMessageResult失败, err=", err)
		return
	}

	mes.Data = string(data)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("发消息序列化Message失败, err=", err)
		return
	}

	t := &Common.TransData{
		Conn: conn,
	}
	err = t.SendMessage(data)
	if err != nil {
		fmt.Println("发消息失败, err= ", err)
		return
	}
}

func (c *ChatProcessor) SendMsgToOffLingUser(ctr model.ChatUserToUserMessage, conn net.Conn) {
	mes := model.Message{}
	cr := model.ChatUserToUserMessageResult{}
	sm := model.ChatUserToUserMessage{}

	sm.Sender = ctr.Sender
	sm.Content = ctr.Content
	sm.UserId = ctr.UserId

	Mm.AddMessage(sm)
	mes.Type = model.ChatUserToUserMessageResultType

	cr.UserId = ctr.UserId
	cr.Content = "离线消息已发送"

	data, err := json.Marshal(&cr)
	if err != nil {
		fmt.Println("序列化ChatUserToUserMessageResult失败,err=", err)
		return
	}

	mes.Data = string(data)

	data, err = json.Marshal(&mes)
	if err != nil {
		fmt.Println("序列化Message失败,err=", err)
		return
	}

	t := &Common.TransData{
		Conn: conn,
	}
	err = t.SendMessage(data)
	if err != nil {
		fmt.Println("发消息失败, err= ", err)
		return
	}

}
