package main

import (
	"flag"
	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	v12 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	_, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	_ = &certmanagerv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "ingress-nginx-kubebuilder-system",
			Name:      "certtest2",
		},
		Spec: certmanagerv1.CertificateSpec{
			DNSNames: []string{
				"bbb.cn",
				"aaa.cn",
			},
			IssuerRef: v12.ObjectReference{
				Name: "ingress-nginx-kubebuilder-selfsigned-issuer",
			},
			SecretName: "all-host-secret-test",
		},
	}
}
