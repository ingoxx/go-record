package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("开始计时...")

	timer := time.AfterFunc(3*time.Second, func() {
		fmt.Println("3秒后触发的任务")
	})

	defer timer.Stop()

	time.Sleep(5 * time.Second) // 确保主协程存活足够长的时间
}
