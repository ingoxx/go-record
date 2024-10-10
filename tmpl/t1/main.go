package main

import (
	"bytes"
	"fmt"
	"text/template"
)

// Info 结构体定义
type Info struct {
	Path string `json:"path"`
	Url  string `json:"url"`
}

func main() {
	i := Info{"/aaa", "http://aaa.com"}

	// 准备模板字符串，这里为了示例直接写死了内容，实际上你可以从文件读取
	tmplStr := `
location {{.Path}} {
    proxy_pass {{.Url}};
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
`

	// 创建一个新的模板并解析模板字符串
	tmpl, err := template.New("config").Parse(tmplStr)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// 创建一个缓冲区来保存渲染结果
	var tpl bytes.Buffer

	// 执行渲染，将i的实例传递给模板
	err = tmpl.Execute(&tpl, i)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	// 输出渲染后的结果
	fmt.Println(tpl.String())
}
