package main

import (
	"fmt"
	"os"
	"time"
)

var (
	Match = 0
	Total = 0
)

func main() {
	path := "/Users"
	// path := "C:/Users/Administrator/Desktop/test/2"
	filename := "test.txt"
	start := time.Now()
	Find(path, filename)
	fmt.Printf("total = %d, match = %d, cost time = %v\n", Total, Match, time.Since(start))
}

func Find(path, filename string) {
	fl, err := os.ReadDir(path)
	if err == nil {
		for _, file := range fl {
			if file.Name() == filename {
				Match++
			}
			if file.IsDir() {
				Find(path+file.Name()+"/", filename)
			} else {
				Total++
			}
		}
	}
}
