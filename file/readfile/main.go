package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func main() {

	//create, err := os.Create("C:\\Users\\Administrator\\Desktop\\up.sh")
	//fmt.Println(create, err)

	err := os.Rename("C:\\Users\\Administrator\\Desktop\\up.sh", "C:\\Users\\Administrator\\Desktop\\up2.sh")
	fmt.Println(err)
}

func fileMd5(file string) (m5 string) {
	f, err := os.Open(file)
	if err != nil {
		return
	}

	defer f.Close()

	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		return
	}

	m5 = hex.EncodeToString(h.Sum(nil))

	return
}
