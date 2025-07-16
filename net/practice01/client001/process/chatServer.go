package process

import (
	"encoding/json"
	"fmt"
	"net"

	Common "github.com/ingoxx/Golang-practise/net/practice01/common"
	"github.com/ingoxx/Golang-practise/net/practice01/model"
)

// 这里是保持跟服务器的连接,接收服务器端发送的消息并显示在界面,如单聊,群聊信息
func LoggedMenu() {
	//显示登录之后的页面
	ml := []string{"显示在线用户", "发送信息", "信息列表", "退出登录"}
	id := 0
	out := true
	// content := ""
	//只要不退出实例就可以一直用,就放在for外边
	cp := &ChatProcess{}

	for out {
		fmt.Println("\t-----------登录成功,欢迎xxx-----------")
		for i, v := range ml {
			fmt.Println("\t", i, v)
		}
		fmt.Printf("\t请输入%d-%d:", 0, len(ml)-1)
		fmt.Scanln(&id)
		if id < 0 || id > len(ml)-1 {
			fmt.Println()
			fmt.Printf("\t输入错误,请输入%d-%d\n", 0, len(ml)-1)
			continue
		}
		switch ml[id] {
		case "显示在线用户":
			ShowOnlineUser()
		case "发送信息":
			fmt.Println("\t----------发送消息---------")
			cp.SendMsg()
		case "信息列表":
			fmt.Println("\t------------信息列表---------")
			Cr.ShowChatRecord()
		case "退出登录":
			fmt.Println("\tbye")
			up := &UserProcessor{}
			up.LoginOut()
			out = false
		}
	}
}

// 和服务器端保持通信接收服务端发来的信息
func KeepConnectServer(conn net.Conn) {
	t := &Common.TransData{
		Conn: conn,
	}

	for {
		mes, err := t.RecvMessage()
		if err != nil {
			fmt.Println("接受服务端消息失败 err=", err)
			return
		}
		switch mes.Type {
		case model.UserStatusChangeType:
			//这里是获取服务端推送的消息,目前推送的是用户上线消息
			//客户端也要维护一个用户在线的map
			ustc := model.UserStatusChange{}
			err := json.Unmarshal([]byte(mes.Data), &ustc)
			if err != nil {
				fmt.Println("反序列化服务端推送的消息失败,err=", err)
				return
			}
			OfflineUser = ustc.OfflineList
			//把在线用户列表保存到客户端的OnLineUser中
			UpdateUserStatus(&ustc)
		case model.ChatMessageResultType:
			//服务端转发的群消息
			RecvGroupMsg(&mes)
		case model.ChatUserToUserMessageResultType:
			//服务端转发的私聊消息
			RecvUserMsg(&mes)
		case model.ChatRrcordMessageResultType:
			//服务端推送的聊天记录
			RecvChatRecord(&mes)
		case model.FileMessageType:
			//接收文件
			RecvFile(&mes)
		default:
			fmt.Println("消息类型不存在")
		}
	}
}
