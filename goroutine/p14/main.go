package main

import (
	"fmt"
	"runtime"
	"sync"
)

var (
	wg sync.WaitGroup
)

type Pool struct {
	Work chan func()
	Sem  chan bool //有缓冲的chan，限制goroutine number
}

func (p *Pool) Worker(task func()) {
	// fmt.Println("gn00000000000 = ", runtime.NumGoroutine())
	defer func() {
		<-p.Sem
		wg.Done()
	}()
	// for task := range p.Work {
	// 	task = <-p.Work
	// 	task()

	// }

	// wg.Done()
	for {
		task()
		task = <-p.Work

	}
}

func (p *Pool) Task(task func()) {
	// fmt.Println("gn33333333333 = ", runtime.NumGoroutine())
	select {
	case p.Work <- task:
	case p.Sem <- true:
		go p.Worker(task)
	}
}

func NewPool(size int) *Pool {
	return &Pool{
		Work: make(chan func()),
		Sem:  make(chan bool, size),
	}
}

func main() {
	pool := NewPool(50)
	for i := 0; i < 1000; i++ {
		pool.Task(func() {
			defer wg.Done()
			// time.Sleep(time.Second)
			fmt.Printf("goroutine num = %d\n", runtime.NumGoroutine())
		})
		// fmt.Println("gn111111111111111 = ", runtime.NumGoroutine())
	}
	wg.Wait()
	fmt.Println("finished")
	// time.Sleep(time.Second)

}
