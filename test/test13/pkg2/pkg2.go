package pkg2

import (
	"fmt"

	"github.com/ingoxx/Golang-practise/test/test13/pkg1"
)

type Pkg2 struct {
}

func (*Pkg2) FromPkg2() {
	fmt.Println("Pkg2")
	pk1 := pkg1.Pkg1{}
	pk1.FromPkg1()
}
