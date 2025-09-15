package bdMapApi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ingoxx/go-record/http/wx/pkg/config"
	"github.com/ingoxx/go-record/http/wx/pkg/form"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type BdMapApi struct {
	City    string
	KeyWord string
	Url     string
}

func NewBdMapApi(url, city, keyWord string) BdMapApi {
	return BdMapApi{
		City:    city,
		KeyWord: keyWord,
		Url:     url,
	}
}

func (t BdMapApi) KeyWordSearch() ([]*form.SaveInRedis, error) {
	var fd form.BdAPIResponse
	var sd = make([]*form.SaveInRedis, 0, 100)
	uri := "/api_place_pro/v1/region?region_limit=true"
	var offset int

	for offset <= 6 {
		fullURL := t.Url + uri + "?" + t.buildData(offset)
		response, err := http.Get(fullURL)
		if err != nil {
			return sd, err
		}

		if response.StatusCode != http.StatusOK {
			return sd, errors.New("请求失败")
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return sd, err
		}

		if err := json.Unmarshal(body, &fd); err != nil {
			return sd, err
		}

		for _, v := range fd.Data {
			u4 := uuid.New()
			addr := fmt.Sprintf("%s%s%s%s", v.City, v.Area, v.Town, v.Name)
			sdd := &form.SaveInRedis{
				Id:     u4.String(),
				UserId: config.Admin,
				Addr:   addr,
				Lat:    v.Location.Lat,
				Lng:    v.Location.Lng,
				Title:  v.Name,
				Tags:   []string{v.Name},
				Img:    "",
				Aid:    v.ID,
				Images: []string{},
			}
			sd = append(sd, sdd)
		}

		offset += 1
	}

	return sd, nil
}

func (t BdMapApi) buildData(offset int) string {
	ost := strconv.Itoa(offset)
	params := url.Values{}
	params.Set("page_num", ost)
	params.Set("page_size", "20")
	params.Set("region_limit", "true")
	params.Set("scope", "2")
	params.Set("output", "json")
	params.Set("extensions_adcode", "true")
	params.Set("region", t.City)
	params.Set("query", t.KeyWord)
	//params.Set("ak", "nMNn6CW6pnnJxB3DIxKdIDLg3iw6RFso")
	params.Set("ak", config.BdKey)

	return params.Encode()
}
