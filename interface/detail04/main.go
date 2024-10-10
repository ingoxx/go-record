package main

//循环导包问题需要借助接口来解决
import (
	"github.com/Lxb921006/Golang-practise/interface/detail04/data"
	"github.com/Lxb921006/Golang-practise/interface/detail04/data2"
)

func main() {
	data2.RequireData()
	nd2 := data2.NewData2()
	data.RequireDta2(nd2)
}
