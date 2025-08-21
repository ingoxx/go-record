package gdMapApi

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/ingoxx/go-record/http/wx/pkg/config"
	"github.com/ingoxx/go-record/http/wx/pkg/form"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GdMapApi struct {
	City    string
	KeyWord string
	Url     string
}

func NewGdMapApi(url, city, keyWord string) GdMapApi {
	return GdMapApi{
		City:    city,
		KeyWord: keyWord,
		Url:     url,
	}
}

func (t GdMapApi) KeyWordSearch() ([]*form.SaveInRedis, error) {
	var fd form.GdAPIResponse
	var sd = make([]*form.SaveInRedis, 0, 100)
	offset := 1

	for offset <= 23 {
		url := fmt.Sprintf("%s/v3/place/text?key=%s&city=%s&keywords=%s&types=&children=1&offset=20&page=%d&extensions=all", t.Url, config.GdKey, t.City, t.KeyWord, offset)
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

		if len(fd.Data) == 0 {
			offset++
			continue
		}

		for _, v := range fd.Data {
			loc, ok := v.Location.(string)
			var lat, lng float64
			if ok {
				ls := strings.Split(loc, ",")
				lat, err = strconv.ParseFloat(ls[1], 64)
				if err != nil {
					log.Fatalln(err)
				}
				lng, err = strconv.ParseFloat(ls[0], 64)
				if err != nil {
					log.Fatalln(err)
				}
			}
			addr := fmt.Sprintf("%s%s%v%s", v.CityName, v.AdName, v.BusinessArea, v.Address)
			var img = make([]string, 0)
			var mainImg string

			for _, v := range v.Photos {
				if v.URL != "" {
					httpsUrl := strings.ReplaceAll(v.URL, "http:", "https:")
					img = append(img, httpsUrl)
				}
			}
			if len(img) > 0 {
				mainImg = img[0]
			}
			u4 := uuid.New()
			sdd := &form.SaveInRedis{
				Id:     u4.String(),
				UserId: config.Admin,
				Addr:   addr,
				Lat:    lat,
				Lng:    lng,
				Title:  v.Name,
				Tags:   []string{v.Name},
				Img:    mainImg,
				Aid:    v.ID,
				Images: img,
			}
			sd = append(sd, sdd)
		}

		offset++
		time.Sleep(time.Second)
	}

	return sd, nil
}
