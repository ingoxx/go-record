// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	var done = make(chan int)
	task(done)
}

func task(done chan int) {
	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel1()
	go func() {
		run()
		done <- 1
	}()

	select {
	case <-ctx1.Done():
		fmt.Println("time out")
	case <-done:
		fmt.Println("ok")
	}

}

func run() {
	time.Sleep(time.Second * 5)
	fmt.Println("i am run()")
}
