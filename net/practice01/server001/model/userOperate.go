package model

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/ingoxx/Golang-practise/net/practice01/model"

	"github.com/go-redis/redis"
)

// 这里设置全局变量是作用是,在main()执行的时候就可以给UserCdus创建实例,这样就不用每次调用UserCdus方法时都创建实例跟工厂模式配合使用
var (
	IuserCdus *UserCdus
)

// 这里是放用户增删改查的对象,也就是对UserInfo结构体的操作
type UserCdus struct {
	pool *redis.Client
}

func NewPoolRds(pool *redis.Client) *UserCdus {
	return &UserCdus{
		pool: pool,
	}
}

func (u *UserCdus) GetUserId(id int) (user *UserInfo, err error) {
	data, err := u.pool.HGet("users", strconv.FormatInt(int64(id), 10)).Result()
	if err != nil {
		err = ErrorUserNotExists
		return
	}

	user = &UserInfo{}

	err = json.Unmarshal([]byte(data), user)
	if err != nil {
		err = errors.New("反序列化失败")
		return
	}
	return
}

func (u *UserCdus) UserLogin(ui int, pwd string) (user *UserInfo, err error) {
	user, err = u.GetUserId(ui)
	if err != nil {
		return
	}

	//userid存在往下校验密码是否正确
	if pwd != user.Passwd {
		err = ErrorUserPwd
		return
	}

	return
}

func (u *UserCdus) UserRegister(ur *model.RegisterMessage) (err error) {
	// defer u.pool.Close() 这里不能关,不然只有第一个连接才能读取redis的数据
	_, err = u.GetUserId(ur.UserId)
	if err == nil {
		err = ErrorUserExists
		return
	}

	//写入到redis
	userId := ur.UserId
	data, err := json.Marshal(ur)
	if err != nil {
		err = ErrorJsonMarshal
		return
	}

	_, err = u.pool.HSet("users", strconv.FormatInt(int64(userId), 10), string(data)).Result()
	if err != nil {
		err = ErrorRegister
		return
	}

	return
}

// 当用户登录成功后,将所有离线用户也推送给所有已经登录的客户端,如果都上线了就不推送
func (u *UserCdus) GetAllUserId() (ul []int, err error) {
	data, err := u.pool.HGetAll("users").Result()
	if err != nil {
		err = errors.New("获取所有用户id失败")
		return
	}
	for id := range data {
		uid, _ := strconv.ParseInt(id, 10, 64)
		ul = append(ul, int(uid))
	}
	return
}
