package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"
)

func main() {
	tmpl := "D:\\project\\github.com\\ingoxx\\Golang-practise\\tmpl\\t3\\server.tmpl"
	proxy := "D:\\project\\github.com\\ingoxx\\Golang-practise\\tmpl\\t3\\proxy.tmpl"
	file, err2 := os.ReadFile(tmpl)
	if err2 != nil {
		log.Fatalln(err2)
	}
	mainTmpl, err := template.New("servers").Parse(string(file))

	proxyStruct := struct {
		HostName        string
		Path            string
		RenderRedirect  bool
		RenderSsl       bool
		RenderAllowList bool
		RenderDenyList  bool
	}{
		HostName:       "aaa.com",
		Path:           "/",
		RenderRedirect: true,
	}

	proxyTmplStr, err := os.ReadFile(proxy)
	if err != nil {
		log.Fatalln(err)
	}
	proxyTmpl, err := template.New("proxyTmpl").Parse(string(proxyTmplStr))
	if err != nil {
		log.Fatalln(err)
	}

	var dynamicTpl01 bytes.Buffer
	err = proxyTmpl.Execute(&dynamicTpl01, proxyStruct)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(dynamicTpl01.String())

	var dynamicTpl02 bytes.Buffer
	err = mainTmpl.Execute(&dynamicTpl02, proxyStruct)
	if err != nil {
		log.Fatalln("e1 >>> ", err)
	}

	var dynamicTpl03 bytes.Buffer
	_, err = mainTmpl.New("servers").Parse(dynamicTpl02.String())
	err = mainTmpl.Execute(&dynamicTpl03, dynamicTpl02)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(dynamicTpl03.String())

}
