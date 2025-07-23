package qqMapApi

import (
	"encoding/json"
	"fmt"
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
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Address  string   `json:"address"`
	Location Location `json:"location"` // 嵌套了 Location 结构体
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
	Id       string   `json:"id"`
	Img      string   `json:"img"`
	Addr     string   `json:"addr"`
	Title    string   `json:"title"`
	UserId   string   `json:"user_id"`
	Lng      float64  `json:"lng"`
	Lat      float64  `json:"lat"`
	Online   string   `json:"online"`
	Distance string   `json:"distance"`
	Tags     []string `json:"tags"`
}

type TxMapApi struct {
	City    string
	KeyWord string
}

func NewTxMapApi(city string) TxMapApi {
	return TxMapApi{
		City: city,
	}
}

func (t TxMapApi) KeyWordSearch() ([]SaveInRedis, error) {
	var fd APIResponse
	var sd = make([]SaveInRedis, 0, 100)
	offset := 1
	for offset <= 5 {
		url := fmt.Sprintf("https://apis.map.qq.com/ws/place/v1/suggestion?region=%s&keyword=篮球场&key=YSRBZ-GSVY3-3P23L-RNWCE-OQB3V-T6BXG&page_size=20&page_index=%d", t.City, offset)
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
			sdd := SaveInRedis{
				Id:     v.ID,
				UserId: "ogR3E62jXXJMbVcImRqMA1gTSegM",
				Addr:   v.Address,
				Lat:    v.Location.Lat,
				Lng:    v.Location.Lng,
				Title:  v.Title,
				Tags:   []string{v.Title},
				Img:    "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/bk3.svg",
			}
			sd = append(sd, sdd)
		}

		offset++
		time.Sleep(time.Second * time.Duration(4))
	}

	return sd, nil
}
