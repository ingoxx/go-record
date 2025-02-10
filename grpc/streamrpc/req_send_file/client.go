package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	addr := flag.String("addr", "", "rpc服务端地址端口") // 定义一个字符串类型的标志
	file := flag.String("file", "30", "文件路径")     // 定义一个整数类型的标志
	flag.Parse()

	if *addr == "" || *file == "" {
		_, err := fmt.Fprintln(os.Stderr, "Error: -addr跟-file都是必须参数，不能为空")
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
