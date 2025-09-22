package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/ingoxx/go-record/http/wx/pkg/config"
	"github.com/ingoxx/go-record/http/wx/pkg/distance"
	cuerr "github.com/ingoxx/go-record/http/wx/pkg/error"
	"github.com/ingoxx/go-record/http/wx/pkg/eva"
	"github.com/ingoxx/go-record/http/wx/pkg/form"
	"github.com/ingoxx/go-record/http/wx/pkg/mapApi"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math"
	"math/rand/v2"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	rds *redis.Client
)

func init() {
	rds = redis.NewClient(
		&redis.Options{
			Addr:         config.RedisAddr,
			DB:           1,
			MinIdleConns: 5,
			Password:     config.RedPwd,
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
	mu   sync.Mutex
	data chan []*form.SaveInRedis
}

func NewRM() *RM {
	return &RM{}
}

func (r *RM) formatKey(key string) string {
	return fmt.Sprintf("%s_%s", config.GroupKey, key)
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

func (r *RM) getAllData(key, cnKey, keyWord string) ([]*form.SaveInRedis, error) {
	data := mapApi.NewMapApi(cnKey, keyWord).Run()
	var result []*form.SaveInRedis
	for _, v := range data {
		if v.Err != nil {
			log.Printf("[ERROR] %s 接口请求失败, 失败信息：%v\n", v.Project, v.Err)
			continue
		}

		result = append(result, v.Data...)
	}

	if len(result) == 0 {
		return make([]*form.SaveInRedis, 0), nil
	}

	nr := r.uniqueByField(result)

	return nr, nil
}

// GetAllData 当前市某个运动的所有场地地址列表, 只保留半年月, 半年月后重新更新, 主要是为了获取最新的场地数据
func (r *RM) GetAllData(key, cnKey, keyWord, lat, lng, sportKey string) ([]*form.SaveInRedis, string, error) {
	// key：shenzhenshi_bks
	r.mu.Lock()
	defer r.mu.Unlock()

	var allData []*form.SaveInRedis

	result, err := r.Get(key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return allData, result, err
	}

	if result == "" {
		res, err := r.getAllData(key, cnKey, keyWord)
		if err != nil {
			return allData, result, err
		}

		// 合并审核列表中的数据
		ld, err := r.mergeData(key)
		if err != nil {
			return allData, result, err
		}

		if len(ld) > 0 {
			res = append(res, ld...)
		}

		// 如果存在的数据就只更新
		sd, err := r.updateLatestData(key, res)
		if err != nil {
			return allData, result, err
		}

		b, err := json.Marshal(&sd)
		if err != nil {
			return allData, result, err
		}

		if err := r.Set(key, b, 0); err != nil {
			return allData, result, err
		}

		gd, err := r.getVenueInfo(sd, lat, lng, sportKey, cnKey)
		if err != nil {
			return gd, string(b), err
		}

		return gd, string(b), nil

	}

	if err := json.Unmarshal([]byte(result), &allData); err != nil {
		return allData, result, err
	}

	nd, err := r.getVenueInfo(allData, lat, lng, sportKey, cnKey)
	if err != nil {
		return nd, result, err
	}

	return nd, result, nil
}

func (r *RM) getVenueInfo(data []*form.SaveInRedis, lat, lng, sportKey, cnKey string) ([]*form.SaveInRedis, error) {
	lat1, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return data, err
	}

	lng1, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		return data, err
	}

	for _, v := range data {
		//统计距离
		dis, err := distance.Distance(lat1, lng1, v.Lat, v.Lng)
		if err != nil {
			return data, err
		}
		v.Distance = fmt.Sprintf("%.1f", dis/1000)
		v.DisVal = math.Trunc(dis*10) / 10
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].DisVal < data[j].DisVal
	})

	var someData []*form.SaveInRedis
	if len(data) >= config.ShowNumber {
		someData = data[:config.ShowNumber]
	} else {
		someData = data
	}

	for _, v := range someData {
		// 统计当前在线人数
		//online, err := r.GetGroupOnline(v.Id)
		//if err != nil {
		//	return data, err
		//}
		//v.Online = online
		v.Online = "0"

		sd := &form.OnlineData{
			Id:       v.Id,
			Title:    v.Title,
			City:     cnKey,
			SportKey: sportKey,
			Online:   0,
		}
		if _, err := r.GetGroupOnline2(sd); err != nil {
			return data, err
		}

		//统计加入组队的人数
		users, err := r.GetJoinGroupUsers(v.Id)
		if err != nil {
			return data, err
		}
		if len(users) > 0 {
			v.JoinUserCount = len(users)
		}
		v.JoinUsers = users

		//统计对当前场地的评价
		board, err := r.GetMsgBoard(v.Id, sportKey)
		if err != nil {
			return data, err
		}

		v.UserReviews = board

		list, err := r.GetAllCheckList(v.Id)
		if err != nil {
			return data, err
		}

		v.VenueUpdateUsers = list
		if len(list) > 7 {
			v.VenueUpdateUsers = list[:7]
		}

		v.VenueUpdateUsersCount = len(list)

	}

	return someData, nil
}

func (r *RM) updateLatestData(key string, data []*form.SaveInRedis) ([]*form.SaveInRedis, error) {
	cl, err := r.GetAddrList() // 遍历获取审核列表，找到对应id将其更新到指定key的数据中
	if err != nil {
		return data, err
	}

	// IsRecord 字段必须是已经通过审核（为true）才能重新加回到场地列表中
	var isFind bool
	for _, v1 := range data {
		for _, v2 := range cl {
			if v1.Aid == v2.Aid && v2.SportKey == key {
				if v2.IsRecord && v2.UpdateType == "2" {
					v1.Images = append(v1.Images, v2.Img)
					v1.Img = v2.Img
				}
				isFind = true
			}
		}
	}

	if !isFind {
		for _, v2 := range cl {
			if v2.SportKey == key && v2.IsRecord && v2.UpdateType == "2" {
				ad := &form.SaveInRedis{
					Id:     v2.Id,
					Tags:   []string{v2.Tags},
					Img:    v2.Img,
					Addr:   v2.Addr,
					Lat:    v2.Lat,
					Lng:    v2.Lng,
					UserId: v2.UserId,
					Title:  v2.Tags,
					Aid:    v2.Id,
				}
				ad.Images = append(ad.Images, v2.Img)
				data = append(data, ad)
			}
		}
	}

	return data, nil
}

func (r *RM) mergeData(key string) ([]*form.SaveInRedis, error) {
	var dataList = make([]*form.SaveInRedis, 0)
	list, err := r.GetAddrList() // 遍历获取审核列表，找到对应id将其更新到指定key的数据中
	if err != nil {
		return dataList, err
	}

	dataList = make([]*form.SaveInRedis, 0, len(list))
	// IsRecord 字段必须是已经通过审核（为true）才能重新加回到场地列表中
	for _, data := range list {
		if data.SportKey == key && data.IsRecord && data.UpdateType == "1" {
			ad := &form.SaveInRedis{
				Id:     data.Id,
				Tags:   []string{data.Tags},
				Img:    data.Img,
				Addr:   data.Addr,
				Lat:    data.Lat,
				Lng:    data.Lng,
				UserId: data.UserId,
				Title:  data.Tags,
				Aid:    data.Id,
			}
			ad.Images = append(ad.Images, data.Img)
			dataList = append(dataList, ad)
		}
	}

	return dataList, nil
}

// Update 将审核通过的新的场地地添加到对应的列表中
func (r *RM) Update(key, id, ut string) ([]form.SaveInRedis, error) {
	var dataList []form.SaveInRedis
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

	// 更新已经存在的场地
	if ut == "2" {
		if _, err := r.UpdateAddrList(id, true); err != nil {
			return dataList, err
		}

		return dataList, nil
	}

	// 添加新的场地
	for _, data := range list {
		if data.Id == id {
			ad := form.SaveInRedis{
				Id:     data.Id,
				Tags:   []string{data.Tags},
				Img:    data.Img,
				Addr:   data.Addr,
				Lat:    data.Lat,
				Lng:    data.Lng,
				UserId: data.UserId,
				Title:  data.Tags,
			}
			ad.Images = append(ad.Images, data.Img)
			dataList = append(dataList, ad)
			break
		}
	}

	b, err := json.Marshal(&dataList)
	if err != nil {
		return dataList, err
	}
	//if err := json.Unmarshal([]byte(result), &dataList); err != nil {
	//	return dataList, err
	//}

	if err := r.Set(key, b, 0); err != nil {
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
	result, err := r.Get(config.AddrListKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return dataList, err
	}

	if errors.Is(err, redis.Nil) {
		dataList = make([]*form.AddrListForm, 0)
		b, err := json.Marshal(&dataList)
		if err != nil {
			return dataList, err
		}
		if err := r.Set(config.AddrListKey, b, 0); err != nil {
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
func (r *RM) UserAddAddrReq(data *form.AddrListForm) error {
	var dataList = make([]*form.AddrListForm, 0)
	result, err := r.Get(config.AddrListKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if result != "" {
		if err := json.Unmarshal([]byte(result), &dataList); err != nil {
			return err
		}
	}

	if data.UpdateType == "1" {
		data.Content = fmt.Sprintf("用户在%s添加了场地图片", time.Now().Format("2006-01-02 15:04:05"))
	} else if data.UpdateType == "2" {
		data.Content = fmt.Sprintf("用户在%s更新了场地图片", time.Now().Format("2006-01-02 15:04:05"))
	}

	dataList = append(dataList, data)
	b, err := json.Marshal(&dataList)
	if err != nil {
		return err
	}

	if err := r.Set(config.AddrListKey, b, 0); err != nil {
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
		if v.Id == id && !v.IsRecord { // 更新场地信息的时候，可能会存在多个相同场地的id，必须是v.IsRecord为false的才能更新
			v.IsRecord = status
			v.IsShow = true
		}
	}

	b, err := json.Marshal(&list)
	if err != nil {
		return list, err
	}

	if err := r.Set(config.AddrListKey, b, 0); err != nil {
		return list, err
	}

	return list, nil
}

// SetWxOpenid 保存微信用户的openid
func (r *RM) SetWxOpenid(wo *form.WxOpenidList) (*form.WxOpenidList, error) {
	var data = make([]*form.WxOpenidList, 0)
	var fd = new(form.WxOpenidList)

	result, err := r.Get(config.WxOPenIdKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return fd, err
	}

	// 用户不存在就添加
	if result == "" {
		if wo.NickName == "" {
			wo.NickName = eva.NewSportType("bks").RandomNickname()
		}
		data = append(data, wo)
		b, err := json.Marshal(&data)
		if err != nil {
			return fd, err
		}

		if err := r.Set(config.WxOPenIdKey, b, 0); err != nil {
			return fd, err
		}

		return fd, nil
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return fd, err
	}

	// 查找用户是否存在
	var isExist bool
	for _, v := range data {
		if v.Openid == wo.Openid {
			if wo.NickName != "" {
				v.NickName = wo.NickName
			}
			if wo.Img != "" {
				v.Img = wo.Img
			}

			fd = v
			isExist = true
			break
		}
	}

	// 用户不存在就添加
	if !isExist {
		if wo.NickName == "" {
			wo.NickName = eva.NewSportType("bks").RandomNickname()
		}
		wo.Img = r.generateRandomImg()
		data = append(data, wo)

		b, err := json.Marshal(&data)
		if err != nil {
			return fd, err
		}

		if err := r.Set(config.WxOPenIdKey, b, 0); err != nil {
			return fd, err
		}

		return wo, nil
	}

	// 更新用户信息
	b, err := json.Marshal(&data)
	if err != nil {
		return fd, err
	}

	if err := r.Set(config.WxOPenIdKey, b, 0); err != nil {
		return fd, err
	}

	return fd, nil
}

func (r *RM) UpdateWxUser(id, city string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var data = make([]*form.WxOpenidList, 0)

	result, err := r.Get(config.WxOPenIdKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Printf("[ERROR] get wx openid error, esg: '%s'\n", err.Error())
		return
	}

	if result == "" {
		log.Println("[INFO] get wx openid list empty")
		return
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		log.Printf("[ERROR] unmarshal openid list data error, esg: '%s'\n", err.Error())
		return
	}

	for _, v := range data {
		if v.Openid == id {
			v.City = city
			break
		}
	}

	b, err := json.Marshal(&data)
	if err != nil {
		log.Printf("[ERROR] marshal openid list data error, esg: '%s'\n", err.Error())
		return
	}

	if err := r.Set(config.WxOPenIdKey, b, 0); err != nil {
		log.Printf("[ERROR] set openid list data error, esg: '%s'\n", err.Error())
		return
	}

}

// GetWxOpenid 微信用户的openid
func (r *RM) GetWxOpenid(id string) error {
	var data = make([]*form.WxOpenidList, 0)
	result, err := r.Get(config.WxOPenIdKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if result == "" {
		b, err := json.Marshal(&data)
		if err != nil {
			return err
		}

		if err := r.Set(config.WxOPenIdKey, b, 0); err != nil {
			return err
		}

		return errors.New("请先登陆微信1")
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
		return errors.New("请先登陆微信2")
	}

	return nil
}

// GetGroupOnline 获取在线人数
func (r *RM) GetGroupOnline(key string) (string, error) {
	gn := fmt.Sprintf("%s_%s", config.OnlineKey, key)
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

// GetGroupOnline2 获取在线人数
func (r *RM) GetGroupOnline2(data *form.OnlineData) (*form.OnlineData, error) {
	gn := fmt.Sprintf("%s_%s", config.OnlineKey, data.Id)
	result, err := r.Get(gn)
	if err != nil && !errors.Is(err, redis.Nil) {
		return data, err
	}

	if result == "" {
		b, err := json.Marshal(&data)
		if err != nil {
			return data, err
		}

		if err := r.Set(gn, b, time.Second*time.Duration(7200)); err != nil {
			return data, err
		}

		return data, nil
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return data, err
	}

	return data, nil
}

// UpdateGroupOnline 更新在线人数
func (r *RM) UpdateGroupOnline(data *form.OnlineData) error {
	var fd *form.OnlineData
	gn := fmt.Sprintf("%s_%s", config.OnlineKey, data.Id)
	result, err := r.Get(gn)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if result == "" {
		b, err := json.Marshal(&data)
		if err != nil {
			return err
		}
		if err := r.Set(gn, b, time.Second*time.Duration(7200)); err != nil {
			return err
		}

		return nil
	}

	if err := json.Unmarshal([]byte(result), &fd); err != nil {
		return err
	}

	fd.Online = data.Online
	b, err := json.Marshal(&fd)
	if err != nil {
		return err
	}
	if err := r.Set(gn, b, time.Second*time.Duration(7200)); err != nil {
		return err
	}

	return nil
}

// VerifyWxUser 验证微信用户openid
func (r *RM) VerifyWxUser(hash string) (string, error) {
	var data = make([]form.WxOpenidList, 0)
	result, err := r.Get(config.WxOPenIdKey)
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
	//sports := `[
	//	{"name": "篮球场", "key": "bks", "checked": false, "icon": "🏀", "img": "https://ai.anythingai.online/static/profile3/main-bk.jpg"},
	//	{"name": "游泳馆", "key": "sws", "checked": false, "icon": "🏊", "img": "https://ai.anythingai.online/static/profile3/swim.png"},
	//	{"name": "羽毛球馆", "key": "bms", "checked": false, "icon": "🏸", "img": "https://ai.anythingai.online/static/profile3/badminton.png"},
	//	{"name": "足球场", "key": "fbs", "checked": false, "icon": "⚽", "img": "https://ai.anythingai.online/static/profile3/football.png"}
	//]`
	sports := `[
		{"title": "篮球场", "name": "🏀篮球场", "key": "bks", "checked": false, "icon": "🏀", "img": "https://ai.anythingai.online/static/profile3/main-bk.jpg", "sport_img": "https://ai.anythingai.online/static/profile3/bks-6.svg"},
		{"title": "攀岩馆", "name": "🧗攀岩馆", "key": "rcg", "checked": false, "icon": "🧗", "img": "https://ai.anythingai.online/static/profile3/rcg.png", "sport_img": "https://ai.anythingai.online/static/profile3/rcg-5.svg"},
		{"title": "游泳馆", "name": "🏊游泳馆", "key": "sws", "checked": false, "icon": "🏊", "img": "https://ai.anythingai.online/static/profile3/swim.png", "sport_img": "https://ai.anythingai.online/static/profile3/swim-6.svg"},
		{"title": "羽毛球馆", "name": "🏸羽毛球馆", "key": "bms", "checked": false, "icon": "🏸", "img": "https://ai.anythingai.online/static/profile3/badminton.png", "sport_img": "https://ai.anythingai.online/static/profile3/bms-6.svg"},
		{"title": "足球场", "name": "⚽足球场", "key": "fbs", "checked": false, "icon": "⚽", "img": "https://ai.anythingai.online/static/profile3/football.png", "sport_img": "https://ai.anythingai.online/static/profile3/fbs-6.svg"}
	]`
	if err := json.Unmarshal([]byte(sports), &data); err != nil {
		return data, err
	}

	return data, nil
}

// GetAllOnlineData 获取所有在线人数的key
func (r *RM) GetAllOnlineData(key string) ([]form.GroupOnlineStatus, error) {
	var cursor uint64
	var matchingKeys []string
	var onlineStatus []form.GroupOnlineStatus
	matchPattern := "group_id_online_*" // 定义你的匹配模式

	for {
		// 使用 Scan 方法，传入游标、匹配模式和建议的单次扫描数量
		keys, nextCursor, err := rds.Scan(cursor, matchPattern, 10).Result()
		if err != nil {
			return onlineStatus, err
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
		//gid := strings.ReplaceAll(v, "group_id_online_", "")

		//name, err := r.getVenueName(key, gid)
		//if err != nil {
		//	return onlineStatus, err
		//}
		if err == nil {
			d := form.GroupOnlineStatus{
				GroupId:    v,
				OnlineUser: online,
				//VenueName:  name,
			}
			onlineStatus = append(onlineStatus, d)
		}
	}

	return onlineStatus, nil
}

// GetAllOnlineData2 获取所有在线人数的key
func (r *RM) GetAllOnlineData2(key string) ([]*form.OnlineData, error) {
	var cursor uint64
	var matchingKeys []string

	var onlineStatus []*form.OnlineData
	matchPattern := "group_id_online_*" // 定义你的匹配模式

	for {
		// 使用 Scan 方法，传入游标、匹配模式和建议的单次扫描数量
		keys, nextCursor, err := rds.Scan(cursor, matchPattern, 10).Result()
		if err != nil {
			return onlineStatus, err
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
		var od *form.OnlineData
		online, err := rds.Get(v).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(online), &od); err != nil {
				return onlineStatus, err
			}
			onlineStatus = append(onlineStatus, od)
		}
	}

	return onlineStatus, nil
}

// GetAllOnlineData3 获取所有在线人数的key
func (r *RM) GetAllOnlineData3(ids []string) ([]*form.OnlineData, error) {
	var (
		cursor       uint64
		onlineStatus []*form.OnlineData
		matchPattern = "group_id_online_*"
	)

	// 将 id slice 转换成 map，方便快速查找
	idSet := make(map[string]struct{}, len(ids))
	for _, v := range ids {
		idSet[v] = struct{}{}
	}

	for {
		keys, nextCursor, err := rds.Scan(cursor, matchPattern, 20).Result() // 每次扫描 20 个
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			// 提取 key 中的 id（假设 key 格式是 group_id_online_<id>）
			parts := strings.Split(key, "_")
			if len(parts) < 4 {
				continue
			}
			keyID := parts[3]

			// 判断 id 是否在目标列表中
			if _, ok := idSet[keyID]; !ok {
				continue
			}

			// 获取并反序列化数据
			val, err := rds.Get(key).Result()
			if err != nil {
				// key 不存在或其他错误，跳过
				continue
			}

			var od form.OnlineData
			if err := json.Unmarshal([]byte(val), &od); err != nil {
				return nil, err
			}
			onlineStatus = append(onlineStatus, &od)
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return onlineStatus, nil
}

// GetJoinGroupUsers 获取每个组加入的人数
func (r *RM) GetJoinGroupUsers(key string) ([]*form.JoinGroupUsers, error) {
	var data []*form.JoinGroupUsers
	gn := fmt.Sprintf("%s_%s", config.JoinGroupKey, key)
	result, err := r.Get(gn)
	if err != nil && !errors.Is(err, redis.Nil) {
		return data, err
	}

	if result == "" {
		return make([]*form.JoinGroupUsers, 0), nil
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return data, err
	}

	return data, nil
}

func (r *RM) JoinGroupUpdate(jd *form.JoinGroupUsers) ([]*form.JoinGroupUsers, error) {
	var data []*form.JoinGroupUsers

	if jd.Oi == "1" {
		return r.exitGroup(jd)
	} else if jd.Oi == "2" {
		return r.UpdateJoinGroupUsers(jd)
	}

	return data, errors.New("invalid parameter")

}

// UpdateJoinGroupUsers 更新每个组加入的人数
func (r *RM) UpdateJoinGroupUsers(jd *form.JoinGroupUsers) ([]*form.JoinGroupUsers, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var data []*form.JoinGroupUsers
	gn := fmt.Sprintf("%s_%s", config.JoinGroupKey, jd.GroupId)
	result, err := r.Get(gn)
	if err != nil && !errors.Is(err, redis.Nil) {
		return data, err
	}

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

// exitGroup 退出组局
func (r *RM) exitGroup(jd *form.JoinGroupUsers) ([]*form.JoinGroupUsers, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var data []*form.JoinGroupUsers
	gn := fmt.Sprintf("%s_%s", config.JoinGroupKey, jd.GroupId)
	result, err := r.Get(gn)
	if err != nil && !errors.Is(err, redis.Nil) {
		return data, err
	}

	if result == "" {
		return make([]*form.JoinGroupUsers, 0), nil
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return data, err
	}

	var fd []*form.JoinGroupUsers
	for _, v := range data {
		if v.User != jd.User {
			fd = append(fd, v)
		}
	}

	if len(fd) == 0 {
		fd = make([]*form.JoinGroupUsers, 0)
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

func (r *RM) checkUserIsJoinGroup(data []*form.JoinGroupUsers, gid, uid string) bool {
	for _, v := range data {
		if v.GroupId == gid && v.User == uid {
			return true
		}
	}
	return false
}

// GetMsgBoard 获取某个场地的所有评价
func (r *RM) GetMsgBoard(gid, sportKey string) ([]*form.MsgBoard, error) {
	var data []*form.MsgBoard
	gn := fmt.Sprintf("%s_%s", config.EvaKey, gid)
	result, err := r.Get(gn)
	if err != nil && !errors.Is(err, redis.Nil) {
		return data, err
	}

	if result == "" {
		board, err := eva.NewSportType(sportKey).DefaultEvaBoard()
		if err != nil {
			return data, err
		}

		nd := r.updateImg(board, gid)
		data = append(data, nd...)
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

	return data, nil
}

// UpdateMsgBoard 用户提交对谋个场地的评价
func (r *RM) UpdateMsgBoard(mb *form.MsgBoard, sportKey string) ([]*form.MsgBoard, error) {
	var data []*form.MsgBoard
	gn := fmt.Sprintf("%s_%s", config.EvaKey, mb.GroupId)
	result, err := r.Get(gn)
	if err != nil && !errors.Is(err, redis.Nil) {
		return data, err
	}

	if result == "" {
		board, err := eva.NewSportType(sportKey).DefaultEvaBoard()
		if err != nil {
			return data, err
		}
		nd := r.updateImg(board, mb.GroupId)
		data = append(data, mb)
		data = append(data, nd...)
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

	if len(mb.LikeUsers) == 0 {
		mb.LikeUsers = make([]string, 0)
	}

	data = append(data, mb)
	b, err := json.Marshal(&data)
	if err != nil {
		return data, err
	}

	if err := r.Set(gn, b, 0); err != nil {
		return data, err
	}

	return data, nil
}

func (r *RM) updateImg(data []*form.MsgBoard, gid string) []*form.MsgBoard {
	m1 := 1001
	m2 := 2904
	for _, v := range data {
		v.Img = fmt.Sprintf("%s/%d.png", config.ImgUrl, rand.IntN(m2-m1+1)+m1)
		v.GroupId = gid
		if len(v.LikeUsers) == 0 {
			v.LikeUsers = make([]string, 0)
		}
	}

	return data
}

func (r *RM) generateRandomImg() string {
	m1 := 1001
	m2 := 2904
	return fmt.Sprintf("%s/%d.png", config.ImgUrl, rand.IntN(m2-m1+1)+m1)
}

// UserLikedReviews 用户点赞留言
func (r *RM) UserLikedReviews(mb *form.MsgBoard, sportKey string) ([]*form.MsgBoard, error) {
	var data []*form.MsgBoard
	gn := fmt.Sprintf("%s_%s", config.EvaKey, mb.GroupId)
	result, err := r.Get(gn)
	if err != nil && !errors.Is(err, redis.Nil) {
		return data, err
	}

	if result == "" {
		board, err := eva.NewSportType(sportKey).DefaultEvaBoard()
		if err != nil {
			return data, err
		}

		nd := r.updateImg(board, mb.GroupId)
		ndd := r.updateLike(nd, mb)
		data = append(data, ndd...)
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

	nd := r.updateLike(data, mb)

	b, err := json.Marshal(&nd)
	if err != nil {
		return data, err
	}

	if err := r.Set(gn, b, 0); err != nil {
		return data, err
	}

	return nd, nil
}

func (r *RM) updateLike(data []*form.MsgBoard, mb *form.MsgBoard) []*form.MsgBoard {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, v := range data {
		if v.EvaluateId == mb.EvaluateId {
			if mb.IsLike {
				v.Like += 1
				v.LikeUsers = append(v.LikeUsers, mb.User)
			} else {
				if v.Like > 0 {
					v.Like -= 1

				}
				v.LikeUsers = r.updateLikeUsers(v.LikeUsers, mb.User)
			}
			break
		}
	}

	return data
}

func (r *RM) updateLikeUsers(u []string, uid string) []string {
	for i, v := range u {
		if v == uid {
			return append(u[:i], u[i+1:]...)
		}
	}

	return u
}

// UpdateVenueInfo 更新运动场地
func (r *RM) UpdateVenueInfo(dt *form.AddrListForm) ([]*form.SaveInRedis, error) {
	var data []*form.SaveInRedis
	result, err := r.Get(dt.SportKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return data, err
	}

	if result == "" {
		return make([]*form.SaveInRedis, 0), nil
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return data, err
	}

	for _, v := range data {
		if v.Id == dt.Id {
			v.Img = dt.Img
		}
	}

	b, err := json.Marshal(&data)
	if err != nil {
		return data, err
	}

	if err := r.Set(dt.SportKey, b, 0); err != nil {
		return data, err
	}

	return data, nil
}

// GetAllWxUsers 获取所有微信用户
func (r *RM) GetAllWxUsers() ([]*form.WxOpenidList, error) {
	var data []*form.WxOpenidList
	result, err := r.Get(config.WxOPenIdKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return data, err
	}

	if result == "" {
		return make([]*form.WxOpenidList, 0), nil
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return data, err
	}

	return data, nil
}

// GetAllCheckList 将相同的场地放在一起
func (r *RM) GetAllCheckList(vid string) ([]*form.AddrListForm, error) {
	var fd = make(map[string][]*form.AddrListForm)

	cl, err := r.GetAddrList() // 遍历获取审核列表，找到对应id将其更新到指定key的数据中
	if err != nil {
		return fd[vid], err
	}

	for _, v := range cl {
		if v.IsRecord {
			i, ok := fd[v.Id]
			if !ok {
				fd[v.Id] = append(fd[v.Id], v)
				continue
			}
			i = append(i, v)
			fd[v.Id] = i
		}
	}

	if len(fd[vid]) > 1 {
		sort.Slice(fd[vid], func(i, j int) bool {
			return r.getTimestamp(fd[vid][i].Time) < r.getTimestamp(fd[vid][j].Time)
		})
	}

	return fd[vid], err
}

func (r *RM) getTimestamp(timeStr string) int64 {
	layout := "2006-01-02 15:04:05"

	t, err := time.ParseInLocation(layout, timeStr, time.Local)
	if err != nil {
		return 0
	}
	return t.UnixNano()
}

func (r *RM) randomPick(s []*form.MsgBoard, n int) []*form.MsgBoard {
	if n > len(s) {
		n = len(s)
	}

	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})

	return s[:n]
}

// 随机获取切片中的随机数量元素
func (r *RM) randomSubset(s []*form.MsgBoard) []*form.MsgBoard {
	n := rand.IntN(len(s)) + 1
	return r.randomPick(s, n)
}

// GetVenueImg 获取场地图片
func (r *RM) GetVenueImg(key, aid, city string) error {
	var allData []*form.SaveInRedis

	result, err := r.Get(key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if result == "" {
		return err
	}

	if err := json.Unmarshal([]byte(result), &allData); err != nil {
		return err
	}

	for _, v := range allData {
		if v.Id == aid && v.Img == "" {
			searchImg, err := mapApi.NewMapApi(city, v.Title).GetGdSinglePlaceSearch()
			if err != nil {
				return err
			}

			if len(searchImg) == 0 {
				return errors.New("未获取到图片")
			}

			if len(searchImg) > 0 {
				v.Img = searchImg[0]
				v.Images = append(v.Images, searchImg...)
			}
			break
		}
	}

	b, err := json.Marshal(&allData)
	if err != nil {
		return err
	}

	if err := r.Set(key, b, 0); err != nil {
		return err
	}

	return nil
}

func (r *RM) FilterVenueData() ([]*form.FilterField, error) {
	var data []*form.FilterField
	var jd = `[
		{"id": 1, "name": "距离最近", "type": 1000 },
		{"id": 2, "name": "组队人数", "type": 1000},
		{"id": 3, "name": "评价数量", "type": 1000},
		{"id": 5, "name": "新发布", "type": 2000},
		{"id": 4, "name": "价格", "type": 2000},
		{"id": 6, "name": "已删除", "type": 2000},
		{"id": 7, "name": "生效中", "type": 2000}
	]`

	if err := json.Unmarshal([]byte(jd), &data); err != nil {
		return data, err
	}

	return data, nil
}

func (r *RM) uniqueByField(data []*form.SaveInRedis) []*form.SaveInRedis {
	m := make(map[string]*form.SaveInRedis)

	for _, p := range data {
		if existing, ok := m[p.Title]; ok {
			// 判断优先级：优先保留 Gender 非空的
			if existing.Img == "" && p.Img != "" {
				m[p.Title] = p // 替换为 Gender 不为空的
			}
			// 如果都为空或都有值，保留第一次出现的即可，不做替换
		} else {
			m[p.Title] = p
		}
	}

	// 转换回切片
	result := make([]*form.SaveInRedis, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

// GetWxBtnText 一些隐藏按钮
func (r *RM) GetWxBtnText() ([]*form.WxBtnText, error) {
	data := `[
		{"id": 1, "name": "获取更多场地图片"},
		{"id": 2, "name": "发布付费陪练"},
		{"id": 3, "name": "场地"},
		{"id": 4, "name": "陪练"}
	]`
	var fd []*form.WxBtnText

	if err := json.Unmarshal([]byte(data), &fd); err != nil {
		return fd, err
	}

	return fd, nil
}

// AddPublish 发布任务
func (r *RM) AddPublish(data *form.PublishData) ([]*form.PublishData, error) {
	var all []*form.PublishData

	data.Id = uuid.NewString()
	week, err := r.getWeek(data.Date)
	if err != nil {
		return all, err
	}

	data.Date = fmt.Sprintf("%s %s", data.Date, week)
	data.PublishDate = time.Now().Format("01-02")
	data.Time = time.Now().Format("2006-01-02 15:04:05")

	// 运动+城市 Hash
	key1 := fmt.Sprintf("publish_%s_%s", data.CityPy, data.SportKey)
	// 用户 Hash
	key2 := fmt.Sprintf("user_publish_%s", data.UserId)

	b, err := json.Marshal(data)
	if err != nil {
		return all, err
	}

	if err := rds.HSet(key1, data.Id, b).Err(); err != nil {
		return all, err
	}
	if err := rds.HSet(key2, data.Id, b).Err(); err != nil {
		return all, err
	}

	sport, err := r.GetTasksByCityAndSport(data.SportKey, data.CityPy)
	if err != nil {
		return sport, err
	}

	return sport, nil
}

func (r *RM) delPublishHistory(userId, tid string) ([]*form.PublishData, error) {
	key := fmt.Sprintf("user_publish_%s", userId)
	values, err := rds.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}

	var result []*form.PublishData
	var uc = make([]*form.UserRoomID, 0)
	for _, v := range values {
		var p *form.PublishData
		if err := json.Unmarshal([]byte(v), &p); err != nil {
			return result, err
		}

		if p.Id == tid {
			p.IsDel = true
			p.UserCount = uc
			p.OnlineNum = 0
			b, err := json.Marshal(&p)
			if err != nil {
				return result, err
			}
			if err := rds.HSet(key, tid, b).Err(); err != nil {
				return result, err
			}
		}

		result = append(result, p)
	}

	if len(result) == 0 {
		return make([]*form.PublishData, 0), nil
	}

	return result, nil
}

// GetUserPublishData 查询用户自己发布的所有任务
func (r *RM) GetUserPublishData(userId string) ([]*form.PublishData, error) {
	key := fmt.Sprintf("user_publish_%s", userId)
	values, err := rds.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}

	var result []*form.PublishData
	for _, v := range values {
		var p *form.PublishData
		if err := json.Unmarshal([]byte(v), &p); err == nil {
			if !p.IsDel {
				result = append(result, p)
			}
		}
	}

	if len(result) == 0 {
		return make([]*form.PublishData, 0), nil
	}

	return result, nil
}

// UpdatePublish 更新任务（同时更新两个 Hash）
func (r *RM) UpdatePublish(data *form.PublishData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 更新运动+城市
	key1 := fmt.Sprintf("publish_%s_%s", data.CityPy, data.SportKey)
	// 更新用户 Hash
	key2 := fmt.Sprintf("user_publish_%s", data.UserId)

	if err := rds.HSet(key1, data.Id, b).Err(); err != nil {
		return err
	}
	if err := rds.HSet(key2, data.Id, b).Err(); err != nil {
		return err
	}
	return nil
}

// UpdateSinglePublishData 更新某个任务，删除或者标记已完成
func (r *RM) UpdateSinglePublishData(ms *form.MissionStatus, uid string) ([]*form.PublishData, error) {
	var fd []*form.PublishData
	hashKey := fmt.Sprintf("publish_%s_%s", ms.CityPy, ms.SportKey)

	// 1. 获取 JSON
	val, err := rds.HGet(hashKey, ms.Id).Result()
	if errors.Is(err, redis.Nil) {
		return fd, fmt.Errorf("task %s not found", ms.Id)
	} else if err != nil {
		return fd, err
	}

	var uc = make([]*form.UserRoomID, 0)
	// 2. 反序列化
	var data *form.PublishData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return fd, err
	}

	if ms.Status == 1 {
		data.Finish = true
	} else if ms.Status == 2 {
		data.IsDel = true
		data.UserCount = uc
		data.OnlineNum = 0
	} else if ms.Status == 3 {
		data.IsDel = false
	} else if ms.Status == 4 {
		data.Finish = false
	}

	// 4. 序列化并写回
	newVal, err := json.Marshal(data)
	if err != nil {
		return fd, err
	}

	if err := rds.HSet(hashKey, ms.Id, newVal).Err(); err != nil {
		return fd, err
	}

	if _, err := r.delPublishHistory(ms.UserId, ms.Id); err != nil {
		return fd, err
	}

	if _, err := r.delPublishTid(ms.Id); err != nil {
		return fd, err
	}

	if uid == config.Admin { // 如果删除发布任务的是管理员就返回管理员能看到的数据
		fd, err = r.GetAllPublishData(ms.SportKey)
		if err != nil {
			return fd, err
		}
	} else { // 如果是用户自己删除的发布任务就返回普通用户才能看到的数据
		fd, err = r.GetTasksByCityAndSport(ms.SportKey, ms.CityPy)
		if err != nil {
			return fd, err
		}
	}

	return fd, nil
}

// GetAllPublishData 只有管理员才能获取不限制城市指定运动类型的所有数据
func (r *RM) GetAllPublishData(sportKey string) ([]*form.PublishData, error) {
	pattern := fmt.Sprintf("publish_*_%s", sportKey)

	var tasks []*form.PublishData

	// 用 SCAN 避免 KEYS 阻塞
	iter := rds.Scan(0, pattern, 0).Iterator()
	for iter.Next() {
		hashKey := iter.Val()

		// 获取该城市下所有任务
		values, err := rds.HGetAll(hashKey).Result()
		if err != nil {
			return nil, err
		}

		for _, v := range values {
			var pd *form.PublishData
			if err := json.Unmarshal([]byte(v), &pd); err == nil {
				tasks = append(tasks, pd)
				tid, err := r.GetPublishTid(pd.Id)
				if err != nil {
					log.Printf("获取tid为：%s接入用户人数时报错: %s", pd.Id, err.Error())
				}
				pd.UserCount = append(pd.UserCount, tid...)
				pd.OnlineNum = len(pd.UserCount)
			}
		}
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return make([]*form.PublishData, 0), nil
	}

	return tasks, nil
}

// GetTasksByCityAndSport 获取指定城市+运动类型的所有任务-当前城市的用户选择了运动类型就可以查看所有悬赏任务
func (r *RM) GetTasksByCityAndSport(sportKey, city string) ([]*form.PublishData, error) {
	hashKey := fmt.Sprintf("publish_%s_%s", city, sportKey) // 例如 publish_shenzhenshi_bks

	allValues, err := rds.HGetAll(hashKey).Result()
	if err != nil {
		return nil, err
	}

	tasks := make([]*form.PublishData, 0, len(allValues))
	for _, v := range allValues {
		var pd *form.PublishData
		if err := json.Unmarshal([]byte(v), &pd); err == nil {
			if !pd.IsDel {
				tasks = append(tasks, pd)
				tid, err := r.GetPublishTid(pd.Id)
				if err != nil {
					log.Printf("获取tid为：%s接入用户人数时报错: %s", pd.Id, err.Error())
				}
				pd.UserCount = append(pd.UserCount, tid...)
				pd.OnlineNum = len(pd.UserCount)
			}
		}
	}

	if len(tasks) == 0 {
		return make([]*form.PublishData, 0), nil
	}

	return tasks, nil
}

func (r *RM) getWeek(dateStr string) (string, error) {
	// 解析日期字符串为 time.Time 对象
	layout := "2006-01-02 15:04:05" // 时间格式必须使用这个参考值
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return dateStr, err
	}

	weekday := t.Weekday()

	// 如果需要中文格式，可以用一个映射
	weekMap := map[time.Weekday]string{
		time.Sunday:    "星期日",
		time.Monday:    "星期一",
		time.Tuesday:   "星期二",
		time.Wednesday: "星期三",
		time.Thursday:  "星期四",
		time.Friday:    "星期五",
		time.Saturday:  "星期六",
	}
	return weekMap[weekday], nil
}

// GenerateId 用户点击沟通生成属于该发布id下的唯一room id
func (r *RM) GenerateId(data *form.UserRoomID) (*form.UserRoomID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("publish_tid_%s", data.Tid) // 例如 publish_tid_aaa-bbb-ccc
	var fds []*form.UserRoomID
	var nd *form.UserRoomID
	var isFind bool
	result, err := rds.HGet(key, data.Tid).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return data, err
	}

	if result == "" {
		data.Rid = uuid.NewString()
		fds = append(fds, data)
	} else {
		if err := json.Unmarshal([]byte(result), &fds); err != nil {
			return data, err
		}
		for _, v := range fds {
			if v.UserId == data.UserId {
				nd = v
				isFind = true
			}
		}

		if !isFind {
			data.Rid = uuid.NewString()
			fds = append(fds, data)
		}
	}

	b, err := json.Marshal(&fds)
	if err != nil {
		return data, err
	}

	if err := rds.HSet(key, data.Tid, b).Err(); err != nil {
		return data, err
	}

	if isFind {
		return nd, nil
	}

	return data, nil
}

// GetPublishTid 返回某个发布id下的所有接入的room id
func (r *RM) GetPublishTid(tid string) ([]*form.UserRoomID, error) {
	hashKey := fmt.Sprintf("publish_tid_%s", tid) // 例如 publish_tid_aaa-bbb-ccc

	allValues, err := rds.HGetAll(hashKey).Result()
	if err != nil {
		return nil, err
	}

	if len(allValues) == 0 {
		return make([]*form.UserRoomID, 0), nil
	}

	tasks := make([]*form.UserRoomID, 0, len(allValues))
	for _, v := range allValues {
		var pd []*form.UserRoomID
		if err := json.Unmarshal([]byte(v), &pd); err == nil {
			tasks = append(tasks, pd...)
		}
	}

	if len(tasks) == 0 {
		return make([]*form.UserRoomID, 0), nil
	}

	return tasks, nil
}

func (r *RM) delPublishTid(tid string) ([]*form.UserRoomID, error) {
	hashKey := fmt.Sprintf("publish_tid_%s", tid) // 例如 publish_tid_aaa-bbb-ccc
	var fds = make([]*form.UserRoomID, 0)

	result, err := rds.HDel(hashKey, tid).Result()
	if err != nil {
		return fds, err
	}

	if result == 0 {
		return fds, nil
	}

	return fds, nil

}
