package customerservice

import (
	"errors"
	"fmt"

	customer01 "github.com/ingoxx/Golang-practise/project/cusmanage/customerModel"
)

type CustomerSystemList []customer01.CustomerSystemField

// 这里的CustomerSystemList必须是指针类型，否则后面每次添加客户都会显示空因为走的是值拷贝
func (c *CustomerSystemList) AddCustomer(a *customer01.CustomerSystemField) {
	fmt.Println("---------客户添加----------")
	fmt.Println()
	fmt.Printf("客户姓名:")
	fmt.Scanln(&a.Name)
	fmt.Printf("客户电话:")
	fmt.Scanln(&a.Tel)
	fmt.Printf("客户地址:")
	fmt.Scanln(&a.Address)
	length := len(*c)
	*c = append(*c, *a)
	fmt.Println()
	if length < len(*c) {
		fmt.Printf("%v客户添加完成", a.Name)
	} else {
		fmt.Printf("%v客户添加失败", a.Name)
	}
}

func (c *CustomerSystemList) EditCustomer(e *customer01.CustomerSystemField) {
	fmt.Println("---------客户修改----------")
	if len(*c) == 0 {
		fmt.Println("请先添加客户")
		return
	}
	var id int
	var name string
	var tel int
	var address string
	fmt.Printf("请选择需要修改的客户编号:")
	fmt.Scanln(&id)
	id -= 1
	check := c.CheckId(id)
	if check != nil {
		fmt.Println(check.Error())
		return
	}
	fmt.Printf("客户姓名(%v):", (*c)[id].Name)
	fmt.Scanln(&name)
	if name != "" {
		(*c)[id].Name = name
		id += 1
		fmt.Printf("编号:%d, 客户姓名已修改成:%v\n", id, name)
		id -= 1
	}
	fmt.Printf("客户电话(%v):", (*c)[id].Tel)
	fmt.Scanln(&tel)
	if tel != 0 {
		(*c)[id].Tel = tel
		id += 1
		fmt.Printf("编号:%d, 客户电话已修改成:%d\n", id, tel)
		id -= 1
	}
	fmt.Printf("客户地址(%v):", (*c)[id].Address)
	fmt.Scanln(&address)
	if address != "" {
		(*c)[id].Address = address
		id += 1
		fmt.Printf("编号:%d, 客户地址已修改成:%v\n", id, address)
		id -= 1
	}
}

func (c *CustomerSystemList) DelCustomer(k *customer01.KeyWord) {
	fmt.Println("---------客户删除----------")
	fmt.Println()
	if len(*c) == 0 {
		fmt.Println("先添加客户吧")
		return
	}
	delId := -1
	fmt.Printf("请输入要删除的客户编号:")
	fmt.Scanln(&delId)
	if delId == -1 {
		return
	}
	delId -= 1
	check := c.CheckId(delId)
	if check != nil {
		fmt.Println(check.Error())
		return
	}
	length := len(*c)
	delIndex := delId
	for {
		fmt.Printf("确定删除吗(y/n)?")
		fmt.Scanln(&k.Confirm)
		if k.Confirm == "y" {
			*c = append((*c)[0:delIndex], (*c)[delIndex+1:length]...)
			delIndex += 1
			if len(*c) < length {
				fmt.Printf("编号%d删除完成", delIndex)
			} else {
				fmt.Printf("编号%d删除失败", delIndex)
			}
			break
		} else if k.Confirm == "n" {
			break
		} else {
			fmt.Println("输入错误")
		}
	}
}

func (c *CustomerSystemList) CustomerList() {
	fmt.Println("---------客户列表----------")
	fmt.Println()
	if len(*c) == 0 {
		fmt.Println("还没有客户哦，先添加吧.")
	} else {
		fmt.Println("客户编号\t客户姓名\t客户电话\t客户地址")
		for i, v := range *c {
			i += 1
			fmt.Printf("%d\t\t%v\t\t%d\t%v\n", i, v.Name, v.Tel, v.Address)
		}
	}
}

func (c *CustomerSystemList) Exists(k *customer01.KeyWord) {
	for {
		fmt.Printf("确定退出吗(输入y/n)?")
		fmt.Scanln(&k.Confirm)
		if k.Confirm == "y" {
			k.Out = true
			fmt.Println("bye...")
			break
		} else if k.Confirm == "n" {
			break
		} else {
			fmt.Println("输入错误...")
		}
	}
}

func (c *CustomerSystemList) CheckId(i int) error {
	for i2 := range *c {
		if i2 == i {
			return nil
		}
	}
	return errors.New("客户编号不存在")
}
