package main

import (
	"flag"
	"fmt"
	"github.com/ingoxx/go-record/kubernetes/client-go/kubectl-plugs/p3/errors"
	"github.com/ingoxx/go-record/kubernetes/client-go/kubectl-plugs/p3/pkg/cmd"
)

func main() {
	flag.Parse()

	namespace := flag.Arg(0)
	if namespace == "" {
		panic(errors.MissDataError)
	}

	if err := cmd.NewKillNameSpace(namespace).KillNS(); err != nil {
		panic(fmt.Errorf("fail to kill namespace '%s', errMsg: %s", namespace, err.Error()))
	}

	fmt.Printf("kill namespace %s successfully.\n", namespace)
}
