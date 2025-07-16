package cmd

import (
	"context"
	ce "github.com/ingoxx/go-record/kubernetes/client-go/kubectl-plugs/p3/errors"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	nsChan = make(chan func() *coreV1.Namespace)
)

type resources struct {
	name   string
	client *kubernetes.Clientset
}

func loadConfig() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		inClusterConfig, err := rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
		config = inClusterConfig
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientSet
}

type KillNameSpace struct {
	rs *resources
}

func NewKillNameSpace(ns string) *KillNameSpace {
	kn := &KillNameSpace{
		rs: &resources{
			name:   ns,
			client: loadConfig(),
		},
	}

	return kn
}

func (n *KillNameSpace) getNS() *coreV1.Namespace {
	ns := new(coreV1.Namespace)
	ns, err := n.rs.client.CoreV1().Namespaces().Get(context.Background(), n.rs.name, metav1.GetOptions{})
	if err != nil {
		panic(ce.NotFoundError)
	}

	if ns.DeletionTimestamp == nil || ns.DeletionTimestamp.IsZero() {
		panic(ce.NewDeleteError(n.rs.name))
	}

	return ns
}

func (n *KillNameSpace) clearNSFinalizers() error {
	ns := n.getNS()
	if len(ns.Spec.Finalizers) != 0 {
		ns.Spec.Finalizers = []coreV1.FinalizerName{}
		if _, err := n.rs.client.CoreV1().Namespaces().Finalize(context.Background(), ns, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}

	return nil
}

func (n *KillNameSpace) KillNS() error {
	if err := n.clearNSFinalizers(); err != nil {
		return err
	}

	if err := n.rs.client.CoreV1().Namespaces().Delete(context.Background(), n.rs.name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	return nil
}
