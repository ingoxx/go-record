package main

import (
	"fmt"
	"sync"
)

const bufferSize = 5

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	buffer := make([]int, 0, bufferSize)

	wg.Add(2)

	// 生产者
	go func() {
		defer wg.Done()

		for i := 0; i < 10; i++ {
			mu.Lock()
			for len(buffer) == bufferSize {
				fmt.Println("缓冲区已满，等待消费者消费数据")
				cond.Wait() // 缓冲区已满，等待消费者消费数据
			}
			buffer = append(buffer, i)
			fmt.Println("Producer: produced", i)
			cond.Signal() // 唤醒消费者
			mu.Unlock()
		}
	}()

	// 消费者
	go func() {
		defer wg.Done()

		for i := 0; i < 10; i++ {
			mu.Lock()
			for len(buffer) == 0 {
				fmt.Println("缓冲区为空，等待生产者生产数据")
				cond.Wait() // 缓冲区为空，等待生产者生产数据
			}
			value := buffer[0]
			buffer = buffer[1:]
			fmt.Println("Consumer: consumed", value)
			cond.Signal() // 唤醒生产者
			mu.Unlock()
		}
	}()

	wg.Wait()
}
