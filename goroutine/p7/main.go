package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"time"
)

var (
	match = 0
	total = 0
)

func main() {

	start := time.Now()
	path := "C:/Users/Administrator/Desktop/test/"
	// path := "C:/Users/"
	FindFile(path, "test.txt")
	fmt.Printf("total = %d, count file = %d,cost time = %v\n", total, match, time.Since(start))
}

func FindFile(path, filename string) {
	fl, err := ioutil.ReadDir(path)
	if err == nil {
		for _, file := range fl {
			if file.IsDir() {
				FindFile(path+file.Name()+"/", filename)
			} else {
				// if file.Name() == filename {
				// 	match++
				// }
				// total++
				fmt.Printf("时间: %v, 文件名: %s, 协程数: %d\n", time.Now().Format("2006-01-02 15:04:05"), path+file.Name(), runtime.NumGoroutine())
			}

		}
	}
}
