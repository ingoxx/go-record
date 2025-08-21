package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// 一般是一次读取4K或者8K的数据到缓冲区
func main() {
	start := time.Now()
	src := "D:\\工作工具\\SQLServer2019-x64-CHS.iso"
	dst := "C:\\Users\\Administrator\\Desktop\\update\\SQLServer2019-x64-CHS.iso"
	//bufferWrite(src, dst)
	normalWrite(src, dst)

	fmt.Println("done >>>", time.Since(start))
}

func bufferWrite(file, path string) {
	var rb = make([]byte, 4096)

	fn, err := os.Create(path)
	if err != nil {
		return
	}

	f, err := os.Open(file)
	if err != nil {
		return
	}

	defer f.Close()

	for {
		n, err := f.Read(rb)
		if err == io.EOF {
			break
		}

		if err != nil {
			return
		}

		_, err = fn.Write(rb[:n])
		if err != nil {
			return
		}
	}
}

func normalWrite(file, path string) {
	rn, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	wn, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = io.Copy(wn, rn); err != nil {
		log.Fatal("copy err >>> ", err)
	}
}
