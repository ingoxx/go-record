package main

import (
	"context"
	"log"
	"math/rand"
	"time"
)

// 在一些请求-响应场景中，由于各种原因，一个请求可能需要很长时间才能响应，有时甚至永远不会响应。对于这种情况，我们应该使用超时解决方案向客户端返回错误信息。这样的超时解决方案可以用select机制来实现

func requestWithTimeout(timeout time.Duration) {
	// c := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	// May need a long time to get the response.
	for range [5]struct{}{} {
		go doRequest(rand.Intn(100))
		select {
		case <-ctx.Done():
			log.Print("timeout111")
		case <-time.After(timeout * time.Second):
			log.Print("timeout222")
		}
	}

}

func doRequest(d int) {
	log.Print("start")
	time.Sleep(3 * time.Second)
}
func main() {
	rand.Seed(time.Now().UnixNano())

	requestWithTimeout(2)

	time.Sleep(30 * time.Second)

}
