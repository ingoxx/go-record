package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ingoxx/go-record/http/wx/pkg/config"
	"github.com/ingoxx/go-record/http/wx/pkg/form"
	"io"
	"log"
	"net/http"
	"time"
)

type TxMapApi struct {
	City    string
	KeyWord string
	Url     string
}

func NewTxMapApi(url, city, keyWord string) TxMapApi {
	return TxMapApi{
		City:    city,
		KeyWord: keyWord,
		Url:     url,
	}
}

func (t TxMapApi) KeyWordSearch() ([]*form.SaveInRedis, error) {
	var fd form.TxAPIResponse
	var sd = make([]*form.SaveInRedis, 0, 100)
	var count int
	totalPages := 1
	offset := 1

	for offset <= totalPages {
		//url := fmt.Sprintf("%s/ws/place/v1/suggestion?region=%s&keyword=%s&key=%s&page_size=20&page_index=%d", t.Url, t.City, t.KeyWord, config.QqKey, offset)
		url := fmt.Sprintf("%s/ws/place/v1/search?boundary=region(%s,)&keyword=%s&key=%s&page_size=20&page_index=%d", t.Url, t.City, t.KeyWord, config.QqKey, offset)
		fmt.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			return sd, err
		}

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return sd, err
		}

		if err := json.Unmarshal(b, &fd); err != nil {
			return sd, err
		}

		if count == 0 {
			count = fd.Count
			if count == 0 {
				return sd, errors.New("没有找到场地数据")
			}
			quotient := count / 20
			remainder := count % 20
			if remainder > 0 {
				totalPages = quotient + 1
			} else {
				totalPages = quotient
			}

			fmt.Printf("找到总：%d条数据, %d页\n", count, totalPages)
		}

		for _, v := range fd.Data {
			u4 := uuid.New()
			sdd := &form.SaveInRedis{
				Id:     u4.String(),
				UserId: config.Admin,
				Addr:   v.Address,
				Lat:    v.Location.Lat,
				Lng:    v.Location.Lng,
				Title:  v.Title,
				Tags:   []string{v.Title},
				Img:    "",
				Aid:    v.ID,
			}
			sd = append(sd, sdd)
		}

		offset++
		time.Sleep(500 * time.Millisecond)
	}

	return sd, nil
}

func uniqueByName(people []*form.SaveInRedis) []*form.SaveInRedis {
	seen := make(map[string]bool) // 记录已经出现的Name
	result := make([]*form.SaveInRedis, 0, len(people))

	for _, p := range people {
		if !seen[p.Title] { // 如果没出现过该Name
			seen[p.Title] = true
			result = append(result, p)
		}
	}
	return result
}

func main() {
	search, err := NewTxMapApi("https://apis.map.qq.com", "深圳市", "篮球场").KeyWordSearch()
	if err != nil {
		log.Fatalln(err)
	}

	ns := uniqueByName(search)

	fmt.Println("去重后的总的数据：", len(ns))

	for _, v := range ns {
		fmt.Println("-------------------------")
		fmt.Println(v.Addr)
		fmt.Println(v.Title)
	}

}
