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
	City     string `json:"city"`
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
	getOpenIdList()
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
		var data []*form.SaveInRedis
		result, err := rdPool.Get(vk).Result()
		if err != nil {
			log.Fatalln(k, err)
		}

		if err := json.Unmarshal([]byte(result), &data); err != nil {
			log.Fatalln(err)
		}

		for _, vd := range data {
			if vd.Id == "c78a5b82-d375-4b0d-998d-4dd8ea8d93ef" {
				vd.Img = ""
			}
			fmt.Println("--------------------------------")
			fmt.Println(vd.Title, vd.Id)
			fmt.Println(vd.Img)
			fmt.Println(vd.Id)

		}

		b, err := json.Marshal(&data)
		if err != nil {
			log.Fatalln(err)
		}

		if _, err := rdPool.Set(vk, b, 0).Result(); err != nil {
			log.Fatalln(err)
		}

		log.Printf("%s , update ok\n", vk)
	}

}

func getOpenIdList() {
	var data []*WxOpenidList
	k := "group_id_wx_open_id_list"
	res, err := rdPool.Get(k).Result()
	if err != nil {
		log.Fatalln(k, err)
	}

	//d2 := `[{"openid":"ogR3E62jXXJMbVcImRqMA1gTSegM","img":"https://ai.anythingai.online/static/profile3/1703.png","nick_name":"xe7xafxaexe7xadx90xe7x9ax84xe5xaex88xe9x97xa8xe5x91x98","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6zO5aFz3RF-CzI14q7U_sHI","img":"https://ai.anythingai.online/static/profile3/2696.png","nick_name":"xe4xbdxa0xe5xa4xa7xe7x88xb7","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E66lDUu0_IKhGyiIn23QpKYE","img":"https://ai.anythingai.online/static/profile3/2047.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6yxl6TFwiQi2UERqANPgleA","img":"https://ai.anythingai.online/static/profile3/2442.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E61KswqU7BmBE3oGcWsybgX0","img":"https://ai.anythingai.online/static/profile3/2048.png","nick_name":"xe4xb8x8axe7xafxaexe6x92x9exe5xa2x99xe7x8ex8b","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6zdA5DjzO5lcht1c3vK83s8","img":"https://ai.anythingai.online/static/profile3/2006.png","nick_name":"xe4xb8x8axe7xafxaexe8xa6x81xe6x89xb6xe6xa2xaf","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E62bbHDCAIPyQl_ZiqYPb9n8","img":"https://ai.anythingai.online/static/profile3/1411.png","nick_name":"xe8x8fx9cxe8xa1xa3xe5x87x9b","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E68WMdzd-iQIp7iwfFWuT4r0","img":"https://ai.anythingai.online/static/profile3/2051.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6-d12WJ4cJfZjON3ADKDcuA","img":"https://ai.anythingai.online/static/profile3/2586.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6zmsARCJ0suIGayyJYgm9Vk","img":"https://ai.anythingai.online/static/profile3/1979.png","nick_name":"xe5x8fx91xe8xb4xa2","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E69bMTdhdZ3TL1qmHlJuqf0U","img":"https://ai.anythingai.online/static/profile3/2168.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6yEY9MRMzAQv4fEK1ifN1sg","img":"https://ai.anythingai.online/static/profile3/2629.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E680VlNOWKttV55A1-SkCsGU","img":"https://ai.anythingai.online/static/profile3/1060.png","nick_name":"-","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E65ij4k9wVn6dzhTACA385js","img":"https://ai.anythingai.online/static/profile3/2153.png","nick_name":"xe7x90x83xe5x9cxbaxe6x8cx87xe5xaex9axe8x83x8cxe9x94x85xe4xbexa0","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E69jUMjr36H3SwDNUjvbbBU8","img":"https://ai.anythingai.online/static/profile3/1665.png","nick_name":"xe8xbfx90xe7x90x83xe5x88xb0xe7x95x8cxe5xa4x96","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E692U7nglaC15oTjIhxI2-gY","img":"https://ai.anythingai.online/static/profile3/2883.png","nick_name":"xe7x90x83xe5x9cxbaxe5x81x87xe5x8axa8xe4xbdx9cxe5xa4xa7xe5xb8x88","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6-Ilm2FKGn1X1P_mnCprXXg","img":"https://ai.anythingai.online/static/profile3/2067.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6502vro1C14oZ01xrvZ8wWo","img":"https://ai.anythingai.online/static/profile3/2294.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E67K4KJfOocda2BwpmExn81I","img":"https://ai.anythingai.online/static/profile3/1672.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E65WpF41OrghI7nTrh7qMQpU","img":"https://ai.anythingai.online/static/profile3/2733.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6w0IjFFcnAKAY6J8CwYU_6k","img":"https://ai.anythingai.online/static/profile3/2489.png","nick_name":"xe7x9bx96xe5xb8xbdxe6x89x93xe8x87xaaxe5xb7xb1xe8x84xb8","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6wFWrvZ_CPwqr5ojm4xeYtg","img":"https://ai.anythingai.online/static/profile3/1277.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E65nkP1qaYADD4uX8dSrrb6M","img":"https://ai.anythingai.online/static/profile3/2688.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E641QbzokQ8i29YixN5JGk_8","img":"https://ai.anythingai.online/static/profile3/1373.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E65YgxSz42t_kdWU0tPjWG9Y","img":"https://ai.anythingai.online/static/profile3/2639.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6-_gIzIoF1qK6YkUYQW7aME","img":"https://ai.anythingai.online/static/profile3/1577.png","nick_name":"xe7x90x83xe8xa1xa3xe6xb0xb8xe8xbfx9cxe6x98xafxe5xb9xb2xe5x87x80xe7x9ax84","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E600qHlIGM_pFkzT18wgK-1I","img":"https://ai.anythingai.online/static/profile3/1507.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E662DDO1LSR0Jj5J_yoW9yY0","img":"https://ai.anythingai.online/static/profile3/2653.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6xMOsaT1UICJXmD2rcE5Rq8","img":"https://ai.anythingai.online/static/profile3/1558.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6wjN0TlTUjoDXZ14YcQ5JT4","img":"https://ai.anythingai.online/static/profile3/2512.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6xHtJYGAHGlW9BFRtdAkNBA","img":"https://ai.anythingai.online/static/profile3/2573.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E69naLQPBEI29uYxIvWhmKXk","img":"https://ai.anythingai.online/static/profile3/2788.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E66nhJe_CisLi7IjtToGxZSw","img":"https://ai.anythingai.online/static/profile3/1816.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E65kVt75GRHrQr7sVha4-qxw","img":"https://ai.anythingai.online/static/profile3/1294.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E68wL8TZFmm2WgB9P-gfKl0E","img":"https://ai.anythingai.online/static/profile3/1833.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E68Xjh-coBg4d62xEEyXfACg","img":"https://ai.anythingai.online/static/profile3/2205.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E61Ja7cVVQbE8dLulajUv-90","img":"https://ai.anythingai.online/static/profile3/1071.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E62aCwdwhAnW45HL7wRnmgSY","img":"https://ai.anythingai.online/static/profile3/2138.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E69xzm3kqEbDpEZWETgee5wQ","img":"https://ai.anythingai.online/static/profile3/2511.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E65NO4z-3qJn6DhqYFMmAnow","img":"https://ai.anythingai.online/static/profile3/1713.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E66V6uUQ1sY0geCRwOS4bPIQ","img":"https://ai.anythingai.online/static/profile3/2348.png","nick_name":"xe7xafxaexe7x90x83xe5x9cxbaxe6x91x84xe5xbdxb1xe5xb8x88","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6080Fa9T9H4MZegGBgw63dU","img":"https://ai.anythingai.online/static/profile3/2409.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E632-ACiWwTiu4gV0JhchFQY","img":"https://ai.anythingai.online/static/profile3/1771.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6yiMBBJNnlHc7MhHwnMnAzI","img":"https://ai.anythingai.online/static/profile3/2828.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6-OBu47vwoDu-yeKnW9Pyxc","img":"https://ai.anythingai.online/static/profile3/2671.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6zQb6fSR6EZA2VN87gHPhA4","img":"https://ai.anythingai.online/static/profile3/1033.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E60kuhCuiwVKN4Gz8Fe9fjqo","img":"https://ai.anythingai.online/static/profile3/2126.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E66JTrncbHixz1zUmk9aewl0","img":"https://ai.anythingai.online/static/profile3/2366.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E68hAQekV3WyYGhRI1pKK3dY","img":"https://ai.anythingai.online/static/profile3/2065.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E60eaJwporZ-55sUPke2Pp40","img":"https://ai.anythingai.online/static/profile3/1119.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E69UBYKxS-DSPOBde-laH9_4","img":"https://ai.anythingai.online/static/profile3/1767.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E67t9GWmUEfYCFdX-HqfF4mE","img":"https://ai.anythingai.online/static/profile3/2497.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E60K2A96SR4gLFQ6ijN2tnGA","img":"https://ai.anythingai.online/static/profile3/1927.png","nick_name":"xe6x89xa3xe7xafxaexe9x9dxa0xe6x84x8fxe5xbfxb5","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6w6N5TVDdE6N9oRnE0e7Ros","img":"https://ai.anythingai.online/static/profile3/1417.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6xykWO6Y7P63qp6Ps8HOpeI","img":"https://ai.anythingai.online/static/profile3/2774.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6-W56jMCYNvaqeKU7ZWgG88","img":"https://ai.anythingai.online/static/profile3/1084.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E60gSsHQg8nMec_868i1iK0w","img":"https://ai.anythingai.online/static/profile3/1258.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E637PWZ45a9W7diR00whzgMU","img":"https://ai.anythingai.online/static/profile3/2382.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E623CcBSpIp6kTRKtIvMkg44","img":"https://ai.anythingai.online/static/profile3/2192.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E65ePjUBZkUR9jJJC00g4liM","img":"https://ai.anythingai.online/static/profile3/1154.png","nick_name":"xe7xafxaexe7xadx90xe5xaex88xe6x8axa4xe7xa5x9e","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6xO5_4MLxyLQX-3cEW98Ddw","img":"https://ai.anythingai.online/static/profile3/1691.png","nick_name":"xe6x8axa2xe6x96xadxe5xa4xb1xe8xb4xa5xe5xa4xa7xe7x8ex8b","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6-9-8peACOPfIt3Xk__7qjs","img":"https://ai.anythingai.online/static/profile3/1559.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E60JI1kOI6C-v-BzQtAHZHy4","img":"https://ai.anythingai.online/static/profile3/2844.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6x3aB8c8JiFYfSCsDjqLMAE","img":"https://ai.anythingai.online/static/profile3/2725.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E61KKn81WSl0ejL0xJO4t3Uc","img":"https://ai.anythingai.online/static/profile3/2352.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E64A6bCOX1hVQ3lXXxTgDFcM","img":"https://ai.anythingai.online/static/profile3/1186.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E68NK4krk1Pll-I5XqdUMmto","img":"https://ai.anythingai.online/static/profile3/1972.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E65jOXOuJ8zDu8QG_y3vz6yY","img":"https://ai.anythingai.online/static/profile3/1045.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6yqBdwS_bOakbMBB8aGDj9s","img":"https://ai.anythingai.online/static/profile3/2478.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6xFdRKL5_3jvKzVo6_SgnsU","img":"https://ai.anythingai.online/static/profile3/2806.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E68jIbeM_GJC376IaTmxOriw","img":"https://ai.anythingai.online/static/profile3/1043.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6x0wLEXztkRMtz3W6TbZ6aY","img":"https://ai.anythingai.online/static/profile3/1063.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6zynCxmC_pX5fW3TaRwhCPY","img":"https://ai.anythingai.online/static/profile3/1669.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6yGxK-b8GAleAnf690G5dTU","img":"https://ai.anythingai.online/static/profile3/2891.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6xoZ1FwhjRDPBP57VGMvy2U","img":"https://ai.anythingai.online/static/profile3/1025.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E69QOIc-TwexKl53UCoNljDU","img":"https://ai.anythingai.online/static/profile3/2472.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6zOHioKHAPAItC-hsAhKCTU","img":"https://ai.anythingai.online/static/profile3/1846.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6ySvjbRhKaX_SrM0PP_w5Kg","img":"https://ai.anythingai.online/static/profile3/2432.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E64U7TgO6WQJpH17Q2ISS4Fc","img":"https://ai.anythingai.online/static/profile3/2225.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E65Ku9QXVb_rl_j3bR0lFKe4","img":"https://ai.anythingai.online/static/profile3/1996.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6-0_RJNPhdaicf6g72h8WPo","img":"https://ai.anythingai.online/static/profile3/1133.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E67lWuXyMeOgcDqoS3q00np8","img":"https://ai.anythingai.online/static/profile3/2851.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E64LcbChwBvv2YJRCDJwLWKY","img":"https://ai.anythingai.online/static/profile3/2852.png","nick_name":"","time":"2025-08-25 14:50:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E61cimsuHXdyIducYRPIIcDg","img":"https://ai.anythingai.online/static/profile3/1748.png","nick_name":"","time":"2025-08-25 15:08:57","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E63SauQTInZzVuEtcrXDDscA","img":"https://ai.anythingai.online/static/profile3/1225.png","nick_name":"xe7x90x83xe5x9cxbaxe8x83x8cxe6x99xafxe6x9dxbf","time":"2025-08-25 15:28:36","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E62VLEzxh67KKJjW53Giv4Fg","img":"https://ai.anythingai.online/static/profile3/2791.png","nick_name":"","time":"2025-08-25 21:16:35","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6x8uyKtJ11qCIZLFzWGnVl0","img":"https://ai.anythingai.online/static/profile3/1970.png","nick_name":"","time":"2025-08-26 16:32:20","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E69hEhHUOb77DrNSR3gHgzx0","img":"https://ai.anythingai.online/static/profile3/2259.png","nick_name":"","time":"2025-08-26 22:18:18","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6wuhhwWZZseb5z-KPjzwAFc","img":"https://ai.anythingai.online/static/profile3/2235.png","nick_name":"","time":"2025-08-27 18:13:41","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6_msl_12_aiPMSVCdjUmd7Q","img":"https://ai.anythingai.online/static/profile3/1638.png","nick_name":"","time":"2025-08-27 22:42:54","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6zoHOuKfO8mXpc7I2C5QE-k","img":"https://ai.anythingai.online/static/profile3/2736.png","nick_name":"","time":"2025-08-28 18:21:00","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E61LPggRK6GgaTua-mJkGOOw","img":"https://ai.anythingai.online/static/profile3/1917.png","nick_name":"","time":"2025-08-28 18:33:08","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6xuwqQ9KO9EtbgxxQ93yAAk","img":"https://ai.anythingai.online/static/profile3/2033.png","nick_name":"","time":"2025-08-28 18:54:55","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E69flBRo7lh4bhzmkLzILV7c","img":"https://ai.anythingai.online/static/profile3/1171.png","nick_name":"","time":"2025-08-28 19:30:51","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E630omuwSatFvsz6L_7CHPpM","img":"https://ai.anythingai.online/static/profile3/2574.png","nick_name":"","time":"2025-08-29 07:13:29","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E65McqcC2EL2Zj-DbhKasd7c","img":"https://ai.anythingai.online/static/profile3/1636.png","nick_name":"","time":"2025-08-29 11:52:52","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6xSEfjerLLA_x2MQxAvoUZc","img":"https://ai.anythingai.online/static/profile3/1520.png","nick_name":"xe4xb8x80xe7xa7x92xe9x92x9fxe6x8ex89xe7x90x83xe4xbexa0","time":"2025-08-29 20:20:26","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E64wtf91iN6hPch9yf7Mptmw","img":"https://ai.anythingai.online/static/profile3/2688.png","nick_name":"xe4xb8x80xe7xa7x92xe9x92x9fxe6x8ex89xe7x90x83xe4xbexa0","time":"2025-08-31 22:24:32","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6zbD-dZdu5aapT6nRgkYfGY","img":"https://ai.anythingai.online/static/profile3/2051.png","nick_name":"xe7x90x83xe5x9cxbaxe4xbcxa0xe7x90x83xe9xbbx91xe6xb4x9e","time":"2025-09-01 13:12:32","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"","img":"https://ai.anythingai.online/static/profile3/2098.png","nick_name":"xe7xafxaexe6x9dxbfxe6xbcx8fxe9xa3x8exe7x8ex8b","time":"2025-09-01 18:46:26","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E61kUsVlIfd1wiTlnan_o64M","img":"https://ai.anythingai.online/static/profile3/1565.png","nick_name":"xe8xa1x97xe7x90x83xe6x99x83xe5x80x92xe8x87xaaxe5xb7xb1","time":"2025-09-01 19:15:56","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6zwGPqL-GpuYy6alLuuYp2o","img":"https://ai.anythingai.online/static/profile3/1408.png","nick_name":"xe7xbdx9axe7x90x83xe4xb8x8dxe8xbfx9bxe7xa0x94xe7xa9xb6xe5x91x98","time":"2025-09-01 20:07:31","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6wsod-C27Sfrd8dxfy27MW0","img":"https://ai.anythingai.online/static/profile3/2524.png","nick_name":"xe6x8axa2xe6x96xadxe5xa4xb1xe8xb4xa5xe5xa4xa7xe7x8ex8b","time":"2025-09-01 22:11:14","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6xexaS6A5yRJ0y-1hxczyJM","img":"https://ai.anythingai.online/static/profile3/1092.png","nick_name":"xe7x90x83xe5x9cxbaxe5x88x92xe6xb0xb4xe5xa4xa7xe5xb8x88","time":"2025-09-01 22:52:47","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6wIQpoSAR_T9vWZLvE6fddk","img":"https://ai.anythingai.online/static/profile3/2117.png","nick_name":"xe7xafxaexe7xadx90xe9x93x81xe6x89x93xe7x9ax84xe5x85x84xe5xbcx9f","time":"2025-09-02 14:17:43","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E65GWsa4et-b19UOz-HFMLJ4","img":"https://ai.anythingai.online/static/profile3/1166.png","nick_name":"xe5xa4xb1xe8xafxafxe5x88xb6xe9x80xa0xe6x9cxba","time":"2025-09-02 16:55:38","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6xJqe_CLpPmKj5ZDxFcrYFY","img":"https://ai.anythingai.online/static/profile3/1742.png","nick_name":"xe4xb8x89xe6xadxa5xe8xb5xb0xe6x88x90xe4xbax94xe6xadxa5","time":"2025-09-02 17:36:09","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6_YEve4G7nBta-2SLjo7dsQ","img":"https://ai.anythingai.online/static/profile3/1067.png","nick_name":"xe8xb7x91xe4xbdx8dxe6xb0xb8xe8xbfx9cxe8xb7x91xe9x94x99","time":"2025-09-02 17:37:43","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6-as3V1ajO4Se9ghTkwOtwQ","img":"https://ai.anythingai.online/static/profile3/1962.png","nick_name":"xe8xb7x91xe4xbdx8dxe6xb0xb8xe8xbfx9cxe8xb7x91xe9x94x99","time":"2025-09-03 19:25:11","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6zODXnproPABUyb2eELMZ9k","img":"https://ai.anythingai.online/static/profile3/2418.png","nick_name":"xe4xbcxa0xe7x90x83xe5xa4xb1xe8xafxafxe8x89xbaxe6x9cxafxe5xaexb6","time":"2025-09-03 21:31:01","city":"xe6xb7xb1xe5x9cxb3"},{"openid":"ogR3E6_63yrxycSAM5ow6lHTye7c","img":"https://ai.anythingai.online/static/profile3/1154.png","nick_name":"xe6x8axa2xe6x96xadxe5xa4xb1xe8xb4xa5xe5xa4xa7xe7x8ex8b","time":"2025-09-03 22:33:19","city":"xe6xb7xb1xe5x9cxb3"}]`

	if err := json.Unmarshal([]byte(res), &data); err != nil {
		log.Fatalln(err)
	}

	for _, v := range data {

		//v.City = "深圳"
		//
		//v.NickName = eva.NewSportType("bks").RandomNickname()
		fmt.Println("--------------------------")
		fmt.Println(v.Openid)
		fmt.Println(v.City)

	}

	//b, err := json.Marshal(&data)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//if _, err := rdPool.Set(k, b, 0).Result(); err != nil {
	//	log.Fatalln(err)
	//}
	//
	//fmt.Println("ok")
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

		fmt.Println(v.Tags, v.Time, v.City)
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
