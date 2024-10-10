package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name string
}

func main() {
	s := &Student{"lxb"}
	getInterfaceVal(s)
}

func getInterfaceVal(i interface{}) {
	vo := reflect.ValueOf(i)
	for i := 0; i < vo.Elem().NumField(); i++ {
		fmt.Println(vo.Elem().Field(i))
	}
}
