package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"
)

func main() {
	// 读取并解析主模板
	mainTmplStr, err := os.ReadFile("D:\\project\\github.com\\ingoxx\\Golang-practise\\tmpl\\t3\\server.tmpl")
	if err != nil {
		log.Fatal("Error reading nginx.conf.tmpl:", err)
	}
	mainTmpl, err := template.New("main").Parse(string(mainTmplStr))
	if err != nil {
		log.Fatal("Error parsing nginx.conf.tmpl:", err)
	}

	redirectTmplStr, err := os.ReadFile("D:\\project\\github.com\\ingoxx\\Golang-practise\\tmpl\\t3\\redirect.tmpl")
	if err != nil {
		log.Fatal("Error reading redirect.tmpl:", err)
	}
	redirectTmpl, err := template.New("subTmplStr").Parse(string(redirectTmplStr))

	proxyTmplStr, err := os.ReadFile("D:\\project\\github.com\\ingoxx\\Golang-practise\\tmpl\\t3\\proxy.tmpl")
	if err != nil {
		log.Fatal("Error reading redirect.tmpl:", err)
	}
	proxyTmpl, err := template.New("subTmplStr").Parse(string(proxyTmplStr))
	if err != nil {
		log.Fatal("Error reading redirect.tmpl:", err)
	}

	redirect := struct {
		HostName        string
		Path            string
		RenderRedirect  bool
		RenderSsl       bool
		RenderAllowList bool
		RenderDenyList  bool
	}{
		HostName:       "aaa.com",
		Path:           "/aaa",
		RenderRedirect: true,
	}

	proxy := struct {
		HostName        string
		Path            string
		RenderRedirect  bool
		RenderSsl       bool
		RenderAllowList bool
		RenderDenyList  bool
	}{
		HostName:       "bbb.com",
		Path:           "/bbb",
		RenderRedirect: true,
	}

	var tpl0 bytes.Buffer
	err = redirectTmpl.Execute(&tpl0, redirect)
	if err != nil {
		log.Fatal("Error executing redirectTmpl:", err)
	}

	var tpl2 bytes.Buffer
	err = proxyTmpl.Execute(&tpl2, proxy)
	if err != nil {
		log.Fatal("Error executing redirectTmpl:", err)
	}

	_, err = mainTmpl.New("redirectTmpl").Parse(tpl0.String())
	if err != nil {
		log.Fatal("Error executing redirectTmpl:", err)
	}

	_, err = mainTmpl.New("backend").Parse(tpl2.String())
	if err != nil {
		log.Fatal("Error executing redirectTmpl:", err)
	}

	// 执行渲染
	var tpl bytes.Buffer
	err = mainTmpl.Execute(&tpl, redirect)
	if err != nil {
		log.Fatal("Error executing template:", err)
	}

	// 输出最终的配置文件内容
	fmt.Println(tpl.String())

}
