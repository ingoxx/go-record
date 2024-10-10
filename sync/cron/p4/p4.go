package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("removed from queue = ", len(queue))
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			// Wait does not block - it suspends the main goroutine until a signal on the condition has been sent.
			fmt.Println("等待队列数据被消费")
			c.Wait()
		}
		fmt.Println("adding to queue:", i)
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()
	}

	fmt.Println("queue len = ", len(queue))
}
