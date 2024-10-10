package account

import (
	"errors"
	"fmt"
	"strconv"
)

type account struct {
	Name      string
	accountId int
	password  string
	balance   float64
}

//工厂模式-构造函数
func NewAccount(acc int, name, pwd string, bal float64) (*account, error) {
	accString := strconv.FormatInt(int64(acc), 10)
	if len(accString) < 6 || len(accString) > 10 {
		return nil, errors.New("账号长度不对")
	}
	if len(pwd) != 6 {
		return nil, errors.New("密码长度不对")
	}
	if bal < 20.0 {
		return nil, errors.New("余额需要大于20.0")
	}
	return &account{
		Name:      name,
		accountId: acc,
		password:  pwd,
		balance:   bal,
	}, nil
}

//方法
func (m *account) SetAccountId(acc01 int) {
	accString := strconv.FormatInt(int64(acc01), 10)
	if len(accString) >= 6 && len(accString) <= 10 {
		m.accountId = acc01
	} else {
		fmt.Println("重置账号失败，账号长度不对...")
	}
}

func (m *account) GetAccountId() int {
	return m.accountId
}
