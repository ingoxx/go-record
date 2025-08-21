package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	taskCount := 5

	wg.Add(taskCount)

	// 任务执行函数
	task := func(id int) {
		defer wg.Done()

		mu.Lock()
		fmt.Println("Task", id, "ready, after 3s scheduling tasks...")
		cond.Wait() // 等待调度信号
		fmt.Println("Task", id, "finished")
		mu.Unlock()
	}

	// 创建并启动任务
	for i := 0; i < taskCount; i++ {
		go task(i)
	}

	// 调度任务
	time.Sleep(time.Second * 3)
	mu.Lock()
	cond.Broadcast() // 发送调度信号，唤醒所有等待的任务
	mu.Unlock()

	wg.Wait()
	fmt.Println("All tasks finished")
}
