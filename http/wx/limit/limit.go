package main

import (
	"fmt"
	"github.com/google/uuid"
)

func main() {
	//filter := sensitive.New()
	//if err := filter.LoadWordDict("D:\\project\\github.com\\ingoxx\\go-record\\http\\wx\\limit\\dict.txt"); err != nil {
	//	log.Fatalln("无法读取脏字库文件", err.Error())
	//}
	//
	//fmt.Println(filter.Replace("你好我草", '*'))

	u4 := uuid.New() // 或者 uuid.NewRandom()
	fmt.Println("UUID v4:", u4.String())
}
