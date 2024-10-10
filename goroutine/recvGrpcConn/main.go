package main

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type Task func(ctx context.Context)

type MultiWork struct {
	Works chan Task
	Limit chan struct{}
	Wg    sync.WaitGroup
}

func NewMultiWork(workers int) *MultiWork {
	nm := &MultiWork{
		Works: make(chan Task),
		Limit: make(chan struct{}, workers),
	}

	go func() {
		for task := range nm.Works {
			nm.Limit <- struct{}{}

			nm.Wg.Add(1)
			go func(task Task) {

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				task(ctx)

				<-nm.Limit
			}(task)

		}
	}()

	return nm
}

func main() {
	rand.Seed(time.Now().UnixNano())

	nm := NewMultiWork(10)

	task := func(ctx context.Context) {
		defer nm.Wg.Done()
		select {
		case <-time.After(time.Hour):
			fmt.Println("task finished ", rand.Intn(1000))
			return
		case <-ctx.Done():
			fmt.Println("task timeout ")
			return
		}
	}

	for range [100]struct{}{} {
		nm.Works <- task
	}

	nm.Wg.Wait()
	close(nm.Works)

	i := 25
	for i > 0 {
		fmt.Println("gn = ", runtime.NumGoroutine())
		i--
		time.Sleep(time.Second * 1)
	}

}
