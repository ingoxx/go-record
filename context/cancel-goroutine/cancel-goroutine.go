package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	var in = make(chan int)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t := time.NewTimer(0)
	go func(ctx context.Context, data <-chan int) {
		for {
			select {
			case <-t.C:
				fmt.Println("timeout")
			case v := <-data:
				if v == 1 {
					t.Reset(time.Second * 2)
				}
				fmt.Println(v)
			case <-ctx.Done():
				return
			}
		}
	}(ctx, in)

	i := 0

	for {
		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
		if i == 6 {
			break
		}
		in <- i
		i++
	}

}
