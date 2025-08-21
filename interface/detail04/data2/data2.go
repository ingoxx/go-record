package data2

import (
	"fmt"

	"github.com/ingoxx/Golang-practise/interface/detail04/data"
)

type Data2 struct {
	Name string
}

func (Data2) ApiPack() {
	fmt.Println("from Data2 ApiPack")
}

func NewData2() *Data2 {
	return &Data2{
		Name: "lxb",
	}
}

func RequireData() {
	data1 := data.NewData()
	data1.DataCheckOne()
}
