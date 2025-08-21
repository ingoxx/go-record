package main

import (
	"os"
	"text/template"
)

type Names struct {
	Names []string `json:"names"`
}

func main() {
	n := Names{Names: []string{"jay", "jay2"}}

	tmpl := template.Must(template.ParseFiles("D:\\project\\github.com\\ingoxx\\Golang-practise\\tmpl\\t2\\a.tmpl"))

	err := tmpl.Execute(os.Stdout, n)
	if err != nil {
		panic(err)
	}
}
