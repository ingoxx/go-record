package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// 文件路径
	filePath := "C:\\Users\\Administrator\\Desktop\\aws_create_cloudfront.sh"

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 使用 bufio.Reader 逐行读取
	reader := bufio.NewReader(file)
	lineNumber := 1
	for {
		line, err := reader.ReadString('\n') // 读取直到遇到换行符
		if err != nil {
			if err.Error() == "EOF" { // 到达文件末尾
				break
			}
			log.Fatalf("读取文件时出错: %v", err)
		}
		// 输出当前读取的行内容
		if strings.TrimSpace(line) == "done" {
			fmt.Printf("%s\n", line)
		}
		//fmt.Printf("%s\n", line)
		lineNumber++
	}
}
