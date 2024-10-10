package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(2))
	defer cancel()

	//go func() {
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			fmt.Println(ctx.Err())
	//			return
	//		}
	//	}
	//}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println(ctx.Err())
					return
				}
			}
		}()

		resp := task(ctx)
		fmt.Println(resp)
	}()

	wg.Wait()

}

func task(ctx context.Context) error {
	time.Sleep(time.Second * time.Duration(14))
	//if ctx != nil {
	//	select {
	//	case <-ctx.Done():
	//		return ctx.Err()
	//	default:
	//	}
	//}

	return nil

}
