package main

import (
	"fmt"

	"github.com/ingoxx/Golang-practise/test/test13/pkg1"
	"github.com/ingoxx/Golang-practise/test/test13/pkg2"
)

type Allpkg struct {
	Pkg1 *pkg1.Pkg1
	Pkg2 *pkg2.Pkg2
	Name string
}

func main() {
	all := Allpkg{}
	all.Name = "lxb"
	all.Pkg1.FromPkg1()
	all.Pkg2.FromPkg2()
	fmt.Println("name = ", all.Name)
}
