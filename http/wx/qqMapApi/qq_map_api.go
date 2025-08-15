package qqMapApi

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/ingoxx/go-record/http/wx/form"
	"io"
	"net/http"
	"time"
)

// Location 对应 JSON 中的 "location" 对象
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// DataItem 对应 JSON 中 "data" 数组的每个元素
type DataItem struct {
	Location Location `json:"location"` // 嵌套了 Location 结构体
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Address  string   `json:"address"`
	Province string   `json:"province"`
	City     string   `json:"city"`
	District string   `json:"district"`
	// 其他字段如果不需要可以不定义
}

// APIResponse 对应整个 JSON 响应
type APIResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []DataItem `json:"data"` // 这是一个 DataItem 的切片(数组)
	Count   int        `json:"count"`
}

// SaveInRedis 写入redis的格式
type SaveInRedis struct {
	JoinUsers        []form.JoinGroupUsers `json:"join_users"` // 某个运动场地，用户点击加入组件的人数
	UserReviews      []form.MsgBoard       `json:"user_reviews"`
	Tags             []string              `json:"tags"`
	Id               string                `json:"id"`
	Img              string                `json:"img"`
	Addr             string                `json:"addr"`
	Title            string                `json:"title"`
	UserId           string                `json:"user_id"`
	Online           string                `json:"online"`
	Distance         string                `json:"distance"`
	JoinUserCount    int                   `json:"join_user_count"`
	UserReviewsCount int                   `json:"user_reviews_count"`
	Lng              float64               `json:"lng"`
	Lat              float64               `json:"lat"`
}

type TxMapApi struct {
	City    string
	KeyWord string
}

func NewTxMapApi(city, keyWord string) TxMapApi {
	return TxMapApi{
		City:    city,
		KeyWord: keyWord,
	}
}

func (t TxMapApi) KeyWordSearch() ([]SaveInRedis, error) {
	var fd APIResponse
	var sd = make([]SaveInRedis, 0, 100)
	offset := 1

	for offset <= 5 {
		url := fmt.Sprintf("https://apis.map.qq.com/ws/place/v1/suggestion?region=%s&keyword=%s&key=YSRBZ-GSVY3-3P23L-RNWCE-OQB3V-T6BXG&page_size=20&page_index=%d", t.City, t.KeyWord, offset)
		resp, err := http.Get(url)
		if err != nil {
			return sd, err
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return sd, err
		}

		if err := json.Unmarshal(b, &fd); err != nil {
			return sd, err
		}

		for _, v := range fd.Data {
			u4 := uuid.New()
			sdd := SaveInRedis{
				Id:     u4.String(),
				UserId: "ogR3E62jXXJMbVcImRqMA1gTSegM",
				Addr:   v.Address,
				Lat:    v.Location.Lat,
				Lng:    v.Location.Lng,
				Title:  v.Title,
				Tags:   []string{v.Title},
				Img:    "",
			}
			sd = append(sd, sdd)
		}

		offset++
		time.Sleep(time.Second)
	}

	return sd, nil
}
