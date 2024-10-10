package utils

import (
	"errors"
	"fmt"
)

type CustomerSystemField struct {
	Id      int
	Name    string
	Tel     int
	Address string
}

type KeyWord struct {
	Confirm string
	Out     bool
	Num     int
}

type CustomerSystemList []CustomerSystemField

func (c CustomerSystemList) AddCustomer(a *CustomerSystemField) {
	fmt.Println("---------客户添加----------")
	fmt.Println()
	fmt.Printf("客户姓名:")
	fmt.Scanln(&a.Name)
	fmt.Printf("客户电话:")
	fmt.Scanln(&a.Tel)
	fmt.Printf("客户地址:")
	fmt.Scanln(&a.Address)
	a.Id += 1
	length := len(c)
	c = append(c, *a)
	fmt.Println()
	if length < len(c) {
		fmt.Printf("%v客户添加完成", a.Name)
	} else {
		fmt.Printf("%v客户添加失败", a.Name)
	}
}

func (c CustomerSystemList) EditCustomer(e *CustomerSystemField) CustomerSystemList {
	fmt.Println("---------客户修改----------")
	if len(c) == 0 {
		fmt.Println("请先添加客户")
		return nil
	}
	var id int
	var name string
	var tel int
	var address string
	fmt.Printf("请选择需要修改的客户编号:")
	fmt.Scanln(&id)
	id -= 1
	check := c.Check(id)
	if check != nil {
		fmt.Println(check)
		return c
	}
	fmt.Printf("客户姓名(%v):", c[id].Name)
	fmt.Scanln(&name)
	if name != "" {
		c[id].Name = name
		id += 1
		fmt.Printf("编号:%d, 客户姓名已修改成:%v", id, name)
	}
	fmt.Printf("客户电话(%v):", c[id].Tel)
	fmt.Scanln(&tel)
	if tel != 0 {
		c[id].Tel = tel
		id += 1
		fmt.Printf("编号:%d, 客户电话已修改成:%d", id, tel)
	}
	fmt.Printf("客户地址(%v):", c[id].Address)
	fmt.Scanln(&address)
	if address != "" {
		c[id].Address = address
		id += 1
		fmt.Printf("编号:%d, 客户地址已修改成:%v", id, address)
	}
	return c
}

func (c CustomerSystemList) DelCustomer(k *KeyWord) CustomerSystemList {
	fmt.Println("---------客户删除----------")
	fmt.Println()
	if len(c) == 0 {
		fmt.Println("先添加客户吧")
		return nil
	}
	var delId int
	fmt.Printf("请输入要删除的客户编号:")
	fmt.Scanln(&delId)
	delId -= 1
	check := c.Check(delId)
	if check != nil {
		fmt.Println(check)
		return c
	}
	length := len(c)
	delIndex := delId
	for {
		fmt.Printf("确定删除吗(y/n)?")
		fmt.Scanln(k.Confirm)
		if k.Confirm == "y" {
			c = append(c[0:delIndex], c[delIndex+1:length]...)
			delIndex += 1
			if len(c) < length {
				fmt.Printf("编号%d删除完成", delIndex)
			} else {
				fmt.Printf("编号%d删除失败", delIndex)
			}
		} else if k.Confirm == "n" {
			return c
		} else {
			fmt.Println("输入错误")
		}
	}
}

func (c CustomerSystemList) CustomerList() {
	fmt.Println("---------客户列表----------")
	fmt.Println()
	if len(c) == 0 {
		fmt.Println("还没有客户哦，先添加吧.")
	} else {
		fmt.Println("客户编号\t客户姓名\t客户电话\t客户地址")
		for _, v := range c {
			fmt.Printf("%d\t\t%v\t\t%d\t%v\n", v.Id, v.Name, v.Tel, v.Address)
		}
	}
}

func (c CustomerSystemList) Exists(k *KeyWord) {
	for {
		fmt.Printf("确定退出吗(输入y/n)?")
		fmt.Scanln(&k.Confirm)
		if k.Confirm == "y" {
			k.Out = true
			fmt.Println("\tbye...")
			break
		} else if k.Confirm == "n" {
			break
		} else {
			fmt.Println("输入错误...")
		}
	}
}

func (c CustomerSystemList) Check(i int) (err error) {
	if i < len(c) && i >= 0 {
		return nil
	}
	return errors.New("客户编号不存在")
}

func (c CustomerSystemList) Run() {
	var cf CustomerSystemList
	keyWord := KeyWord{}
	customer := CustomerSystemField{}
	cl := SelectList()
	for !keyWord.Out {
		fmt.Println()
		fmt.Println("--------客户管理系统--------")
		fmt.Println()
		keyWord.Num = 0
		for i, v := range cl {
			if i == 0 {
				continue
			}
			fmt.Printf("\t%d %v\n", i, v)
		}
		fmt.Println()
		fmt.Printf("\t请选择(%d-%d):", 1, len(cl)-1)
		fmt.Scanln(&keyWord.Num)
		fmt.Println()
		if keyWord.Num <= 0 || keyWord.Num > len(cl)-1 {
			fmt.Printf("\t输入错误,请选择(%d-%d)", 1, len(cl)-1)
			fmt.Println()
			continue
		}
		fmt.Println()
		switch cl[keyWord.Num] {
		case "退出":
			cf.Exists(&keyWord)
		case "客户列表":
			cf.CustomerList()
		case "客户添加":
			cf.AddCustomer(&customer)
		case "客户修改":
			cf = cf.EditCustomer(&customer)
		case "客户删除":
			cf = cf.DelCustomer(&keyWord)
		}
		fmt.Println()
		fmt.Println("-------------------------------------------")
		fmt.Println()
	}
}

func SelectList() []string {
	c := []string{
		1: "客户列表",
		2: "客户添加",
		3: "客户修改",
		4: "客户删除",
		5: "退出"}
	return c
}
