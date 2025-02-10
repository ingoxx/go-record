package main

import "fmt"

// 如果有一个或多个非阻塞case操作，Go运行时会随机选择其中一个非阻塞操作执行，然后继续执行相应的case分支

func main() {
	c := make(chan struct{})
	close(c)
	// fmt.Println(<-c)
	// c <- struct_copy{}{}
	//select-case的分支是如果都是非阻塞则随机选择一个
	select {
	// Panic if the first case is selected.
	case c <- struct{}{}:
		fmt.Println("send")
	case <-c:
		fmt.Println("recv")
	}
}
