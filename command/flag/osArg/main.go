package main

import (
	"flag"
	"fmt"
)

var (
	iniFile  = flag.String("i", "", "ini file path")
	iniFile2 = flag.String("t", "", "ini file path")
)

func main() {
	flag.Parse()

	fmt.Println(*iniFile, *iniFile2)
	fmt.Println(flag.NArg())

}
