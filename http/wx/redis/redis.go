package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ingoxx/go-record/http/wx/form"
	"github.com/ingoxx/go-record/http/wx/qqMapApi"
	"log"
	"sync"
	"time"
)

var (
	rds         *redis.Client
	groupKey    = "group-id"
	AddrListKey = "addr-check-list"
	//cusAddrKey   = "group-id-cus"
)

func init() {
	rds = redis.NewClient(
		&redis.Options{
			Addr:         "193.112.111.237:6378",
			DB:           1,
			MinIdleConns: 5,
			Password:     "chatai",
			PoolSize:     5,
			PoolTimeout:  30 * time.Second,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	)

	if err := rds.Ping().Err(); err != nil {
		log.Println("fail to connect redis, error msg: ", err)
	}

	log.Println("redis init completed")

}

type RM struct {
	mu sync.Mutex
}

func NewRM() *RM {
	return &RM{}
}

func (r *RM) formatKey(key string) string {
	return fmt.Sprintf("%s-%s", groupKey, key)
}

func (r *RM) Set(key string, b interface{}, expire time.Duration) error {
	return rds.Set(r.formatKey(key), b, expire).Err()
}

func (r *RM) Get(key string) (string, error) {
	result, err := rds.Get(r.formatKey(key)).Result()
	if err != nil {
		return result, err
	}

	if result == "" {
		return result, errors.New("null")
	}

	return result, nil
}

// GetAllData 当前市的所有篮球场地址, 只保留两周, 两周后重新更新
func (r *RM) GetAllData(key string, cnKey string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result, err := r.Get(key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return result, err
	}

	if result == "" {
		search, err := qqMapApi.NewTxMapApi(cnKey).KeyWordSearch()
		if err != nil {
			return result, err
		}
		ld, err := r.mergeData(key)
		if err != nil {
			return result, err
		}
		if len(ld) > 0 {
			search = append(search, ld...)
		}
		b, err := json.Marshal(&search)
		if err != nil {
			return result, err
		}
		if err := r.Set(key, b, time.Second*time.Duration(1209600)); err != nil {
			return result, err
		}

		return string(b), nil
	}

	return result, nil
}

func (r *RM) mergeData(key string) ([]qqMapApi.SaveInRedis, error) {
	var dataList = make([]qqMapApi.SaveInRedis, 0)
	list, err := r.GetAddrList() // 遍历获取审核列表，找到对应id将其更新到指定key的数据中
	if err != nil {
		return dataList, err
	}

	dataList = make([]qqMapApi.SaveInRedis, 0, len(list))
	for _, data := range list {
		if data.CityPy == key {
			ad := qqMapApi.SaveInRedis{
				Id:     data.Id,
				Tags:   []string{"篮球场"},
				Img:    "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/bk3.svg",
				Addr:   data.Addr,
				Lat:    data.Lat,
				Lng:    data.Lng,
				UserId: data.UserId,
				Title:  "篮球场",
			}
			dataList = append(dataList, ad)
		}
	}

	return dataList, nil
}

func (r *RM) Update(key, id string) ([]qqMapApi.SaveInRedis, error) {
	var dataList []qqMapApi.SaveInRedis
	result, err := r.Get(key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return dataList, err
	}

	if err := json.Unmarshal([]byte(result), &dataList); err != nil {
		return dataList, err
	}

	list, err := r.GetAddrList() // 遍历获取审核列表，找到对应id将其更新到指定key的数据中
	if err != nil {
		return dataList, err
	}

	for _, data := range list {
		if data.Id == id && !data.IsRecord && data.CityPy == key {
			ad := qqMapApi.SaveInRedis{
				Id:     data.Id,
				Tags:   []string{"篮球场"},
				Img:    "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/bk3.svg",
				Addr:   data.Addr,
				Lat:    data.Lat,
				Lng:    data.Lng,
				UserId: data.UserId,
				Title:  "篮球场",
			}
			dataList = append(dataList, ad)
			break
		}
	}

	b, err := json.Marshal(dataList)
	if err := json.Unmarshal([]byte(result), &dataList); err != nil {
		return dataList, err
	}

	if err := r.Set(key, b, time.Second*time.Duration(1209600)); err != nil {
		return dataList, err
	}

	if _, err := r.UpdateAddrList(id); err != nil {
		return dataList, err
	}

	return dataList, nil
}

// GetAddrList 所有用户添加的篮球场地址列表，不过期长期保存用户添加的篮球场地址
func (r *RM) GetAddrList() ([]*form.AddAddrForm, error) {
	var dataList []*form.AddAddrForm
	result, err := r.Get(AddrListKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return dataList, err
	}

	if errors.Is(err, redis.Nil) {
		dataList = make([]*form.AddAddrForm, 0)
		b, err := json.Marshal(&dataList)
		if err != nil {
			return dataList, err
		}
		if err := r.Set(AddrListKey, b, 0); err != nil {
			return dataList, err
		}
		return dataList, nil
	}

	if err := json.Unmarshal([]byte(result), &dataList); err != nil {
		return dataList, err
	}

	return dataList, nil
}

// UserAddAddrReq 用户提交添加篮球场地址的请求
func (r *RM) UserAddAddrReq(data form.AddAddrForm) error {
	var dataList = make([]form.AddAddrForm, 0)
	result, err := r.Get(AddrListKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if result != "" {
		if err := json.Unmarshal([]byte(result), &dataList); err != nil {
			return err
		}
	}

	dataList = append(dataList, data)
	b, err := json.Marshal(&dataList)
	if err != nil {
		return err
	}

	if err := r.Set(AddrListKey, b, 0); err != nil {
		return err
	}

	return nil
}

// UpdateAddrList 更新审核列表
func (r *RM) UpdateAddrList(id string) ([]*form.AddAddrForm, error) {
	list, err := r.GetAddrList() // 遍历获取审核列表，找到对应id将其更新到指定key的数据中
	if err != nil {
		return list, err
	}

	for _, v := range list {
		if v.Id == id {
			v.IsRecord = true
		}
	}

	b, err := json.Marshal(&list)
	if err != nil {
		return list, err
	}

	if err := r.Set(AddrListKey, b, 0); err != nil {
		return list, err
	}

	return list, nil
}
