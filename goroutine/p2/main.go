package main

import (
	"fmt"
	"os"
	"time"
)

// 搜索某个目录下的指定名字文件有多少个
type FindFiles struct {
	Path           string
	FileName       string
	Workers        int
	MaxWorkers     int
	Match          int
	Total          int
	SearchChan     chan string
	WorkerDoneChan chan bool
	ResChan        chan bool
	TotalChan      chan bool
}

func (ff *FindFiles) Ergodic(path string, signle bool) {
	// fmt.Println("gn = ", runtime.NumGoroutine())
	fl, err := os.ReadDir(path)
	if err == nil {
		for _, file := range fl {
			// if file.Name() == ff.FileName {
			// 	ff.ResChan <- true
			// }
			if file.IsDir() {
				if ff.Workers < ff.MaxWorkers {
					ff.SearchChan <- path + file.Name() + "/"
				} else {
					ff.Ergodic(path+file.Name()+"/", false)
				}
			} else {
				// if runtime.NumGoroutine() == 22 {
				// 	fmt.Println(runtime.NumGoroutine())
				// }
				ff.TotalChan <- true
			}
		}
	}
	if signle {
		ff.WorkerDoneChan <- true
	}
}

func (ff *FindFiles) Run() {
	go ff.Ergodic(ff.Path, true)

	for {
		select {
		case path := <-ff.SearchChan:
			ff.Workers++
			go ff.Ergodic(path, true)
		case <-ff.WorkerDoneChan:
			ff.Workers--
			if ff.Workers == 0 {
				close(ff.SearchChan)
				close(ff.WorkerDoneChan)
				close(ff.ResChan)
				return
			}
		case <-ff.TotalChan:
			ff.Total++
		case <-ff.ResChan:
			ff.Match++
		}
	}

}

func NewFindFiles(path, filename string) *FindFiles {
	return &FindFiles{
		Path:           path,
		FileName:       filename,
		Workers:        1,
		MaxWorkers:     200,
		SearchChan:     make(chan string),
		WorkerDoneChan: make(chan bool),
		ResChan:        make(chan bool),
		TotalChan:      make(chan bool),
	}
}

func main() {
	//统计test目录下文件名为test.txt的文件数量, 以及耗时
	start := time.Now()
	ff := NewFindFiles("C:\\Windows", "test.txt")
	// ff := NewFindFiles("C:/Users/Administrator/Desktop/test2/", "test.txt")
	ff.Run()
	fmt.Printf("total=%d, count file = %d,cost time = %v\n", ff.Total, ff.Match, time.Since(start))
}
