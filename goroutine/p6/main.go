package main

import "fmt"

//注意：当没有数据再传入管道时，应该关闭channel，或者是退出主线程
func main() {
	i := 0
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)
	// 开启goroutine将0~100的数发送到ch1中
	go func() {
		// defer close(ch1)
		for {

			i++
			ch1 <- i
		}
	}()
	// 开启goroutine从ch1中接收值，并将该值的平方发送到ch2中
	go func() {
		for {
			i, ok := <-ch1 // 通道关闭后再取值ok=false
			if !ok {
				break
			}
			ch2 <- i * i
		}
		// close(ch2)
	}()
	// 在主goroutine中从ch2中接收值打印
	// for i := range ch2 { // 通道关闭后会退出for range循环
	// 	fmt.Println(i)
	// }
	for {
		fmt.Println(<-ch1)
		fmt.Println(<-ch2)
	}
}
