package process

import (
	"fmt"
)

//当用户登录成功显示在线用户列表,目的是为了获取每个用户的conn
var (
	UserMgr *UserManage
)

func init() {
	UserMgr = &UserManage{
		OnLineUser:  make(map[int]*UserProcessor, 1024),
		OffLineUser: []int{},
	}
}

//map的key是唯一的,如果存在key则会更新,不存在则添加
type UserManage struct {
	OnLineUser  map[int]*UserProcessor
	OffLineUser []int
}

//添加在线用户
func (u *UserManage) AddOnlineUser(up *UserProcessor) {
	u.OnLineUser[up.UserId] = up
}

//删除在线用户
func (u *UserManage) DelOnlineUser(uid int) {
	delete(u.OnLineUser, uid)
}

//获取所有在线用户
func (u *UserManage) GetAllOnlineUser(uid int) map[int]*UserProcessor {
	return u.OnLineUser
}

//根据userid返回对应的UserProcessor(其实主要是为了获取Conn   net.Conn)
func (u *UserManage) GetOnlineUserByUid(uid int) (up *UserProcessor, err error) {
	up, ok := u.OnLineUser[uid]
	if !ok {
		err = fmt.Errorf("%d用户不在线", uid)
		return
	}
	return
}

func (u *UserManage) DelOffLineUserList(uid int, ul []int) {
	delIndex := 0
	fmt.Println("ul= ", ul)
	fmt.Println("u.OffLineUser= ", u.OffLineUser)
	if len(u.OffLineUser) == 0 {
		for k, v := range ul {
			if v == uid {
				delIndex = k
				break
			}
		}
		u.OffLineUser = append(ul[0:delIndex], ul[delIndex+1:]...)
	} else {
		for k, v := range u.OffLineUser {
			if v == uid {
				delIndex = k
				break
			}
		}
		u.OffLineUser = append(u.OffLineUser[0:delIndex], u.OffLineUser[delIndex+1:]...)
	}
}

func (u *UserManage) AddOffLineUserList(uid int) {
	u.OffLineUser = append(u.OffLineUser, uid)
}
