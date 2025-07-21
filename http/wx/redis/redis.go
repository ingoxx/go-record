package redis

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

var (
	rds       *redis.Client
	keyPrefix = "group-id"
	KeyGroups = "bk_list"
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
}

func NewRM() *RM {
	return &RM{}
}

func (r *RM) formatKey(key string) string {
	return fmt.Sprintf("%s-%s", keyPrefix, key)
}

func (r *RM) Set(key string, b interface{}) error {
	return rds.Set(r.formatKey(key), b, 0).Err()
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

func (r *RM) GetAllData() (string, error) {
	result, err := rds.Get(KeyGroups).Result()
	if err != nil {
		return result, err
	}

	if result == "" {
		return result, errors.New("null")
	}

	return result, nil
}
