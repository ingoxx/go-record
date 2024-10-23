package main

import (
	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	v12 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	cert := &certmanagerv1.Certificate{
		ObjectMeta: v1.ObjectMeta{
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

	sc := certmanagerv1.SchemeGroupVersion

}
