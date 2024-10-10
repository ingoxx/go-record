package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	src := "D:\\工作工具\\SQLServer2019-x64-CHS.iso"
	dst := "C:\\Users\\Administrator\\Desktop\\update\\SQLServer2019-x64-CHS.iso"

	r, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}

	nr := bufio.NewReader(r)

	w, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}

	defer w.Close()

	rb := make([]byte, 4096)

	nw := bufio.NewWriter(w)

	for {
		n, err1 := nr.Read(rb)
		if err1 == io.EOF {
			break
		}

		if err1 != nil {
			log.Fatal(err1)
		}

		nw.Write(rb[:n])
	}

	if err = nw.Flush(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("耗时 >>>", time.Since(start))

}
