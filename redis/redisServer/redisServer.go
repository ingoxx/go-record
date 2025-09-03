package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ingoxx/go-record/http/wx/pkg/config"
	"github.com/ingoxx/go-record/http/wx/pkg/form"
	"log"
	"strings"
	"time"
)

var rdPool = redis.NewClient(
	&redis.Options{
		Addr:         config.RedisAddr,
		DB:           1,
		MinIdleConns: 5,
		Password:     config.RedPwd,
		PoolSize:     5,
		PoolTimeout:  30 * time.Second,
		DialTimeout:  1 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	},
)

// SaveInRedis 写入redis的格式
type SaveInRedis struct {
	JoinUsers              []form.JoinGroupUsers `json:"join_users"` // 某个运动场地，用户点击加入组件的人数
	UserReviews            []*form.MsgBoard      `json:"user_reviews"`
	VenueUpdateUsers       []*form.AddrListForm  `json:"venue_update_users"`
	Tags                   []string              `json:"tags"`
	Images                 []string              `json:"images"`
	Id                     string                `json:"id"`
	Img                    string                `json:"img"`
	Addr                   string                `json:"addr"`
	Title                  string                `json:"title"`
	UserId                 string                `json:"user_id"`
	Online                 string                `json:"online"`
	Distance               string                `json:"distance"`
	Aid                    string                `json:"aid"` // 接口返回的地址唯一id，再次请求接口返回的id是一致的，更新的时候有用
	JoinUserCount          int                   `json:"join_user_count"`
	UserReviewsCount       int                   `json:"user_reviews_count"`
	VenueUpdateUsersCount  int                   `json:"venue_update_users_count"`
	Lng                    float64               `json:"lng"`
	Lat                    float64               `json:"lat"`
	DisVal                 float64               `json:"dis_val"`
	IsShow                 bool                  `json:"is_show"`
	IsShowUserReviews      bool                  `json:"is_show_user_reviews"`
	IsShowJoinUsers        bool                  `json:"is_show_join_users"`
	IsShowVenueUpdateUsers bool                  `json:"is_show_venue_update_users"`
}

type SportInfo struct {
	Name    string `json:"name"`
	Key     string `json:"key"`
	Img     string `json:"img"`
	Icon    string `json:"icon"`
	Checked bool   `json:"checked"`
}

type WxOpenidList struct {
	Openid   string `json:"openid"`
	Img      string `json:"img"`
	NickName string `json:"nick_name"`
	Time     string `json:"time"`
}

type AddrListForm struct {
	UserImg    string  `json:"user_img"` // 用户头像
	Content    string  `json:"content"`  // 更新内容,目前只能统一更新图片,这里都写: 更新了场地图片
	NickName   string  `json:"nick_name"`
	Tags       string  `json:"tags"  validate:"required"`
	Id         string  `json:"id" validate:"required"`
	Addr       string  `json:"addr" validate:"required"`
	UserId     string  `json:"user_id" validate:"required"`      // 添加场地的用户id
	City       string  `json:"city"  validate:"required"`        // 前端传入的是中文
	CityPy     string  `json:"city_py"`                          // 前端传入的中文转成拼音
	SportKey   string  `json:"sport_key" validate:"required"`    // 运动分类，篮球：shenzhenshi_bks,足球：shenzhenshi_fbs...
	UpdateType string  `json:"update_type"  validate:"required"` // 更新类型：1.用户添加的新场地，2.用户更新了场地
	Aid        string  `json:"aid"`                              // api返回的场地的唯一id，就是再次请求返回的id都是一样的
	Img        string  `json:"img"`                              // 场地图片
	Time       string  `json:"time"`                             // 更新时间
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	IsRecord   bool    `json:"is_record"` // true：已记录（审核通过），false：未记录（还未审核通过）
	IsShow     bool    `json:"is_show"`   // 审核列表中的数据，true：隐藏，false：不隐藏
}

func main() {
	getAllData()
}

func updateSquare() {
	var data []*SaveInRedis
	var sportData []SportInfo
	var sportMap = make(map[string]string)
	sports := `[
		{"name": "篮球场", "key": "bks", "checked": false, "icon": "🏀", "img": "https://ai.anythingai.online/static/profile3/main-bk.jpg"},
		{"name": "游泳馆", "key": "sws", "checked": false, "icon": "🏊", "img": "https://ai.anythingai.online/static/profile3/swim.png"},
		{"name": "羽毛球馆", "key": "bms", "checked": false, "icon": "🏸", "img": "https://ai.anythingai.online/static/profile3/badminton.png"},
		{"name": "足球场", "key": "fbs", "checked": false, "icon": "⚽", "img": "https://ai.anythingai.online/static/profile3/football.png"},
		{"name": "网球场", "key": "tns", "checked": false, "icon": "🎾", "img": "https://ai.anythingai.online/static/profile3/tennis.png"},
		{"name": "高尔夫球场", "key": "gos", "checked": false, "icon": "🏌️", "img": "https://ai.anythingai.online/static/profile3/golf.png"},
		{"name": "滑雪场", "key": "hxc", "checked": false, "icon": "⛷️", "img": "https://ai.anythingai.online/static/profile3/ski.png"},
		{"name": "瑜伽馆", "key": "yjg", "checked": false, "icon": "🧘", "img": "https://ai.anythingai.online/static/profile3/yoga.png"},
		{"name": "跆拳道馆", "key": "tqd", "checked": false, "icon": "🥋", "img": "https://ai.anythingai.online/static/profile3/taekwondo.png"},
		{"name": "健身房", "key": "gym", "checked": false, "icon": "🏋️‍♂️", "img": "https://ai.anythingai.online/static/profile3/gym.png"}
	]`
	if err := json.Unmarshal([]byte(sports), &sportData); err != nil {
		log.Fatalln(err)
	}
	for _, v := range sportData {
		sportMap[v.Key] = v.Img
	}

	keys := []string{
		//"group_id_zhuhaishi_bms",
		//"group_id_zhuhaishi_tqd",
		//"group_id_zhuhaishi_gos",
		//"group_id_zhuhaishi_yjg",
		"group_id_zhuhaishi_bks",
		//"group_id_zhuhaishi_fbs",
		//"group_id_zhuhaishi_tns",
		//"group_id_zhuhaishi_sws",
		//"group_id_zhuhaishi_gym",
	}

	for _, k := range keys {
		sport := strings.Split(k, "_")

		result, err := rdPool.Get(k).Result()
		if err != nil {
			log.Fatalln(k, err)
		}

		if err := json.Unmarshal([]byte(result), &data); err != nil {
			log.Fatalln(err)
		}

		for _, v := range data {
			kd, ok := sportMap[sport[3]]
			fmt.Println()
			if ok {
				v.Img = kd
			}
		}

		b, err := json.Marshal(&data)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = rdPool.Set(k, b, 0).Result()
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("%s update ok\n", k)
		//return
	}
}

func getAllData() {
	k := []string{
		//"group_id_guangzhoushi_fbs",
		//"group_id_zhangshashi_bms",
		//"group_id_anyangshi_bks",
		//"group_id_shenzhenshi_bms",
		//"group_id_guangzhoushi_sws",
		//"group_id_shanghaishi_bks",
		//"group_id_shenzhenshi_gos",
		//"group_id_zhongqingshi_bks",
		//"group_id_zhangchunshi_bks",
		//"group_id_shenzhenshi_yjg",
		//"group_id_daqingshi_bks",
		//"group_id_guangzhoushi_bks",
		//"group_id_huizhoushi_bms",
		//"group_id_shenzhenshi_gym",
		//"group_id_shenzhenshi_fbs",
		//"group_id_jinhuashi_bks",
		//"group_id_chuzhoushi_bks",
		//"group_id_shenzhenshi_sws",
		//"group_id_taizhoushi_bks",
		//"group_id_zhongqingshi_fbs",
		"group_id_shenzhenshi_bks",
		//"group_id_heyuanshi_bks",
		//"group_id_changzhoushi_bks",
		//"group_id_shanghaishi_bms",
		//"group_id_chengdoushi_bks",
		//"group_id_guangzhoushi_yjg",
		//"group_id_shenzhenshi_tns",
		//"group_id_huizhoushi_bks",
		//"group_id_zhongqingshi_bms",
	}

	for _, vk := range k {
		var data []*SaveInRedis
		result, err := rdPool.Get(vk).Result()
		if err != nil {
			log.Fatalln(k, err)
		}

		if err := json.Unmarshal([]byte(result), &data); err != nil {
			log.Fatalln(err)
		}

		total := len(data)
		var countImg int

		for _, vd := range data {
			if vd.Img != "" {
				countImg++
			}
			fmt.Println("--------------------------------")
			fmt.Println(vd.Title)

		}

		fmt.Printf("总：%d, 没有图片的：%d\n", total, total-countImg)
		//b, err := json.Marshal(&data)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//
		//if _, err := rdPool.Set(vk, b, 0).Result(); err != nil {
		//	log.Fatalln(err)
		//}

		//log.Printf("%s , ok\n", vk)
	}

}

func getOpenIdList() {
	var data []*WxOpenidList
	k := "group_id_wx_open_id_list"
	result, err := rdPool.Get(k).Result()
	if err != nil {
		log.Fatalln(k, err)
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		log.Fatalln(err)
	}

	for _, v := range data {
		v.Img = strings.ReplaceAll(v.Img, "anyhingai", "anythingai")
	}

	b, err := json.Marshal(&data)
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := rdPool.Set(k, b, 0).Result(); err != nil {
		log.Fatalln(err)
	}

}

func getCheckList() {
	var data []*AddrListForm
	k := "group_id_addr_check_list"
	result, err := rdPool.Get(k).Result()
	if err != nil {
		log.Fatalln(k, err)
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		log.Fatalln(err)
	}

	for _, v := range data {
		//user := getUserInfo(v.UserId)
		//if user.Img != "" {
		//	v.UserImg = user.Img
		//}
		//if user.NickName != "" {
		//	v.NickName = user.NickName
		//}
		//if v.UpdateType == "1" {
		//	v.Content = "用户添加了新的场地"
		//}
		//if v.UpdateType == "2" {
		//	v.Content = "用户更新了场地图片"
		//}
		if v.Time == "" {
			v.Time = time.Now().Format("2006-01-02 15:04:05")
		}

		fmt.Println(v.Tags, v.Time)
	}

	b, err := json.Marshal(&data)
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := rdPool.Set(k, b, 0).Result(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("save ok")
}

func getUserInfo(uid string) *WxOpenidList {
	var data []*WxOpenidList
	k := "group_id_wx_open_id_list"
	result, err := rdPool.Get(k).Result()
	if err != nil {
		log.Fatalln(k, err)
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		log.Fatalln(err)
	}

	for _, v := range data {
		if v.Openid == uid {
			fmt.Println(v.NickName, v.Openid)
			return v
		}
	}

	return new(WxOpenidList)
}
