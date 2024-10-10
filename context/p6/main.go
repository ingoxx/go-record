package main

import (
	"context"
	"log"
	"math/rand"
	"time"
)

var (
	res = make(chan int)
)

func main() {
	do := make(chan int)

	rand.Seed(time.Now().Unix())
	s := []int{100, 200, 300, 400, 500, 600}
	ctx := context.Background()

	work := func() {
		for v := range do {
			run(v, ctx)
		}
	}

	for range [3]struct{}{} {
		go work()
	}

	for _, v := range s {
		do <- v
	}

	time.Sleep(time.Second * 10)

}

func run(v int, ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	//这里会有内存泄露问题
	// go func(ctx context.Context) {
	// 	for {
	// 		time.Sleep(time.Second)
	// 	}
	// 	res <- v
	// }(ctx)

	select {
	case <-ctx.Done():
		log.Print(v)
		log.Print("TIME OUT111")
		return
	case <-time.After(time.Second * 2):
		log.Print("TIME OUT222")
		return
	case v2 := <-res:
		log.Print(v2)
	}

}
