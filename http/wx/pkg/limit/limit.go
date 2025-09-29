package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ingoxx/go-record/http/wx/pkg/config"
	"github.com/ingoxx/go-record/http/wx/pkg/form"
	"io"
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
		url := fmt.Sprintf("%s/ws/place/v1/suggestion?region=%s&keyword=%s&key=%s&page_size=20&page_index=%d", t.Url, t.City, t.KeyWord, config.QqKey, offset)
		//url := fmt.Sprintf("%s/ws/place/v1/search?boundary=region(%s,)&keyword=%s&key=%s&page_size=20&page_index=%d", t.Url, t.City, t.KeyWord, config.QqKey, offset)
		resp, err := http.Get(url)
		if err != nil {
			return sd, err
		}

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return sd, err
		}

		fmt.Println(string(b))

		if err := json.Unmarshal(b, &fd); err != nil {
			fmt.Println(err)
			return sd, err
		}

		fmt.Println(fd)
		return sd, err
		if count == 0 {
			count = fd.Count
			if count == 0 {
				return sd, errors.New(fd.Message)
			}
			quotient := count / 20
			remainder := count % 20
			if remainder > 0 {
				totalPages = quotient + 1
			} else {
				totalPages = quotient
			}
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
		return sd, nil
	}

	return sd, nil
}

func main() {
	NewTxMapApi("https://restapi.amap.com", "深圳市", "公共厕所").KeyWordSearch()
}
