package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	forever := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				forever <- struct{}{}
				log.Print(ctx.Err().Error())
				return
			default:
				fmt.Println("for loop")
				time.Sleep(time.Second * 1)
			}

			time.Sleep(500 * time.Millisecond)
		}
	}(ctx)

	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()

	<-forever
	fmt.Println("finish")
}
