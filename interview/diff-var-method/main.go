package main

import (
	"log"
	"time"
)

// 方法值接收者和指针接收者区别

type Person struct {
	Age int
}

func (p *Person) Add1() {
	p.Age += 10
}

func (p Person) Add2() {
	p.Age += 10
}

func (p Person) Get() int {
	return p.Age
}

func main() {
	p := Person{Age: 31}

	p.Add1()
	log.Print(p.Get())

	p.Add2()
	log.Print(p.Get())

	p2 := &Person{Age: 41}
	p2.Add1()
	log.Print(p2.Get())

	p2.Add2()
	log.Print(p2.Get())

	t2 := time.Now().Add(10 * time.Second)
	log.Print(t2.Format("2006-01-02 15:04:05"))

	t3 := <-time.After(10)
	log.Print(t3.Format("2006-01-02 15:04:05"))

	time.Sleep(time.Second * 2)
	t1 := time.Now()

	log.Print(t1.Format("2006-01-02 15:04:05"))

}
