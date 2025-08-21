package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {

	context.WithCancel(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	data, err := ReadFileContext(ctx, "D:\\工作工具\\SQLServer2019-x64-CHS.iso")
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			// 读取文件超时
			fmt.Println("ctx err >>> ", err)
		} else if errors.Is(err, os.ErrClosed) {
			fmt.Println("file close err >>> ", err)
			// 读取文件被取消
		} else {
			fmt.Println("unknown err >>> ", err)
			// 其他错误
		}
		return
	}

	// 处理读取的数据
	fmt.Println(string(data))

	//i := 25
	//for i > 0 {
	//	fmt.Println("gn = ", runtime.NumGoroutine())
	//	i--
	//	time.Sleep(time.Second * 1)
	//}
}

func ReadFileContext(ctx context.Context, filename string) ([]byte, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建一个 buffer
	buf := bytes.NewBuffer(nil)

	// 启动一个 goroutine，在 ctx.Done() 事件发生时取消读取
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			file.Close()
			close(done)
		case <-done:
		}
	}()

	// 从文件中读取数据，并写入 buffer
	_, err = io.Copy(buf, file)
	if err != nil {
		fmt.Println("copy err >>> ", err)
		return nil, err
	}

	// 检查是否被取消
	select {
	case <-done:
		fmt.Println("done err")
		return nil, ctx.Err()
	default:
	}

	// 返回数据
	return buf.Bytes(), nil
}
