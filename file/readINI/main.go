package main

import (
	"fmt"
	"gopkg.in/gcfg.v1"
)

type HuawenConfig struct {
	Key        string
	Secret     string
	Region     string
	PutSrcPath string
	Bucket     string
}

type Config struct {
	HuaWen HuawenConfig
}

func main() {
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, "C:\\Users\\Administrator\\Desktop\\update\\aws.ini")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(cfg.HuaWen.Key)
		fmt.Println(cfg.HuaWen.Secret)
		fmt.Println(cfg.HuaWen.PutSrcPath)
	}
}
