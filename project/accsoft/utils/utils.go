package utils

import "fmt"

// type functype func() []string

type AccountDetail struct {
	InOut       string
	Blance      float64
	InOutBlance float64
	Directions  string
}

type AccountKeyWord struct {
	Confirm string
	Out     bool
}

type accountDetailSlice []AccountDetail

func NewAccountDetailSlice() *accountDetailSlice {
	return &accountDetailSlice{
		AccountDetail{
			InOut:       "",
			Blance:      0,
			InOutBlance: 0,
			Directions:  "",
		},
	}
}

func (a accountDetailSlice) Query(d accountDetailSlice, d2 *AccountDetail) {
	fmt.Println("-----------------当前收支明细-----------------")
	fmt.Println()
	if d2.InOut == "" {
		fmt.Println("\t去赚点钱吧...")
	} else {
		fmt.Println("收支\t账户金额\t收支金额\t说明")
		for _, v := range d {
			fmt.Printf("%v\t%f\t%f\t%v\n", v.InOut, v.Blance, v.InOutBlance, v.Directions)
		}
	}
}

func (a accountDetailSlice) Add(d *AccountDetail) accountDetailSlice {
	fmt.Println("-----------------当前登记收入-----------------")
	fmt.Println()
	fmt.Printf("\t收入金额:")
	fmt.Scanln(&d.InOutBlance)
	fmt.Printf("\t收入说明:")
	fmt.Scanln(&d.Directions)
	d.Blance += d.InOutBlance
	d.InOut = "收入"
	a = append(a, *d)
	return a
}

func (a accountDetailSlice) Reduce(d *AccountDetail) accountDetailSlice {
	fmt.Println("-----------------当前登记支出-----------------")
	fmt.Println()
	fmt.Printf("\t支出金额:")
	fmt.Scanln(&d.InOutBlance)
	fmt.Printf("\t支出说明:")
	fmt.Scanln(&d.Directions)
	fmt.Println()
	if d.Blance < d.InOutBlance || d.Blance == 0.0 {
		fmt.Println("\t超支啦...")
		return a
	} else {
		d.Blance -= d.InOutBlance
		d.InOut = "支出"
		a = append(a, *d)
		return a
	}

}

func (a accountDetailSlice) Exist(d *AccountKeyWord) {
	for {
		fmt.Printf("确定退出吗(输入y/n)?")
		fmt.Scanln(&d.Confirm)
		if d.Confirm == "y" {
			d.Out = true
			fmt.Println("\tbye...")
			break
		} else if d.Confirm == "n" {
			break
		} else {
			fmt.Println("输入错误...")
		}
	}
}

func (a accountDetailSlice) Run() {
	var num int
	var detailSlice accountDetailSlice
	income := AccountDetail{}
	keyWord := AccountKeyWord{}
	fc := FuncList()
	for !keyWord.Out {
		fmt.Println("--------家庭收支记账软件--------")
		fmt.Println()
		for i, v := range fc {
			if i == 0 {
				continue
			}
			fmt.Printf("\t%d %v\n", i, v)
		}
		fmt.Println()
		fmt.Printf("\t请选择(%d-%d):", 1, len(fc)-1)
		fmt.Scanln(&num)
		fmt.Println()
		if num <= 0 || num > len(fc)-1 {
			fmt.Println("\t输入错误")
			fmt.Println()
			continue
		}
		fmt.Println()
		switch fc[num] {
		case "退出":
			detailSlice.Exist(&keyWord)
		case "收支明细":
			detailSlice.Query(detailSlice, &income)
		case "登记收入":
			detailSlice = detailSlice.Add(&income)
		case "登记支出":
			detailSlice = detailSlice.Reduce(&income)
		}
		fmt.Println()
		fmt.Println("-------------------------------------------")
		fmt.Println()
	}
}

func FuncList() []string {
	fl := []string{1: "收支明细", 2: "登记收入", 3: "登记支出", 4: "退出"}
	return fl
}
