package customerview

import (
	"errors"
	"fmt"

	customer02 "github.com/ingoxx/Golang-practise/project/cusmanage/customerModel"
	customer03 "github.com/ingoxx/Golang-practise/project/cusmanage/customerService"
)

type CustomerList struct {
	CustomerService customer03.CustomerSystemList
	Key             customer02.KeyWord
	CustomerField   customer02.CustomerSystemField
}

func (c *CustomerList) List() {
	c.CustomerService.CustomerList()
}

func (c *CustomerList) Add() {
	c.CustomerService.AddCustomer(&c.CustomerField)
}

func (c *CustomerList) Edit() {
	c.CustomerService.EditCustomer(&c.CustomerField)

}

func (c *CustomerList) Del() {
	c.CustomerService.DelCustomer(&c.Key)
}

func (c *CustomerList) Exists() {
	c.CustomerService.Exists(&c.Key)
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

type customerSelect = func() []string //等同于 type customerSelect func() []string

func FindId(i int, c customerSelect) error {
	for i2 := range c() {
		if i2 == i {
			return nil
		}
	}
	return errors.New("输入错误")
}

func Run() {
	customer := CustomerList{}
	cl := SelectList()
	for !customer.Key.Out {
		fmt.Println()
		fmt.Println("--------客户管理系统--------")
		fmt.Println()
		customer.Key.Num = 0
		for i, v := range cl {
			if i == 0 {
				continue
			}
			fmt.Printf("\t%d %v\n", i, v)
		}
		fmt.Println()
		fmt.Printf("\t请选择(%d-%d):", 1, len(cl)-1)
		fmt.Scanln(&customer.Key.Num)
		fmt.Println()
		res := FindId(customer.Key.Num, SelectList)
		if res != nil {
			fmt.Println(res.Error())
			continue
		}
		fmt.Println()
		switch cl[customer.Key.Num] {
		case "退出":
			customer.Exists()
		case "客户列表":
			customer.List()
		case "客户添加":
			customer.Add()
		case "客户修改":
			customer.Edit()
		case "客户删除":
			customer.Del()
		}
		fmt.Println()
		fmt.Println("-------------------------------------------")
		fmt.Println()
	}
}
