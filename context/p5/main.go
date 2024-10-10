package main

import (
	"context"
	"fmt"
)

//<-chan int 表示只读管道
func gen(ctx context.Context) <-chan int {
	dst := make(chan int)
	n := 1
	go func() {
		for {
			select {
			case <-ctx.Done():
				return // return结束该goroutine，防止泄露
			case dst <- n:
				n++
			}
		}
	}()
	return dst
}

func main() {
	ctx, cancel := context.WithCancel(context.Background()) //最顶层的父上下文, 为了衍生出更多的子上下文，主要用于main，初始化，测试中
	defer cancel()                                          // 当我们取完需要的整数后调用cancel

	for n := range gen(ctx) {
		fmt.Println(n)
		// if n == 5 {
		// 	break
		// }
	}
}
