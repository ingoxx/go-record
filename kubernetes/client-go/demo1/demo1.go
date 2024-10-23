package main

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	scheme2 "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func breakPoint() {
	fmt.Println("aaa")
}

func main() {
	breakPoint()
	//config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatal("config err >>> ", err)
	}

	//client
	config.GroupVersion = &v1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme2.Codecs
	config.APIPath = "/api"
	clientFor, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatal("clientFor err >>> ", err)
	}

	// get data
	pod := v1.Pod{}
	err = clientFor.Get().Namespace("web").Resource("pods").Name("nginx-deployment-74bdddb69b-7nljw").Do(context.Background()).Into(&pod)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("pod sources >>> ", pod)

}
