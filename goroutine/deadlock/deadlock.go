package main

import (
	"fmt"
	"time"
)

// 死锁例子:死锁（deadlock）通常发生在一组goroutine之间，它们都在等待一些事件，而这些事件只有通过这组goroutine中的某个goroutine才能发生。换句话说，死锁发生在一组goroutine相互等待对方做某事的情况
func main() {
	var workChan = make(chan int)
	go task(workChan)
	go task(workChan)
	go task(workChan)
	for i := 0; i < 4; i++ {
		workChan <- i
	}
	time.Sleep(time.Second * 5)
}

func task(d chan int) {

	fmt.Println(<-d)

}
