package process

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	Common "github.com/ingoxx/Golang-practise/net/practice01/common"
	UserMessage "github.com/ingoxx/Golang-practise/net/practice01/model"
)

type ChatProcess struct {
}

// 发送群聊的消息
func (c *ChatProcess) SendAllToUser(content string) (err error) {
	mes := UserMessage.Message{}
	mes.Type = UserMessage.ChatMessageType

	chat := UserMessage.ChatMessage{}
	chat.Content = content
	chat.UserId = CurUsers.UserId //这里是登录成功后初始化的CurrentUser实例
	chat.UserStatus = CurUsers.UserStatus

	data, err := json.Marshal(&chat)
	if err != nil {
		err = fmt.Errorf("群发消息序列化失败 err=%v", err)
		return
	}

	mes.Data = string(data)
	data, err = json.Marshal(&mes)
	if err != nil {
		err = fmt.Errorf("传输消息序列化失败 err=%v", err)
		return
	}

	t := &Common.TransData{
		Conn: CurUsers.Conn,
	}

	err = t.SendMessage(data)
	if err != nil {
		err = errors.New("发送群消息失败")
		return
	}

	return

}

// 发送私聊的消息
func (c *ChatProcess) SendMsgUserToUser(uid int, content string) (err error) {
	mes := UserMessage.Message{}
	mes.Type = UserMessage.ChatUserToUserMessageType

	chat := UserMessage.ChatUserToUserMessage{}
	chat.Content = content
	chat.UserId = uid
	chat.Sender = CurUsers.UserId

	data, err := json.Marshal(&chat)
	if err != nil {
		err = fmt.Errorf("群发消息序列化失败 err=%v", err)
		return
	}

	mes.Data = string(data)
	data, err = json.Marshal(&mes)
	if err != nil {
		err = fmt.Errorf("传输消息序列化失败 err=%v", err)
		return
	}

	t := &Common.TransData{
		Conn: CurUsers.Conn,
	}

	err = t.SendMessage(data)
	if err != nil {
		err = errors.New("发送群消息失败")
		return
	}

	return

}

// 专门用于发送消息的方法
func (c *ChatProcess) SendMsg() (err error) {
	id := 0
	uid := 0
	messagetype := 0
	file := ""
	out := true
	cp := &ChatProcess{}
	for out {
		fmt.Println("\t1 在线用户")
		fmt.Println("\t2 离线用户")
		fmt.Println("\t3 群聊")
		fmt.Println("\t4 退出")
		fmt.Printf("请选择:")
		fmt.Scanln(&id)
		switch id {
		case 1:
			//在线用户
			status := false
			ShowOnlineUser()
			fmt.Println()
			fmt.Printf("\t请选择用户id进行聊天:")
			fmt.Scanln(&uid)
			for u := range OnLineUser {
				if uid == u {
					status = true
				}
			}
			if !status {
				fmt.Println("\t你选择的用户id不存在,请重新选择")
			} else {
				fmt.Println("\t1 发送文件")
				fmt.Println("\t2 发送内容")
				fmt.Println("请选择发送的消息类型:")
				fmt.Scanln(&messagetype)
				switch messagetype {
				case 1:
					fmt.Printf("请输入文件路径:")
					fmt.Scanln(&file)
					fp := &FileProcess{}
					fp.SendFile(file, uid)
				case 2:
					fmt.Printf("对%d发消息:", uid)
					reader := bufio.NewReader(os.Stdin)
					content, _ := reader.ReadString('\n')
					content = strings.ReplaceAll(content, "\r\n", "")
					fmt.Println(content)
					c.SendMsgUserToUser(uid, content)
				default:
					fmt.Println("选择错误")
				}

			}
		case 2:
			//离线用户
			status := false
			ShowOfflineListUser()
			fmt.Println()
			fmt.Println()
			fmt.Printf("\t请选择用户id进行聊天:")
			fmt.Scanln(&uid)
			for _, u := range OfflineUser {
				if uid == u {
					status = true
				}
			}
			if !status {
				fmt.Println("\t你选择的用户id不存在,请重新选择")
			} else {
				fmt.Printf("对%d发离线消息:", uid)
				reader := bufio.NewReader(os.Stdin)
				content, _ := reader.ReadString('\n')
				content = strings.ReplaceAll(content, "\r\n", "")
				fmt.Println(content)
				c.SendMsgUserToUser(uid, content)
			}
		case 3:
			fmt.Printf("对大家群发消息:")
			reader := bufio.NewReader(os.Stdin)
			content, _ := reader.ReadString('\n')
			content = strings.ReplaceAll(content, "\r\n", "")
			cp.SendAllToUser(content)
		case 4:
			out = false
		default:
			fmt.Println("选择错误,请重新选择")
		}
	}
	return
}
