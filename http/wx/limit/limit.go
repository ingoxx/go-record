package main

import (
	"fmt"
	"github.com/importcjj/sensitive"
	"log"
)

func main() {
	filter := sensitive.New()
	if err := filter.LoadWordDict("D:\\project\\github.com\\ingoxx\\go-record\\http\\wx\\limit\\dict.txt"); err != nil {
		log.Fatalln("无法读取脏字库文件", err.Error())
	}

	fmt.Println(filter.Replace("你好我草", '*'))
}
