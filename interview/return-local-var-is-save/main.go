package main

import "log"

func main() {
	// Go函数返回局部变量指针是否安全? 在golang里边是安全的
	log.Print(add(1, 2))
}

func add(x, y int) *int {
	res := x + y
	return &res
}

// 编译时go build -gcflags="-m -m -l" -o main.exe，可以看到是否发生逃逸的分析，或者go run -gcflags="-m -m -l" add.go
