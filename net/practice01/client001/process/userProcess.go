package process

import (
	"encoding/json"
	"errors"
	"fmt"

	"net"

	Common "github.com/ingoxx/Golang-practise/net/practice01/common"
	UserMessage "github.com/ingoxx/Golang-practise/net/practice01/model"
)

type UserProcessor struct {
	//需要再加字段
}

func (up *UserProcessor) Login(u int, p string) (err error) {

	//定协议
	c, e1 := net.Dial("tcp", "101.35.143.86:8092")
	if e1 != nil {
		// fmt.Println(e1)
		err = fmt.Errorf("登录失败, err=%v", e1)
		return
	}
	// defer c.Close()

	m := UserMessage.Message{}
	m.Type = UserMessage.LoginMessAgeType

	lm := UserMessage.LoginMessage{}
	lm.UserId = u
	lm.Passwd = p

	//将消息本体(LoginMessage)序列化
	data, e2 := json.Marshal(&lm)
	if e2 != nil {
		fmt.Println("序列化失败=", e2)
		return
	}

	//将序列化好的消息本体放到Message.Data
	m.Data = string(data)

	//最终将Message序列化发送到服务器
	data, e2 = json.Marshal(&m)
	if e2 != nil {
		fmt.Println("序列化失败=", e2)
		return
	}

	t := &Common.TransData{
		Conn: c,
	}

	//发送消息本体
	err = t.SendMessage(data)
	if err != nil {
		err = errors.New("SendMessage errors")
		return
	}

	//处理服务器发回来的消息
	mes, e3 := t.RecvMessage()
	if e3 != nil {
		err = errors.New("RecvMessage errors")
		return
	}

	rm := UserMessage.LoginResult{}
	e3 = json.Unmarshal([]byte(mes.Data), &rm)
	if e3 != nil {
		err = fmt.Errorf("序列化mes.Data失败, err = %v", e3)
		return
	}

	if rm.Code != 200 {
		err = errors.New(rm.Err)
		return
	}

	//初始化CurUsers(群发消息用)
	CurUsers.Conn = c //保存跟服务器建立的Conn
	CurUsers.UserId = u
	CurUsers.UserStatus = UserMessage.OnLine

	//这里初始化离线用户变量
	OfflineUser = rm.AllUsers

	//登录成功显示其他在线用户(除了自己)列表
	for _, v := range rm.UsersId {
		if v == u {
			continue
		}
		us := &UserMessage.UserInfor{
			UserId:     v,
			UserStatus: UserMessage.OnLine,
		}
		//初始化OnLineUser,服务端推送一个就添加一个
		OnLineUser[v] = us
	}

	//登录成功后要启动一个协程, 专门用来保持和服务端通信,接收服务器端发送的消息如单聊,群聊信息
	go KeepConnectServer(c)

	return

}

func (up *UserProcessor) Register(u int, p string) (err error) {

	//定协议
	c, e1 := net.Dial("tcp", "192.168.11.188:8092")
	if e1 != nil {
		fmt.Println(e1)
		return
	}

	// defer c.Close()

	m := UserMessage.Message{}
	m.Type = UserMessage.RegisterType

	lm := UserMessage.RegisterMessage{}
	lm.UserId = u
	lm.Passwd = p

	//将消息本体(LoginMessage)序列化
	data, e2 := json.Marshal(&lm)
	if e2 != nil {
		fmt.Println("序列化失败=", e2)
		return
	}

	//将序列化好的消息本体放到Message.Data
	m.Data = string(data)

	//最终将Message序列化发送到服务器
	data, e2 = json.Marshal(&m)
	if e2 != nil {
		fmt.Println("序列化失败=", e2)
		return
	}

	fmt.Println("客户端发送的消息本体=", string(data))

	t := &Common.TransData{
		Conn: c,
	}

	//转成可以表示长度的byte切片
	err = t.SendMessage(data)
	if err != nil {
		err = errors.New("SendMessage errors")
		return
	}

	//处理服务器发回来的消息
	mes, e3 := t.RecvMessage()
	if e3 != nil {
		err = errors.New("RecvMessage errors")
		return
	}

	rm := UserMessage.RegisterResult{}
	e3 = json.Unmarshal([]byte(mes.Data), &rm)
	if e3 != nil {
		err = errors.New("RecvMessage errors")
		return
	}

	if rm.Code != 200 {
		err = errors.New(rm.Err)
		return
	}

	//注册成功后要启动一个协程, 专门用来保持和服务端通信,接收服务器端发送的消息并显示在界面,如单聊,群聊信息
	go KeepConnectServer(c)

	return

}

func (up *UserProcessor) LoginOut() (err error) {
	mes := UserMessage.Message{}
	us := UserMessage.UserStatusChange{}
	us.UserStatus = 1
	us.UserId = CurUsers.UserId

	data, err := json.Marshal(&us)
	if err != nil {
		err = fmt.Errorf("UserStatusChange消息序列化失败,err=%v", err)
		return
	}
	mes.Type = UserMessage.UserStatusChangeType
	mes.Data = string(data)

	data, err = json.Marshal(&mes)
	if err != nil {
		err = fmt.Errorf("Message消息序列化失败,err=%v", err)
		return
	}

	t := &Common.TransData{
		Conn: CurUsers.Conn,
	}

	t.SendMessage(data)
	if err != nil {
		err = errors.New("发送退出消息失败")
		return
	}

	return
}
