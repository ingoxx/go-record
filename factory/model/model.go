package model

type student struct {
	Name string
	age  int
}

//场景：当定义个结构体变量名首字母是小写时，但是又要可以在别的包里边可以创建这个结构体的实例，这里就需要用到工厂模式
//下面就是工厂模式
func Newstudent(n string, a int) *student {
	return &student{
		Name: n,
		age:  a,
	}
}

func (m *student) Newage() int { //封装
	if m.age != 20 {
		return 0
	}
	return m.age //相当于私有字段，跟python的私有属性__val，双下划线表示私有
}
