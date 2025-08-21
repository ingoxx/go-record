package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	total = 0
)

func sendwork(path string) {
	fl, err := os.ReadDir(path)
	if err == nil {
		for _, file := range fl {
			if file.IsDir() {
				// log.Print("path = ", filepath.Join(path, file.Name()))
				sendwork(filepath.Join(path, file.Name()))
			} else {
				total++
			}
		}
	}
}

func main() {
	start := time.Now()
	path := "C:/Windows/"
	// path := "C:/Users/Administrator/Desktop/test/"
	sendwork(path)
	log.Printf("total = %d, cost time = %v", total, time.Since(start))
}
