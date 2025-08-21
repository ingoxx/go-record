package main

import (
	"context"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
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

	secret, err := clientSet.CoreV1().Secrets("ingress-nginx").Get(context.TODO(), "ingress-nginx-admission", v12.GetOptions{})
	if err != nil {
		log.Fatalln("secret err >>>", err)
	}

	dir := "C:\\Users\\Administrator\\Desktop"

	for k, v := range secret.Data {
		os.WriteFile(filepath.Join(dir, k), v, 0777)
	}

}
