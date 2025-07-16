package data

import (
	"fmt"

	"github.com/ingoxx/Golang-practise/interface/detail04/api"
)

type Data1 struct {
}

func (Data1) DataCheckOne() {
	fmt.Println("check one")
}

func NewData() *Data1 {
	return &Data1{}
}

// 如果data package需要调用data2的ApiPack()需要创建一个接口
func RequireDta2(api api.Api) {
	api.ApiPack()

}
