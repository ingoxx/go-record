package main

import (
	"fmt"
	"time"
)

// 在协程里, channel的收发都会阻塞代码运行, 解决方法go提供了select, 这样就可以channel解决收发阻塞代码运行的问题
func main() {
	c1 := make(chan string)
	c2 := make(chan string)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				c1 <- "lxb"
				time.Sleep(time.Second / 2)
			}
		}()
	}

	go func() {
		for {
			d := <-c1
			c2 <- d + "lqm"
			time.Sleep(time.Second * 2)
		}
	}()

	for {
		select {
		case c11 := <-c1:
			// time.Sleep(time.Second * 5)
			fmt.Println(c11)
		case c22 := <-c2:
			fmt.Println(c22)
		}
	}
	// for {
	// 	fmt.Println(<-c1)
	// 	fmt.Println(<-c2)
	// }
	// 当第一次获取c1,c2的值,由于c2管道在写入数据后要等待2秒，此时c1已经写入数据，但是没有还没有被获取，所以也会跟着阻塞到2秒，为了不影响c1的运行需要用到select
}
