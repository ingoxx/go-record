package redisServer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

type Rds struct {
	conn   *redis.Client
	openid string
	Err    error
	lock   *sync.Mutex
}

func NewRds(openid string) *Rds {
	var rdPool = redis.NewClient(
		&redis.Options{
			Addr:         "127.0.0.1:6377",
			DB:           10,
			MinIdleConns: 5,
			Password:     "chatai",
			PoolSize:     5,
			PoolTimeout:  30 * time.Second,
			DialTimeout:  1 * time.Second,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	)

	return &Rds{
		conn:   rdPool,
		openid: openid,
		Err:    rdPool.Ping().Err(),
		lock:   new(sync.Mutex),
	}
}

func (r *Rds) CheckInvite(userid string) (err error) {
	iE := r.CheckKey("invite")
	if !iE {
		return errors.New("已被限制访问")
	}

	var data = make(map[string][]string)
	check, err := r.conn.HGet("invite", r.openid).Result()
	err = json.Unmarshal([]byte(check), &data)
	if err != nil {
		return
	}

	var isExists bool
	user, ok := data["user"]
	if ok {
		for _, v := range user {
			if v == userid {
				isExists = true
				break
			}
		}

		if isExists {
			return errors.New("同一个用户只能邀请一次")
		}

		user = append(user, userid)
		data["user"] = user

	}

	b, err := json.Marshal(&data)
	if err != nil {
		return
	}

	err = r.conn.HSet("invite", r.openid, b).Err()
	if err != nil {
		return
	}

	return
}

func (r *Rds) isNextDay(n, o time.Time) bool {
	ny, nm, nd := n.Date()
	oy, om, od := o.Date()
	if ny != oy || nm != om || nd != od {
		return true
	}
	return false
}

func (r *Rds) updateInviteDate() (err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	var data = make(map[string]int)
	quota, err := r.conn.HGet("quota", r.openid).Result()
	err = json.Unmarshal([]byte(quota), &data)
	if err != nil {
		return
	}

	date := data["time"]
	now := time.Now().Unix()
	nft := time.Unix(now, 0)
	rft := time.Unix(int64(date), 0)

	if r.isNextDay(nft, rft) {
		data["invite"] = 0
		data["add"] = 2
		data["finished"] = 2
		data["time"] = int(now)
	}

	b, err := json.Marshal(&data)
	if err != nil {
		return
	}

	err = r.conn.HSet("quota", r.openid, b).Err()
	if err != nil {
		return
	}

	return
}

func (r *Rds) Invite(userid string) (i int, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if err = r.CheckInvite(userid); err != nil {
		return
	}

	iE := r.CheckKey("quota")
	if !iE {
		return 0, errors.New("已被限制访问")
	}

	var data = make(map[string]int)
	quota, err := r.conn.HGet("quota", r.openid).Result()
	err = json.Unmarshal([]byte(quota), &data)
	if err != nil {
		return
	}

	invite, ok := data["invite"]
	date := data["time"]
	chatgpt := data["chatgpt"]
	bd := data["bd"]
	gemini := data["gemini"]
	qw := data["qw"]
	isAdd := data["add"]
	if !ok {
		data["invite"] = 1
		data["time"] = int(time.Now().Unix()) + 86400
		// 1：当天额度已加，2：当天额度未加
		data["add"] = 2
		// 1：当天已完成邀请数，2：当天未完成邀请数
		data["finished"] = 2
	}

	// 第二天后邀请数重新归0
	now := int(time.Now().Unix())
	nft := time.Unix(time.Now().Unix(), 0)
	rft := time.Unix(int64(date), 0)

	if invite < 5 && !r.isNextDay(nft, rft) {
		invite++
		if invite == 5 && isAdd == 2 {
			chatgpt++
			data["chatgpt"] = chatgpt
			bd++
			data["bd"] = bd
			gemini++
			data["gemini"] = gemini
			qw++
			data["qw"] = qw
			data["add"] = 1
			data["finished"] = 1
		}
		data["invite"] = invite
	} else {
		if r.isNextDay(nft, rft) {
			data["invite"] = 0
			data["add"] = 2
			data["finished"] = 2
			data["time"] = now
		}
	}

	b, err := json.Marshal(&data)
	if err != nil {
		return
	}

	err = r.conn.HSet("quota", r.openid, b).Err()
	if err != nil {
		return
	}

	i = invite

	return
}

func (r *Rds) CheckKey(key string) bool {
	isExists, _ := r.conn.HExists(key, r.openid).Result()
	return isExists
}

func (r *Rds) SetInvite() (err error) {
	isExists := r.CheckKey("invite")
	if !isExists {
		var initData = []string{"lxb"}
		var data = map[string][]string{
			"user": initData,
		}

		b, err := json.Marshal(&data)
		if err != nil {
			return err
		}
		err = r.conn.HSet("invite", r.openid, b).Err()
		if err != nil {
			return err
		}
	}

	return
}

func (r *Rds) SetQuota() (err error) {
	isExists := r.CheckKey("quota")
	if !isExists {
		var data = map[string]int{
			"chatgpt": 0,
			"gemini":  20,
			"bd":      5,
			"qw":      5,
		}

		b, err := json.Marshal(&data)
		if err != nil {
			return err
		}
		err = r.conn.HSet("quota", r.openid, b).Err()
		if err != nil {
			return err
		}
	}

	return
}

func (r *Rds) GetQuota() (data map[string]int, err error) {
	if err = r.updateInviteDate(); err != nil {
		return
	}

	quota, err := r.conn.HGet("quota", r.openid).Result()
	err = json.Unmarshal([]byte(quota), &data)
	if err != nil {
		return
	}

	return
}

func (r *Rds) UpdateQuota(model string) (err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	iE := r.CheckKey("quota")
	if !iE {
		return errors.New("已被限制访问")
	}

	var data = make(map[string]int)
	quota, err := r.conn.HGet("quota", r.openid).Result()
	err = json.Unmarshal([]byte(quota), &data)
	if err != nil {
		return
	}

	q := data[model]
	if q <= 0 {
		return errors.New(fmt.Sprintf("%s额度已用完, 请在个人中心里点击联系我们增加额度", model))
	}

	q--
	data[model] = q

	b, err := json.Marshal(&data)
	if err != nil {
		return
	}

	err = r.conn.HSet("quota", r.openid, b).Err()
	if err != nil {
		return
	}

	return
}

func (r *Rds) Get() (err error) {
	if err = r.Set(); err != nil {
		return
	}
	return
}

func (r *Rds) Check() (dd string, err error) {
	var isExists bool
	var data []string
	dd, err = r.conn.Get("openid").Result()
	if err != nil {
		err = nil
		return
	}

	if err = json.Unmarshal([]byte(dd), &data); err != nil {
		return
	}

	for _, v := range data {
		if v == r.openid {
			isExists = true
			break
		}
	}

	if isExists {
		err = errors.New("openid已存在")
		return
	}

	return
}

func (r *Rds) Set() (err error) {
	var ds []string
	data, err := r.Check()
	if err != nil {
		return
	}

	if data != "" {
		if err = json.Unmarshal([]byte(data), &ds); err != nil {
			return
		}
	} else {
		ds = make([]string, 0, 10)
	}

	ds = append(ds, r.openid)

	bs, err := json.Marshal(&ds)
	if err != nil {
		return
	}

	err = r.conn.Set("openid", bs, time.Duration(31536*10000)*time.Second).Err()
	if err != nil {
		return
	}

	return

}

func (r *Rds) CheckOpenId() (err error) {
	var isExists bool
	var data []string
	ds, err := r.conn.Get("openid").Result()
	if err != nil {
		return
	}

	if err = json.Unmarshal([]byte(ds), &data); err != nil {
		return
	}

	for _, v := range data {
		if v == r.openid {
			isExists = true
			break
		}
	}

	if isExists {
		return
	}

	return errors.New("openid不合法")
}
