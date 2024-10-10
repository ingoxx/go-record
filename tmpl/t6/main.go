package main

import (
	"os"
	"text/template"
)

// CIDRRange 结构体表示一个IP范围，包含起始和结束IP
type CIDRRange struct {
	StartIP string
	EndIP   string
}

// Location 结构体匹配模板中的$location，其中CIDR字段现在是一个CIDRRange结构体的切片
type Location struct {
	Allowlist struct {
		CIDR []*CIDRRange
	}
}

func main() {
	file := "C:\\Users\\Administrator\\Desktop\\out.pem"
	str := "111232sadad"
	os.WriteFile(file, []byte(str), 0777)
	// 初始化数据
	locationData := Location{
		Allowlist: struct{ CIDR []*CIDRRange }{
			CIDR: []*CIDRRange{
				{"192.168.1.0", "192.168.1.255"},
				{"10.0.0.0", "10.0.255.255"},
			},
		},
	}

	// 定义模板字符串
	tmplStr := `
{{ if gt (len .Allowlist.CIDR) 0 }}
{{ range $range := .Allowlist.CIDR }}
{{ if e

allow from {{ $range.StartIP }}
{{ else if ne $range.EndIP "" }}
to {{ $range.EndIP }};
{{ end }}
{{ end }}
deny all;
{{ end }}
`

	// 创建并解析模板
	tmpl, err := template.New("aclTemplate").Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	// 渲染模板到stdout
	err = tmpl.Execute(os.Stdout, locationData)
	if err != nil {
		panic(err)
	}
}
