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
	JoinUsers        []form.JoinGroupUsers `json:"join_users"` // 某个运动场地，用户点击加入组件的人数
	UserReviews      []form.MsgBoard       `json:"user_reviews"`
	Tags             []string              `json:"tags"`
	Id               string                `json:"id"`
	Img              string                `json:"img"`
	Images           []string              `json:"images"`
	Addr             string                `json:"addr"`
	Title            string                `json:"title"`
	UserId           string                `json:"user_id"`
	Online           string                `json:"online"`
	Distance         string                `json:"distance"`
	Aid              string                `json:"aid"` // 接口返回的地址唯一id，再次请求接口返回的id是一致的，更新的时候有用
	JoinUserCount    int                   `json:"join_user_count"`
	UserReviewsCount int                   `json:"user_reviews_count"`
	Lng              float64               `json:"lng"`
	Lat              float64               `json:"lat"`
	IsShow           bool                  `json:"is_show"`
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
		//"group_id_hangzhoushi_bks",
		//"group_id_heyuanshi_fbs",
		//"group_id_dongguanshi_sws",
		//"group_id_heyuanshi_tns",
		//"group_id_shenzhenshi_bms",
		//"group_id_shanghaishi_bks",
		//"group_id_shenzhenshi_tqd",
		//"group_id_shenzhenshi_gos",
		//"group_id_hangzhoushi_bms",
		//"group_id_heyuanshi_sws",
		//"group_id_shenzhenshi_yjg",
		"group_id_shenzhenshi_bks",
		//"group_id_heyuanshi_bks",
		//"group_id_heyuanshi_bms",
		//"group_id_guangzhoushi_bks",
		//"group_id_shenzhenshi_fbs",
		//"group_id_huizhoushi_tns",
		//"group_id_jinanshi_tns",
		//"group_id_shenzhenshi_tns",
		//"group_id_chuzhoushi_bks",
		//"group_id_zhuhaishi_bks",
		//"group_id_huizhoushi_bks",
		//"group_id_heyuanshi_gym",
		//"group_id_shenzhenshi_sws",
		//"group_id_heyuanshi_yjg",
		//"group_id_suzhoushi_bks",
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

		fmt.Println("count >>> ", len(data))

		for _, vd := range data {
			fmt.Println("------------------------------------------------------")
			fmt.Println(vd.Title, vd.Addr, vd.Images, vd.Lat, vd.Lng)
		}

		//b, err := json.Marshal(&data)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//
		//if _, err := rdPool.Set(vk, b, 0).Result(); err != nil {
		//	log.Fatalln(err)
		//}
		//
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
		v.Time = time.Now().Format("2006-01-02 15:04:05")
	}

	b, err := json.Marshal(&data)
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range data {
		fmt.Println(v.Openid, v.NickName, v.Img, v.Time)
	}

	if _, err := rdPool.Set(k, b, 0).Result(); err != nil {
		log.Fatalln(err)
	}

}
