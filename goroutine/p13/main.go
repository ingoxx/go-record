package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
	// limit = make(chan bool, 20)
	totalchan = make(chan bool)
	total     = 0
)

func Work(path string, size int) {

	// fmt.Printf("gn1111: %d\n", runtime.NumGoroutine())

	defer func() {
		wg.Done()

	}()

	fl, err := os.ReadDir(path)
	if err == nil {
		for _, file := range fl {
			if file.IsDir() {
				wg.Add(1)
				go Work(path+file.Name()+"/", size)
			} else {
				// fmt.Printf("gn2222: %d\n", runtime.NumGoroutine())
				totalchan <- true
			}
		}
	}
}

func count() {
	for {
		_, ok := <-totalchan
		if ok {
			total++
		} else {
			break
		}
	}
}

func main() {
	path := "C:/Windows/"
	size := 0
	start := time.Now()
	wg.Add(1)
	go count()
	wg.Add(1)
	go Work(path, size)
	wg.Wait()
	fmt.Printf("total= %d, cost: %v\n", total, time.Since(start))
}
