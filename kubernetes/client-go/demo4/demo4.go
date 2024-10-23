package main

import (
	"context"
	"fmt"
	"github.com/imdario/mergo"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {

	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		inClusterConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalln("config err >>>", err)
		}
		config = inClusterConfig
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("clientSet err >>>", err)
	}

	list, err := clientSet.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range list.Items {
		fmt.Println(v.Name)
	}

	//informers.NewSharedInformerFactory(clientSet, time.Second*60)
	stur1 := struct {
		Name string
		Age  int
	}{
		Name: "lxb",
		Age:  32,
	}

	data := map[string]interface{}{
		"Name": "lqm",
		"Age":  18,
	}

	err = mergo.MapWithOverwrite(&stur1, data)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(stur1)
}

func cert() {

}
