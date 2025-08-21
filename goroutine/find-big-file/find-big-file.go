package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const (
	//mb
	size = 500
)

var (
	wg      sync.WaitGroup
	totalCh = make(chan string)
	total   = 0
)

func main() {
	start := time.Now()

	defer func() {
		fmt.Printf("文件大于%dm总共有: %d, 耗时: %v\n", size, total, time.Since(start))
	}()

	limitCh := make(chan struct{}, runtime.NumCPU()/2)

	root := "D:\\"

	go func() {
		for {
			select {
			case file, ok := <-totalCh:
				if !ok {
					return
				}
				fmt.Printf("文件路径: %s, 当前协程数量: %d\n", file, runtime.NumGoroutine())
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
				s, err := os.Stat(filepath.Join(root, file.Name()))
				if err == nil {
					if s.Size()/1024/1024 > size {
						totalCh <- filepath.Join(root, file.Name())
					}
				}
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
