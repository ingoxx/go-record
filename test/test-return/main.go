package main

import (
	"fmt"
	"time"
)

func main() {
	i := 10
	for {
		if i == 4 {
			return // 如果这里return相当于退出了整个程序， 最后的fmt.Println("ok")是不会执行的
		}
		fmt.Println(i)
		time.Sleep(time.Second)
		i--
	}

	fmt.Println("ok")
}
