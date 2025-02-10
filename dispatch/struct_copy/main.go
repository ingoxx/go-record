package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type StdA struct {
	Name   string
	Age    uint
	Gender uint
}

type StdB struct {
	Name   string
	Age    uint
	Gender uint
	Addr   string // 目标结构体多出的字段
}

func main() {
	structToMap()
}

func structToStruct() {
	sa := StdA{Name: "jack", Age: 18, Gender: 1}
	var sb StdB
	sb.Addr = "default address" // Addr 有初始值

	err := mapstructure.Decode(sa, &sb)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Result: %+v\n", sb)
}

func structToMap() {
	var m = make([]map[string]interface{}, 0)
	s := []struct {
		Name string
		Age  uint
	}{
		{
			Name: "lxb",
			Age:  18,
		},
		{
			Name: "lqm",
			Age:  18,
		},
	}

	b, _ := json.Marshal(&s)
	json.Unmarshal(b, &m)
	fmt.Println(m)
}
