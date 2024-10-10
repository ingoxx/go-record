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
	limitCh = make(chan struct{}, runtime.NumCPU())
	total   = 0
)

func main() {
	start := time.Now()

	defer func() {
		fmt.Printf("total = %d, time = %v\n", total, time.Since(start))
	}()

	root := "C:\\Windows"

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
		wg.Done()
		<-limit
	}
}
