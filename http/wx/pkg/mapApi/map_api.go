package mapApi

import (
	"github.com/ingoxx/go-record/http/wx/pkg/bdMapApi"
	"github.com/ingoxx/go-record/http/wx/pkg/form"
	"github.com/ingoxx/go-record/http/wx/pkg/gdMapApi"
	"github.com/ingoxx/go-record/http/wx/pkg/qqMapApi"
	"sync"
)

var (
	goroutineNum = 10
)

type MapApi struct {
	urls    map[string]string
	city    string
	keyWord string
	limit   chan struct{}
}

type MapData struct {
	Data    []*form.SaveInRedis
	Err     error
	Project string
}

// NewMapApi 新的api接口
func NewMapApi(city, keyWord string) *MapApi {
	urls := map[string]string{
		"gd": "https://restapi.amap.com",
		"qq": "https://apis.map.qq.com",
		//"bd": "https://api.map.baidu.com",
	}

	return &MapApi{
		urls:    urls,
		city:    city,
		keyWord: keyWord,
		limit:   make(chan struct{}, goroutineNum),
	}
}

func (ma *MapApi) Run() []*MapData {
	var wg sync.WaitGroup
	resultChan := make(chan *MapData, len(ma.urls))

	for v := range ma.urls {
		ma.limit <- struct{}{}
		wg.Add(1)

		go func(key string) {
			defer func() {
				wg.Done()
				<-ma.limit
			}()

			var data []*form.SaveInRedis
			var err error
			if key == "gd" {
				data, err = gdMapApi.NewGdMapApi(ma.urls[key], ma.city, ma.keyWord).KeyWordSearch()
				if err != nil {
					resultChan <- &MapData{Err: err, Project: key}
					return
				}
				resultChan <- &MapData{Data: data, Project: key}
			} else if key == "qq" {
				data, err = qqMapApi.NewTxMapApi(ma.urls[key], ma.city, ma.keyWord).KeyWordSearch()
				if err != nil {
					resultChan <- &MapData{Err: err, Project: key}
					return
				}
				resultChan <- &MapData{Data: data, Project: key}
			} else if key == "bd" {
				data, err = bdMapApi.NewBdMapApi(ma.urls[key], ma.city, ma.keyWord).KeyWordSearch()
				if err != nil {
					resultChan <- &MapData{Err: err, Project: key}
					return
				}
				resultChan <- &MapData{Data: data, Project: key}
			}
		}(v)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []*MapData
	for d := range resultChan {
		results = append(results, d)
	}

	return results
}

func (ma *MapApi) GetGdSinglePlaceSearch() ([]string, error) {
	return gdMapApi.NewGdMapApi(ma.urls["gd"], ma.city, ma.keyWord).SinglePlaceSearch()
}
