package pkg1

import "fmt"

type Pkg1 struct{}

func (*Pkg1) FromPkg1() {
	fmt.Println("pkg1")
}
