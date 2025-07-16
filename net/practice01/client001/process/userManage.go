package process

import (
	"fmt"

	CurUsersModel "github.com/ingoxx/Golang-practise/net/practice01/client001/model"
	UserMessage "github.com/ingoxx/Golang-practise/net/practice01/model"
)

// 客户端也维护的一个用户列表,map类型
var OnLineUser map[int]*UserMessage.UserInfor = make(map[int]*UserMessage.UserInfor, 10)
var OfflineUser []int

var CurUsers CurUsersModel.CurrentUser //在用户登录成功后,就初始化CurUsers, 维护自己的Conn,及要把当前跟服务端的Conn保存起来

// 在客户端展示在线用户列表
func ShowOnlineUser() {
	fmt.Println("\t----------在线用户---------")
	for uid, up := range OnLineUser {
		if up.UserStatus == 0 {
			fmt.Printf("\t在线用户id=%d\n", uid)
		}
	}
}

// 在客户端展示离线用户列表
func ShowOfflineUser() {
	fmt.Println("\t----------离线用户---------")
	for uid, up := range OnLineUser {
		if up.UserStatus == 1 {
			fmt.Printf("\t离线用户id=%d\n", uid)
		}
	}
}

func ShowOfflineListUser() {
	fmt.Println("\t----------离线用户列表---------")
	for _, v := range OfflineUser {
		fmt.Printf("\t离线用户id=%d\n", v)
	}
}

// 获取服务端推送的UserStatusChange消息
func UpdateUserStatus(mes *UserMessage.UserStatusChange) {
	u, ok := OnLineUser[mes.UserId]
	if !ok {
		u = &UserMessage.UserInfor{
			UserId: mes.UserId,
		}
	}
	u.UserStatus = mes.UserStatus
	OnLineUser[mes.UserId] = u
	//接收到用户在线/离线的消息后直接展示
	if mes.UserStatus == 0 {
		ShowOnlineUser()
	} else if mes.UserStatus == 1 {
		ShowOfflineUser()
	}
}
