package redis

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

var (
	Rds       *redis.Client
	keyPrefix = "group-id"
)

func init() {
	var Rds = redis.NewClient(
		&redis.Options{
			Addr:         "193.112.111.237:6378",
			DB:           4,
			MinIdleConns: 5,
			Password:     "chatai",
			PoolSize:     5,
			PoolTimeout:  30 * time.Second,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	)

	if err := Rds.Ping(); err != nil {
		log.Fatalln("fail to connect redis, error msg: ", err)
	}

}

type RM struct {
	rds *redis.Client
}

func NewRM() *RM {
	return &RM{
		rds: Rds,
	}
}

func (r *RM) formatKey(key string) string {
	return fmt.Sprintf("%s-%s", keyPrefix, key)
}

func (r *RM) Set(key string, b interface{}) error {
	return r.rds.Set(r.formatKey(key), b, 0).Err()
}

func (r *RM) Get(key string) (string, error) {
	result, err := r.rds.Get(r.formatKey(key)).Result()
	if err != nil {
		return result, err
	}

	if result == "" {
		return result, errors.New("null")
	}

	return result, nil
}
