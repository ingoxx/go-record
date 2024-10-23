package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatal(err)
	}

	dynamicSet, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	certRes := schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1", Resource: "issuers"}
	issuer := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "cert-manager.io/v1",
			"kind":       "Issuer",
			"metadata": map[string]interface{}{
				"name":      "test" + "-issuer",
				"namespace": "ingress-nginx-kubebuilder-system",
			},
			"spec": map[string]interface{}{
				"selfSigned": map[string]interface{}{},
			},
		},
	}
	create, err := dynamicSet.Resource(certRes).Namespace("ingress-nginx-kubebuilder-system").Create(context.Background(), issuer, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(create.GetName())
}
