package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	cuerr "github.com/ingoxx/go-record/http/wx/error"
	"github.com/ingoxx/go-record/http/wx/form"
	"github.com/ingoxx/go-record/http/wx/qqMapApi"
	"golang.org/x/crypto/bcrypt"
	"log"
	"sync"
	"time"
)

var (
	rds          *redis.Client
	groupKey     = "group_id"
	AddrListKey  = "addr_check_list"
	WxOPenIdKey  = "wx_open_id_list"
	onlineKey    = "online"
	joinGroupKey = "join"
	defaultImg   = "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/bk3.svg"
	//cusAddrKey   = "group-id-cus"
)

func init() {
	rds = redis.NewClient(
		&redis.Options{
			Addr:         "193.112.111.237:6378",
			DB:           1,
			MinIdleConns: 5,
			Password:     "chatai",
			PoolSize:     5,
			PoolTimeout:  30 * time.Second,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	)

	if err := rds.Ping().Err(); err != nil {
		log.Println("fail to connect redis, error msg: ", err)
	}

	log.Println("redis init completed")

}

type RM struct {
	mu sync.Mutex
}

func NewRM() *RM {
	return &RM{}
}

func (r *RM) formatKey(key string) string {
	return fmt.Sprintf("%s_%s", groupKey, key)
}

func (r *RM) Set(key string, b interface{}, expire time.Duration) error {
	return rds.Set(r.formatKey(key), b, expire).Err()
}

func (r *RM) Get(key string) (string, error) {
	result, err := rds.Get(r.formatKey(key)).Result()
	if err != nil {
		return result, err
	}

	if result == "" {
		return result, errors.New("null")
	}

	return result, nil
}

// GetAllData 当前市某个运动的所有场地地址列表, 只保留1个月, 1个月后重新更新, 主要是为了获取最新的场地数据
func (r *RM) GetAllData(key, cnKey, keyWord string) (string, error) {
	// key：shenzhenshi_bks
	r.mu.Lock()
	defer r.mu.Unlock()

	result, err := r.Get(key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return result, err
	}

	if result == "" {
		search, err := qqMapApi.NewTxMapApi(cnKey, keyWord).KeyWordSearch()
		if err != nil {
			return result, err
		}
		ld, err := r.mergeData(key)
		if err != nil {
			return result, err
		}
		if len(ld) > 0 {
			search = append(search, ld...)
		}
		b, err := json.Marshal(&search)
		if err != nil {
			return result, err
		}
		if err := r.Set(key, b, time.Second*time.Duration(2592000)); err != nil {
			return result, err
		}

		return string(b), nil
	}

	return result, nil
}

func (r *RM) mergeData(key string) ([]qqMapApi.SaveInRedis, error) {
	var dataList = make([]qqMapApi.SaveInRedis, 0)
	list, err := r.GetAddrList() // 遍历获取审核列表，找到对应id将其更新到指定key的数据中
	if err != nil {
		return dataList, err
	}

	dataList = make([]qqMapApi.SaveInRedis, 0, len(list))
	for _, data := range list {
		if data.SportKey == key && data.IsRecord {
			ad := qqMapApi.SaveInRedis{
				Id:     data.Id,
				Tags:   []string{data.Tags},
				Img:    defaultImg,
				Addr:   data.Addr,
				Lat:    data.Lat,
				Lng:    data.Lng,
				UserId: data.UserId,
				Title:  data.Tags,
			}
			dataList = append(dataList, ad)
		}
	}

	return dataList, nil
}

// Update 将审核通过的新的场地地添加到对应的列表中
func (r *RM) Update(key, id string) ([]qqMapApi.SaveInRedis, error) {
	var dataList []qqMapApi.SaveInRedis
	result, err := r.Get(key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return dataList, err
	}

	if err := json.Unmarshal([]byte(result), &dataList); err != nil {
		return dataList, err
	}

	list, err := r.GetAddrList() // 遍历获取审核列表，找到对应id将其更新到指定key的数据中
	if err != nil {
		return dataList, err
	}

	for _, data := range list {
		if data.Id == id {
			ad := qqMapApi.SaveInRedis{
				Id:     data.Id,
				Tags:   []string{data.Tags},
				Img:    "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/bk3.svg",
				Addr:   data.Addr,
				Lat:    data.Lat,
				Lng:    data.Lng,
				UserId: data.UserId,
				Title:  data.Tags,
			}
			dataList = append(dataList, ad)
			break
		}
	}

	b, err := json.Marshal(dataList)
	if err := json.Unmarshal([]byte(result), &dataList); err != nil {
		return dataList, err
	}

	if err := r.Set(key, b, time.Second*time.Duration(1209600)); err != nil {
		return dataList, err
	}

	if _, err := r.UpdateAddrList(id, true); err != nil {
		return dataList, err
	}

	return dataList, nil
}

// GetAddrList 所有用户添加的篮球场地址列表，不过期长期保存用户添加的篮球场地址
func (r *RM) GetAddrList() ([]*form.AddrListForm, error) {
	var dataList []*form.AddrListForm
	result, err := r.Get(AddrListKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return dataList, err
	}

	if errors.Is(err, redis.Nil) {
		dataList = make([]*form.AddrListForm, 0)
		b, err := json.Marshal(&dataList)
		if err != nil {
			return dataList, err
		}
		if err := r.Set(AddrListKey, b, 0); err != nil {
			return dataList, err
		}
		return dataList, nil
	}

	if err := json.Unmarshal([]byte(result), &dataList); err != nil {
		return dataList, err
	}

	return dataList, nil
}

// UserAddAddrReq 用户提交添加篮球场地址的请求
func (r *RM) UserAddAddrReq(data form.AddrListForm) error {
	var dataList = make([]form.AddrListForm, 0)
	result, err := r.Get(AddrListKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if result != "" {
		if err := json.Unmarshal([]byte(result), &dataList); err != nil {
			return err
		}
	}

	dataList = append(dataList, data)
	b, err := json.Marshal(&dataList)
	if err != nil {
		return err
	}

	if err := r.Set(AddrListKey, b, 0); err != nil {
		return err
	}

	return nil
}

// UpdateAddrList 更新审核列表
func (r *RM) UpdateAddrList(id string, status bool) ([]*form.AddrListForm, error) {
	list, err := r.GetAddrList() // 遍历获取审核列表，找到对应id将其更新到指定key的数据中
	if err != nil {
		return list, err
	}

	for _, v := range list {
		if v.Id == id {
			v.IsRecord = status
			v.IsShow = true
		}
	}

	b, err := json.Marshal(&list)
	if err != nil {
		return list, err
	}

	if err := r.Set(AddrListKey, b, 0); err != nil {
		return list, err
	}

	return list, nil
}

// SetWxOpenid 保存微信用户的openid
func (r *RM) SetWxOpenid(id string) error {
	var data = make([]form.WxOpenidList, 0)
	result, err := r.Get(WxOPenIdKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	if result == "" {
		d := form.WxOpenidList{
			Openid: id,
		}
		data = append(data, d)
		b, err := json.Marshal(&data)
		if err != nil {
			return err
		}

		if err := r.Set(WxOPenIdKey, b, 0); err != nil {
			return err
		}

		return nil
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return err
	}

	var isExist bool
	for _, v := range data {
		if v.Openid == id {
			isExist = true
		}
	}

	if !isExist {
		d := form.WxOpenidList{
			Openid: id,
		}
		data = append(data, d)

		b, err := json.Marshal(&data)
		if err != nil {
			return err
		}

		if err := r.Set(WxOPenIdKey, b, 0); err != nil {
			return err
		}
	}

	return nil
}

// GetWxOpenid 微信用户的openid
func (r *RM) GetWxOpenid(id string) error {
	var data = make([]form.WxOpenidList, 0)
	result, err := r.Get(WxOPenIdKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if result == "" {
		b, err := json.Marshal(&data)
		if err != nil {
			return err
		}

		if err := r.Set(WxOPenIdKey, b, 0); err != nil {
			return err
		}

		return errors.New("请先登陆微信")
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return err
	}

	var isFind bool
	for _, v := range data {
		if v.Openid == id {
			isFind = true
			break
		}
	}

	if !isFind {
		return errors.New("请先登陆微信")
	}

	return nil
}

// GetGroupOnline 获取在线人数
func (r *RM) GetGroupOnline(key string) (string, error) {
	gn := fmt.Sprintf("%s_%s", onlineKey, key)
	result, err := r.Get(gn)
	if err != nil && !errors.Is(err, redis.Nil) {
		return result, err
	}

	if result == "" {
		if err := r.Set(gn, 0, time.Second*time.Duration(7200)); err != nil {
			return result, err
		}

		return "0", nil
	}

	return result, nil
}

// VerifyWxUser 验证微信用户openid
func (r *RM) VerifyWxUser(hash string) (string, error) {
	var data = make([]form.WxOpenidList, 0)
	result, err := r.Get(WxOPenIdKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return result, err
	}

	if result == "" {
		return result, errors.New("用户不存在")
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return result, err
	}
	var oid string
	for _, v := range data {
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(v.Openid)); err == nil {
			oid = v.Openid
			break
		}
	}

	if oid == "" {
		return result, errors.New("用户不存在")
	}

	return oid, nil
}

// GetSportList 运动场地列表
func (r *RM) GetSportList() ([]form.SportList, error) {
	var data []form.SportList
	sports := `[
		{"name": "篮球场", "key": "bks", "checked": false, "icon": "🏀", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/main-bk.jpg"},
		{"name": "游泳馆", "key": "sws", "checked": false, "icon": "🏊", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/swim.png"},
		{"name": "羽毛球馆", "key": "bms", "checked": false, "icon": "🏸", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/badminton.png"},
		{"name": "足球场", "key": "fbs", "checked": false, "icon": "⚽", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/football.png"},
		{"name": "网球场", "key": "tns", "checked": false, "icon": "🎾", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/tennis.png"},
		{"name": "高尔夫球场", "key": "gos", "checked": false, "icon": "🏌️", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/golf.png"},
		{"name": "滑雪场", "key": "hxc", "checked": false, "icon": "⛷️", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/hxc.png"},
		{"name": "瑜伽馆", "key": "yjg", "checked": false, "icon": "🧘", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/yjg.png"},
		{"name": "跆拳道馆", "key": "tqd", "checked": false, "icon": "🥋", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/tqdg.png"},
		{"name": "健身房", "key": "gym", "checked": false, "icon": "🏋️‍♂️", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/gym.png"}
	]`
	if err := json.Unmarshal([]byte(sports), &data); err != nil {
		return data, err
	}

	return data, nil
}

// GetAllOnlineData 获取所有在线人数的key
func (r *RM) GetAllOnlineData() ([]form.GroupOnlineStatus, error) {
	var cursor uint64
	var matchingKeys []string
	var onlineStatus []form.GroupOnlineStatus
	matchPattern := "group_id_online_*" // 定义你的匹配模式

	for {
		// 使用 Scan 方法，传入游标、匹配模式和建议的单次扫描数量
		keys, nextCursor, err := rds.Scan(cursor, matchPattern, 10).Result()
		if err != nil {
			panic(err)
		}

		// 将本次扫描到的 key 追加到结果列表
		matchingKeys = append(matchingKeys, keys...)

		// 如果游标返回 0，说明迭代完成
		if nextCursor == 0 {
			break
		}

		// 更新游标以进行下一次迭代
		cursor = nextCursor
	}

	for _, v := range matchingKeys {
		online, err := rds.Get(v).Result()
		if err == nil {
			d := form.GroupOnlineStatus{
				GroupId:    v,
				OnlineUser: online,
			}
			onlineStatus = append(onlineStatus, d)
		}
	}

	return onlineStatus, nil
}

// GetJoinGroupUsers 获取每个组加入的人数
func (r *RM) GetJoinGroupUsers(key string) ([]form.JoinGroupUsers, error) {
	var data []form.JoinGroupUsers
	gn := fmt.Sprintf("%s_%s", joinGroupKey, key)
	result, err := r.Get(gn)
	if err != nil && !errors.Is(err, redis.Nil) {
		return data, err
	}

	if result == "" {
		var data = make([]form.JoinGroupUsers, 0)
		return data, nil
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return data, err
	}

	return data, nil
}

func (r *RM) JoinGroupUpdate(jd form.JoinGroupUsers) ([]form.JoinGroupUsers, error) {
	var data []form.JoinGroupUsers

	if jd.Oi == "1" {
		return r.exitGroup(jd)
	} else if jd.Oi == "2" {
		return r.UpdateJoinGroupUsers(jd)
	}

	return data, errors.New("invalid parameter")

}

// UpdateJoinGroupUsers 更新每个组加入的人数
func (r *RM) UpdateJoinGroupUsers(jd form.JoinGroupUsers) ([]form.JoinGroupUsers, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var data []form.JoinGroupUsers
	gn := fmt.Sprintf("%s_%s", joinGroupKey, jd.GroupId)
	result, err := r.Get(gn)
	if err != nil && !errors.Is(err, redis.Nil) {
		return data, err
	}

	//var d = form.JoinGroupUsers{
	//	GroupId: key,
	//	User:    uid,
	//	Img:     img,
	//}

	if result == "" {
		data = append(data, jd)
		b, err := json.Marshal(&data)
		if err != nil {
			return data, err
		}

		if err := r.Set(gn, b, time.Second*time.Duration(86400)); err != nil {
			return data, err
		}

		return data, nil
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return data, err
	}

	if r.checkUserIsJoinGroup(data, jd.GroupId, jd.User) {
		return data, cuerr.NewDuplicateError("已在该球局")
	}

	data = append(data, jd)
	b, err := json.Marshal(&data)
	if err != nil {
		return data, err
	}

	if err := r.Set(gn, b, time.Second*time.Duration(86400)); err != nil {
		return data, err
	}

	return data, nil
}

func (r *RM) exitGroup(jd form.JoinGroupUsers) ([]form.JoinGroupUsers, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var data []form.JoinGroupUsers
	gn := fmt.Sprintf("%s_%s", joinGroupKey, jd.GroupId)
	result, err := r.Get(gn)
	if err != nil && !errors.Is(err, redis.Nil) {
		return data, err
	}

	if result == "" {
		return make([]form.JoinGroupUsers, 0), nil
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return data, err
	}

	var fd []form.JoinGroupUsers
	for _, v := range data {
		if v.User != jd.User {
			fd = append(fd, v)
		}
	}

	if len(fd) == 0 {
		fd = make([]form.JoinGroupUsers, 0)
	}

	b, err := json.Marshal(&fd)
	if err != nil {
		return fd, err
	}

	if err := r.Set(gn, b, time.Second*time.Duration(86400)); err != nil {
		return fd, err
	}

	return fd, nil
}

func (r *RM) checkUserIsJoinGroup(data []form.JoinGroupUsers, gid, uid string) bool {
	for _, v := range data {
		if v.GroupId == gid && v.User == uid {
			return true
		}
	}
	return false
}
