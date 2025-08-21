package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var (
	wg      sync.WaitGroup
	totalCh = make(chan struct{})
	total   = 0
)

func main() {
	start := time.Now()
<<<<<<< HEAD
	limitCh := make(chan struct{}, 100)
	root := "D:\\"
=======
	root := "C:\\Windows"
>>>>>>> main

	go func() {
		for {
			select {
			case _, ok := <-totalCh:
				if !ok {
					return
				}
				total++
			default:
			}
		}
	}()

	Loop(root, limitCh, true)

	wg.Wait()

	close(totalCh)

	fmt.Printf("total = %d, time = %v/n", total, time.Since(start))

	var i = 10

	for i > 0 {
		fmt.Println(runtime.NumGoroutine())
		i--
		time.Sleep(time.Second)
	}

}

func Loop(root string, limit chan struct{}, f bool) {
	fd, err := os.ReadDir(root)
	if err == nil {
		for _, file := range fd {
			if !file.IsDir() {
				totalCh <- struct{}{}
			} else {
				select {
				case limit <- struct{}{}:
					wg.Add(1)
					go Loop(filepath.Join(root, file.Name()), limit, false)
				default:
					Loop(filepath.Join(root, file.Name()), limit, true)
				}
			}
		}
	}

	if !f {
		<-limit
		wg.Done()
	}
}
