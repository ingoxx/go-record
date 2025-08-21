package process

//这里从总控process.go分发到这里, 是处理跟用户相关的逻辑,如登录,注册等

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	Common "github.com/ingoxx/Golang-practise/net/practice01/common"
	UserMessage "github.com/ingoxx/Golang-practise/net/practice01/model"
	UserModeler "github.com/ingoxx/Golang-practise/net/practice01/server001/model"
)

type UserProcessor struct {
	Conn   net.Conn
	UserId int //这个字段的作用是标明当前的UserProcessor属于哪个UserId
}

// 用户登录
func (u *UserProcessor) LoginProcess(mes *UserMessage.Message) (err error) {
	lm := UserMessage.LoginMessage{}
	data := json.Unmarshal([]byte(mes.Data), &lm)
	if data != nil {
		err = errors.New("反序列化登录消息错误")
		return
	}

	//重新封装传输的消息数据结构
	resMsg := UserMessage.Message{}
	resMsg.Type = UserMessage.LoginResultType

	//返回的处理完后的LoginResult结构体
	rm := UserMessage.LoginResult{}

	//从redis获取数据进行校验
	user, err := UserModeler.IuserCdus.UserLogin(lm.UserId, lm.Passwd)
	fmt.Println("登录请求获取的UserInfo = ", user)
	if err != nil {
		fmt.Println("登录失败err = ", err)
		switch err {
		case UserModeler.ErrorUserExists:
			rm.Code = 300
			rm.Err = err.Error()
		case UserModeler.ErrorUserNotExists:
			rm.Code = 400
			rm.Err = err.Error()
		case UserModeler.ErrorUserPwd:
			rm.Code = 500
			rm.Err = err.Error()
		default:
			rm.Code = 600
			rm.Err = "未知错误"
		}
	} else {
		rm.Code = 200
		//把已经登录成功的user conn放到UserManage
		u.UserId = lm.UserId
		allUsers, _ := UserModeler.IuserCdus.GetAllUserId()
		UserMgr.AddOnlineUser(u)
		UserMgr.DelOffLineUserList(u.UserId, allUsers)
		//将当前用户id放入到rm.Users中然后一起将LoginResult返回给客户端
		for uid := range UserMgr.OnLineUser {
			rm.UsersId = append(rm.UsersId, uid)
		}
		rm.AllUsers = UserMgr.OffLineUser

		fmt.Println("rm.UsersId = ", rm.UsersId)
		fmt.Println("登录成功")
	}

	//将rm序列化放到resMsg.Data中
	j, e := json.Marshal(&rm)
	if e != nil {
		err = errors.New("rm 序列化失败")
		return
	}
	resMsg.Data = string(j)

	//重新序列化数据传输的唯一结构
	j, e = json.Marshal(&resMsg)
	if e != nil {
		err = errors.New("resMsg 序列化失败")
		return
	}

	//返回登录消息验证后的结果给客户端
	t := &Common.TransData{
		Conn: u.Conn,
	}
	err = t.SendMessage(j)
	if err != nil {
		fmt.Println("SendMessage err = ", err)
	}

	//登录成功通知所有在线的用户自己上线了
	if rm.Code == 200 {
		u.PushUserStatus(lm.UserId, UserMessage.OnLine)
		u.PushOffLineMessage()
	}

	return
}

// 用户注册
func (u *UserProcessor) RegisterProcess(mes *UserMessage.Message) (err error) {
	lm := UserMessage.RegisterMessage{}
	data := json.Unmarshal([]byte(mes.Data), &lm)
	if data != nil {
		err = errors.New("反序列化消息错误")
		return
	}

	//重新封装传输的消息数据结构
	resMsg := UserMessage.Message{}
	resMsg.Type = UserMessage.RegisterResultType

	//返回的处理完后的RegisterResult结构体
	rm := UserMessage.RegisterResult{}

	//从redis获取数据进行校验
	err = UserModeler.IuserCdus.UserRegister(&lm)
	if err != nil {
		fmt.Println("注册失败, err = ", err)
		switch err {
		case UserModeler.ErrorUserExists:
			rm.Code = 300
			rm.Err = err.Error()
		case UserModeler.ErrorUserNotExists:
			rm.Code = 400
			rm.Err = err.Error()
		case UserModeler.ErrorUserPwd:
			rm.Code = 500
			rm.Err = err.Error()
		case UserModeler.ErrorRegister:
			rm.Code = 700
			rm.Err = err.Error()
		default:
			rm.Code = 600
			rm.Err = "注册出现未知错误"
		}
	} else {
		rm.Code = 200
		fmt.Println("注册成功")
	}

	//将rm序列化放到resMsg.Data中
	j, e := json.Marshal(&rm)
	if e != nil {
		err = errors.New("rm 序列化失败")
		return
	}
	resMsg.Data = string(j)

	//重新序列化数据传输的唯一结构
	j, e = json.Marshal(&resMsg)
	if e != nil {
		err = errors.New("resMsg 序列化失败")
		return
	}

	//返回注册结果给客户端
	t := &Common.TransData{
		Conn: u.Conn,
	}
	err = t.SendMessage(j)
	if err != nil {
		fmt.Println("SendMessage err = ", err)
	}
	return
}

func (u *UserProcessor) LoginOut(mes *UserMessage.Message) (err error) {
	us := UserMessage.UserStatusChange{}
	err = json.Unmarshal([]byte(mes.Data), &us)
	if err != nil {
		err = errors.New("接收的用户状态修改消息反序列化失败")
		return
	}
	UserMgr.DelOnlineUser(us.UserId)
	UserMgr.AddOffLineUserList(us.UserId)
	u.PushUserStatus(us.UserId, UserMessage.Offline)

	return
}

// 推送所有在线用户给所有已登录的用户(除了自己)
func (u *UserProcessor) PushUserStatus(uid, status int) {
	//遍历UserManage的OnLineUser,然后推送UserStatusChange消息给所有在线用户
	for id, up := range UserMgr.OnLineUser { //UserMgr.OnLineUser这里存的就是所有在线的用户,遍历他们目的是为了通知他们有用户登录上线了
		//这里的up是服务端保存的每个客户端的Conn
		//过滤掉自己
		if id == uid {
			continue
		}
		//起单独的方法来推送
		up.PushMsgToAllUser(uid, status) //将自己推送给除了自己以外的所有在线的用户
	}
}

// 推送给所有在线用户
func (u *UserProcessor) PushMsgToAllUser(uid int, status int) {
	mes := UserMessage.Message{}
	mes.Type = UserMessage.UserStatusChangeType

	us := UserMessage.UserStatusChange{}
	us.UserId = uid
	us.UserStatus = status
	us.OfflineList = UserMgr.OffLineUser
	fmt.Println("us.OfflineList = ", UserMgr.OffLineUser)

	data, err := json.Marshal(&us)
	if err != nil {
		fmt.Println("推送消息序列化us失败 ,err = ", err)
		return
	}

	mes.Data = string(data)

	fmt.Println("推送消息 = ", mes.Data)

	data, err = json.Marshal(&mes)
	if err != nil {
		fmt.Println("推送消息序列化mes失败 ,err = ", err)
		return
	}

	t := &Common.TransData{ //这里每次都要去创建实例的原因是,每个已经登录的用户conn都是不同的,如果将这个conn设置为全局,登录的永远都是只有一个人
		Conn: u.Conn,
	}
	err = t.SendMessage(data)
	if err != nil {
		fmt.Println("推送消息失败 err= ", err)
		return
	}
}

// 给登录的用户推送离线消息
func (u *UserProcessor) PushOffLineMessage() {
	mes := UserMessage.Message{}
	ochr := UserMessage.ChatRrcordMessageResult{}
	mes.Type = UserMessage.ChatRrcordMessageResultType
	ochr.ChatRecord = Mm.MessageList[u.UserId]

	data, err := json.Marshal(&ochr)
	if err != nil {
		fmt.Println("推送消息序列化us失败 ,err = ", err)
		return
	}

	mes.Data = string(data)

	fmt.Println("推送消息 = ", mes.Data)

	data, err = json.Marshal(&mes)
	if err != nil {
		fmt.Println("推送消息序列化mes失败 ,err = ", err)
		return
	}

	t := &Common.TransData{ //这里每次都要去创建实例的原因是,每个已经登录的用户conn都是不同的,如果将这个conn设置为全局,登录的永远都是只有一个人
		Conn: u.Conn,
	}
	err = t.SendMessage(data)
	if err != nil {
		fmt.Println("推送消息失败 err= ", err)
		return
	}

}
