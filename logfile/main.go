package main

import (
	"github.com/ingoxx/Golang-practise/logger"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("C:\\Users\\Administrator\\Desktop\\my.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("errors opening file: %v", err)
	}

	defer f.Close()

	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	log.Println(" Orders API Called")

	fi, err := os.Stat("C:\\Users\\Administrator\\Desktop\\my.log")
	if err != nil {
		return
	}

	log.Println("size >>>,", fi.Size(), fi.Name())

	logger.SetLogFile("C:\\Users\\Administrator\\Desktop\\my2.log")
	logger.SetLogLevel(logger.InfoLevel)
	logger.Info("test logger %s", "lxb")
	logger.Debug("test logger %s", "lqm")
	logger.Error("test logger %s", "lyy")
}
