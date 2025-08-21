package main

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

const nginxTmpl = `
worker_processes  4;
error_log  /var/log/nginx/errors.log notice;
pid        /var/run/nginx.pid;
worker_rlimit_nofile 1047552;
worker_shutdown_timeout 240s ;

events {
	multi_accept        on;
	worker_connections  16384;
	use                 epoll;
}

http {
	{{ template "server" }}
}
`

const serverTmpl = `# start {{ .HostName }}
server {
    listen 80;
    listen [::]:80;
	listen 443;
	listen [::]:443;
	server_name {{ .HostName }};

    {{ if .RenderRedirect }}
        location {{ .Path }} {
            return 301 http://{{ .HostName }};
        }
    {{ end }}
    {{ if .RenderSsl }}{{ end }}
    {{ if .RenderAllowList }}{{ template "allowListTmpl" }}{{ end }}
    {{ if .RenderDenyList }}{{ template "denyListTmpl" }}{{ end }}

    # allow cos
    add_header 'Access-Control-Allow-Origin' '*';
    add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
    add_header 'Access-Control-Allow-Headers' 'DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,xfilecategory,xfilename,xfilesize';
    add_header 'Access-Control-Expose-Headers' 'Content-Length,Content-Range';
    if ($request_method = 'OPTIONS') {
        return 204;
    }

    # webSocket enable
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";

     location {{ .Path }} {
		proxy_pass http://{{ .HostName }};
		proxy_set_header Host $host;
		proxy_set_header X-Real-IP $remote_addr;
	}
}
` // 省略，直接粘贴server.tmpl的内容
const backendTmpl = ` location {{ .Path }} {
    proxy_pass http://{{ .HostName }};
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}` // 省略，直接粘贴backend.tmpl的内容

func main() {

	redirect := []struct {
		HostName        string
		Path            string
		RenderRedirect  bool
		RenderSsl       bool
		RenderAllowList bool
		RenderDenyList  bool
	}{
		{HostName: "ccc.com", RenderRedirect: false},
		{HostName: "aaa.com", Path: "/api", RenderRedirect: false},
		{HostName: "bbb.com", Path: "/mg", RenderRedirect: false},
	}

	fmt.Println(len(redirect))

	b := struct {
		Bytes bytes.Buffer
	}{}

	mainTmpl, err := template.New("main").Parse(nginxTmpl)
	if err != nil {
		log.Fatal("Error parsing nginx.conf.tmpl:", err)
	}

	//var tpl0 bytes.Buffer
	//b.Bytes = tpl0
	for _, v := range redirect[1:] {
		subTmpl, err := template.New("main").Parse(serverTmpl)
		if err != nil {
			log.Fatal("Error parsing nginx.conf.tmpl:", err)
		}

		err = subTmpl.Execute(&b.Bytes, v)
		if err != nil {
			log.Fatal("Error executing redirectTmpl:", err)
		}
	}

	_, err = mainTmpl.New("server").Parse(b.Bytes.String())
	if err != nil {
		log.Fatal("Error executing redirectTmpl:", err)
	}

	// 执行渲染
	var tpl bytes.Buffer
	err = mainTmpl.Execute(&tpl, nil)
	if err != nil {
		log.Fatal("Error executing template:", err)
	}

	// 输出最终的配置文件内容
	fmt.Println(tpl.String())
}
