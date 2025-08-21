package main

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"
)

//搜索某个目录下的指定文件有多少个
var (
	Match          = 0
	MatchChan      = make(chan bool)
	MaxWorkersChan = make(chan string, 20)
	wg             sync.WaitGroup
)

func main() {
	start := time.Now()
	path := "C:/Users/"
	// path := "C:/Users/Administrator/Desktop/test/2"
	filename := "test.txt"
	go func() {
		for {
			_, ok := <-MatchChan
			if ok {
				Match++
			} else {
				break
			}
		}
	}()
	FindFile(path, filename, true)
	wg.Wait()
	close(MatchChan)
	close(MaxWorkersChan)
	fmt.Printf("count = %d, cost = %v\n", Match, time.Since(start))
}

func FindFile(path, filename string, s bool) {
	// fmt.Printf("time = %v, goroutine num = %d\n", time.Now().Unix(), runtime.NumGoroutine())
	if !s {
		<-MaxWorkersChan
	}
	fl, err := ioutil.ReadDir(path)
	if err == nil {
		for _, file := range fl {
			if file.Name() == filename {
				MatchChan <- true
			}
			if file.IsDir() {
				MaxWorkersChan <- path + file.Name() + "/" //这里超过容量，就会报deadlock!
				wg.Add(1)
				// time.Sleep(time.Second)
				go FindFile(path+file.Name()+"/", filename, false)
			}
		}
	}
	if !s {
		wg.Done()
	}
}
