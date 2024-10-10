package main

import (
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("https://192.168.3.20:6443", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(config)
}
