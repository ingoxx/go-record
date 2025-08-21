package main

import (
	_ "embed"
	"fmt"
)

//go:embed test.txt
var data string

// 嵌入功能, 编译时, 默认不会把静态文件也一起编译进去, 现在可以使用embed来使其编译也能将静态文件加入
func main() {
	fmt.Println(data)
}
