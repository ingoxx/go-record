package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	// 加载主模板和子模板
	tmpl := template.Must(template.ParseFiles("C:\\Users\\Administrator\\Desktop\\main.tmpl", "C:\\Users\\Administrator\\Desktop\\render.tmpl"))

	data := struct {
		Name string
	}{
		Name: "张三",
	}

	// 执行渲染
	err := tmpl.ExecuteTemplate(os.Stdout, "render", data)
	if err != nil {
		log.Fatalf("执行模板渲染失败: %v", err)
	}

}
