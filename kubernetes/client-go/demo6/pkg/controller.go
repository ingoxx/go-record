package controller

import (
	"context"
	"k8s.io/api/autoscaling/v2beta1"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	v1ServiceInformer "k8s.io/client-go/informers/core/v1"
	v1IngressInformer "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	v1ServiceLister "k8s.io/client-go/listers/core/v1"
	v1IngressLister "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"reflect"
	"time"
)

type Controller struct {
	client        kubernetes.Interface
	ingressLister v1IngressLister.IngressLister
	serviceLister v1ServiceLister.ServiceLister
	queue         workqueue.RateLimitingInterface
}

func (c *Controller) addServiceHandler(obj interface{}) {
	c.inQueue(obj)
}

func (c *Controller) updateServiceHandler(oldObj interface{}, newObj interface{}) {
	equal := reflect.DeepEqual(oldObj, newObj)
	if !equal {
		c.inQueue(newObj)
	}
}

func (c *Controller) deleteServiceHandler(obj interface{}) {
	c.inQueue(obj)
}

func (c *Controller) inQueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}

	c.queue.Add(key)
}

func (c *Controller) deleteIngressHandler(obj interface{}) {
	ingress, ok := obj.(*v1.Ingress)
	if !ok {
		return
	}

	ownerReference := metav1.GetControllerOf(ingress)
	if ownerReference == nil {
		return
	}

	if ownerReference.Kind != "Service" {
		return
	}

	c.queue.Add(ingress.Namespace + "/" + ingress.Name)

}

func (c *Controller) Run(workers int, stopCh chan struct{}) {
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
}

func (c *Controller) runWorker() {
	for c.processNextItem() {

	}
}

func (c *Controller) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	defer c.queue.Done(key)

	err := c.serviceIngressCURD(key.(string))
	if err != nil {
		c.handleErr(err, key)
	}

	return true
}

func (c *Controller) serviceNotExistsDeleteIngress(igsErr error, ingress *v1.Ingress) error {
	if !errors.IsNotFound(igsErr) {
		err := c.client.NetworkingV1().Ingresses(ingress.Namespace).Delete(context.Background(), ingress.Name, metav1.DeleteOptions{})
		return err
	}
	return nil
}

func (c *Controller) serviceIngressCURD(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	serviceName := namespace + "/" + name
	if serviceName != "web/service-ingress" {
		return nil
	}

	ingress, igsErr := c.ingressLister.Ingresses(namespace).Get(name)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	service, err := c.serviceLister.Services(namespace).Get(name)
	if err != nil {
		err = c.serviceNotExistsDeleteIngress(igsErr, ingress)
		if err != nil {
			return err
		}
		return err
	}

	_, ok := service.GetAnnotations()["ingress/http"]
	if ok && errors.IsNotFound(igsErr) {
		// service的annotation存在，但是ingress不存在，就创建ingress
		addIngress := c.constructIngress(service)
		_, err = c.client.NetworkingV1().Ingresses(namespace).Create(context.Background(), addIngress, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else if !ok && ingress != nil {
		// 删除ingress
		err = c.client.NetworkingV1().Ingresses(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	}

	return nil

}

func (c *Controller) constructIngress(service *v12.Service) *v1.Ingress {
	var IngressClassName = "nginx"
	var PathType = v1.PathTypePrefix
	var addIngress = &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: service.Namespace,
			Name:      service.Name,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(service, v2beta1.SchemeGroupVersion.WithKind("Service")),
			},
		},
		Spec: v1.IngressSpec{
			IngressClassName: &IngressClassName,
			Rules: []v1.IngressRule{
				{
					Host: "lxb.pkg.com",
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &PathType,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: service.Name,
											Port: v1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return addIngress
}

func (c *Controller) handleErr(err error, key interface{}) {
	// 放回队列重试5次
	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("Error syncing key %v: %v, retrying...", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	// 重试了5次都不行就记录错误,并退出重试
	c.queue.Forget(key)
	runtime.HandleError(err)
	klog.Infof("%v handle failed.\n", key)
}

func NewController(client kubernetes.Interface, serviceInformer v1ServiceInformer.ServiceInformer, ingressInformer v1IngressInformer.IngressInformer) *Controller {
	c := &Controller{
		client:        client,
		ingressLister: ingressInformer.Lister(),
		serviceLister: serviceInformer.Lister(),
		queue: workqueue.NewRateLimitingQueueWithConfig(workqueue.DefaultControllerRateLimiter(),
			workqueue.RateLimitingQueueConfig{
				Name: "ingressMgr",
			}),
	}

	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addServiceHandler,
		UpdateFunc: c.updateServiceHandler,
		DeleteFunc: c.deleteServiceHandler},
	)

	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngressHandler,
	})

	return c
}
