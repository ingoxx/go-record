package main

import (
	"context"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 创建一个新的RESTClientConfig，它会自动从kubeconfig文件读取配置
	configFlags := genericclioptions.NewConfigFlags(true)
	configLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configFlags, &clientcmd.ConfigOverrides{})
	restConfig, err := configLoader.ClientConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client configuration: %v\n", err)
		os.Exit(1)
	}

	// 创建Kubernetes客户端
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	// 获取Pod列表
	pods, err := clientset.CoreV1().Pods("").List(context.Background(), v1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing pods: %v\n", err)
		os.Exit(1)
	}

	// 打印Pod列表
	for _, pod := range pods.Items {
		fmt.Printf("Name: %s, Namespace: %s\n", pod.Name, pod.Namespace)
	}
}
