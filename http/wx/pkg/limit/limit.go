package main

import (
	"encoding/json"
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
	offset := 1

	for offset <= 5 {
		url := fmt.Sprintf("%s/ws/place/v1/suggestion?region=%s&keyword=%s&key=%s&page_size=20&page_index=%d", t.Url, t.City, t.KeyWord, config.QqKey, offset)
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
		time.Sleep(time.Second)
	}

	return sd, nil
}

func main() {
	search, err := NewTxMapApi("https://apis.map.qq.com", "深圳市", "攀岩馆").KeyWordSearch()
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range search {
		fmt.Println("-------------------------")
		fmt.Println(v.Addr)
		fmt.Println(v.Title)
	}

}
