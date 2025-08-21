// You can edit this code!
// Click here and start typing.
package main

import (
	"log"
	"time"
)

func main() {
	data := make(chan error, 1)
	exit := make(chan struct{})
	go doReq(data)
	go doResp(data, exit)
	<-exit
}

func doReq(data chan error) {
	log.Println("begin")
	time.Sleep(time.Second * 2)
	log.Println("finished")
	data <- nil

}

func doResp(data chan error, exit chan struct{}) {
	for {
		select {
		case <-data:
			log.Println("resp ok")
			exit <- struct{}{}
			return
		default:
		}
	}
}
