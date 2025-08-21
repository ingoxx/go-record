package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type Ps struct {
	Age  uint
	Name string
}

func main() {
	p := Ps{
		Age:  10,
		Name: "lll",
	}

	b, _ := json.Marshal(&p)
	var data map[string]interface{}
	json.Unmarshal(b, &data)
	fmt.Println(data)

	m, _ := Dispatch(p)
	fmt.Println(m)
}

func Dispatch(v interface{}) (m map[string]interface{}, err error) {
	m1 := make(map[string]interface{})
	t := reflect.TypeOf(v)
	vd := reflect.ValueOf(v)

	if t.Kind() != reflect.Struct {
		err = errors.New("不是期待的类型")
		return
	}

	for i := 0; i < t.NumField(); i++ {
		m1[t.Field(i).Name] = vd.Field(i).Interface()
	}

	m = m1

	return
}
