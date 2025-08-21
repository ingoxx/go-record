package main

import (
	"github.com/ingoxx/k8s-client-go/demo6/pkg"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {
	var stopCh = make(chan struct{})

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

	factory := informers.NewSharedInformerFactoryWithOptions(clientSet, 0, informers.WithNamespace("web"))
	services := factory.Core().V1().Services()
	ingresses := factory.Networking().V1().Ingresses()

	controllers := controller.NewController(clientSet, services, ingresses)
	// 启动informer
	factory.Start(stopCh)
	// 同步所有数据到delta fifo queue
	factory.WaitForCacheSync(stopCh)
	// 开始监控指定资源
	controllers.Run(5, stopCh)
}

// config=>创建listwatch=>注册事件并启动informer=>同步所有数据到delta fifo queue=>开始监听指定资源变化
