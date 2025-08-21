package main

import (
	"encoding/json"
	"fmt"
)

// 有用json技巧

type People struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Men struct {
	People
	Age string `json:"age"` // 这里注意，在Men里设置了同名但不同类型的字段Age，会替换掉People字段Age的int类型
}

type Woman struct {
	Men
	Skill  string `json:"skill"`
	Hobby  string `json:"-"`                //这里的-表示在序列化的时候忽略这个字段
	Weight int    `json:"weight,omitempty"` //这里表示这个字段是空值也就是没有赋值就在序列化的时候忽略它
	IsGood bool   `json:"is_good,omitempty"`
}

type Respone struct {
	Code int    `json:"code,string"` //这里表示序列化的时候会自动把int类型转换成string类型
	Resp string `json:"resp"`
}

type UnKnowFieldType struct {
	Code int             `json:"code"`
	Resp json.RawMessage `json:"resp"` //如果确实不知道结构体字段的类型，可以先声明为json.RawMessage类型，这个类型可以存储任意格式的json文本
}

func main() {
	// ----------------好用1------------------
	fmt.Println("-----men-------")
	m := new(Men)

	m.Age = "31"
	m.Name = "lxb"

	r1, _ := json.Marshal(m)
	fmt.Println(string(r1)) // output: {"name":"lxb","age":"31"}

	// ----------------好用2------------------
	fmt.Println("-----Woman-------")
	w := new(Woman)
	w.Name = "lqm"
	w.Age = "31"
	w.Skill = "Filial piety and kindness"
	//w.Weight = 0
	w.IsGood = false
	r2, _ := json.Marshal(w)
	fmt.Println(string(r2)) // output: {"name":"lqm","age":"31","skill":"Filial piety and kindness"}

	// ----------------好用3------------------
	fmt.Println("-----respone-------")
	resp := new(Respone)

	resp.Code = 200
	resp.Resp = "succeed"

	r3, _ := json.Marshal(resp)
	fmt.Println(string(r3)) // output: {"code":"200","resp":"succeed"}

	// ----------------好用4------------------
	jsonStr := `{"code": 200, "resp": {"name": "lxb", "age": 31}}`
	var uk UnKnowFieldType
	json.Unmarshal([]byte(jsonStr), &uk)
	fmt.Println(string(uk.Resp)) // output: {"name": "lxb", "age": 31}

}
