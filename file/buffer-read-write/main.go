package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

// 一般是一次读取4K或者8K的数据到缓冲区
func main() {

	start := time.Now()

	var rb = make([]byte, 4096)
	file := "D:\\工作工具\\SQLServer2019-x64-CHS.iso"
	path := "C:\\Users\\Administrator\\Desktop\\SQLServer2019-x64-CHS.iso"

	f, err := os.Open(file)
	if err != nil {
		return
	}

	defer f.Close()

	fn, err := os.Create(path)
	if err != nil {
		return
	}

	for {
		n, err1 := f.Read(rb)
		if err1 == io.EOF {
			break
		}

		if err1 != nil {
			return
		}

		_, err = fn.Write(rb[:n])
		if err != nil {
			return
		}
	}

	fmt.Println("耗时 = ", time.Since(start))
}
