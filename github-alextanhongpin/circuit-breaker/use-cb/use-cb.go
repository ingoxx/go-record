package main

import (
	"errors"
	"fmt"
	"github.com/Lxb921006/Golang-practise/github-alextanhongpin/circuit-breaker/cb"

	"time"
)

func main() {
	g := cb.New()
	defer g.Stop()

	for i := 0; i < 5; i++ {
		err := g.Do(func() error {
			return errors.New("error")
		})
		fmt.Println(err)
	}

	time.Sleep(1100 * time.Millisecond)

	for i := 0; i < 5; i++ {
		err := g.Do(func() error {
			return nil
		})
		fmt.Println(err)
	}

}
