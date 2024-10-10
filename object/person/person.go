package person

import "fmt"

type person struct {
	Name   string
	age    int
	salary float64
}

func Newperson(n string) *person {
	return &person{
		Name: n,
	}
}

//setter
func (s *person) Setage(age int) {
	if age >= 18 && age < 100 {
		s.age = age
	} else {
		fmt.Println("未成年")
	}
}

func (s *person) Getage() int {
	return s.age
}

//getter
func (s *person) Setsalary(sal float64) {
	if sal > 5000 && sal < 20000 {
		s.salary = sal
	} else {
		fmt.Println("薪水不合理")
	}
}

func (s *person) Getsalary() float64 {
	return s.salary
}
