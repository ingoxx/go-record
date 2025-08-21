package main

import "fmt"

//在Go语言中，状态模式用于处理对象在其生命周期内根据内部状态改变其行为。下面是一个使用状态模式的简单示例，模拟ATM机根据账户余额的状态进行取款操作

// 定义状态接口

type AccountState interface {
	Withdraw(account *Account, amount int) error
	Deposit(account *Account, amount int)
}

// 账户结构体，包含状态字段

type Account struct {
	balance int
	state   AccountState
}

func (a *Account) ChangeState(state AccountState) {
	a.state = state
}

func (a *Account) Withdraw(amount int) error {
	return a.state.Withdraw(a, amount)
}

func (a *Account) Deposit(amount int) {
	a.state.Deposit(a, amount)
}

// 状态实现：正常状态（BalanceOK）

type BalanceOK struct{}

func (bk BalanceOK) Withdraw(account *Account, amount int) error {
	if account.balance >= amount {
		account.balance -= amount
		fmt.Printf("Withdrawal successful. New balance: %d\n", account.balance)
		return nil
	}
	return fmt.Errorf("insufficient balance for withdrawal")
}

func (bk BalanceOK) Deposit(account *Account, amount int) {
	account.balance += amount
	fmt.Printf("Deposit successful. New balance: %d\n", account.balance)
}

// 状态实现：透支状态（Overdraft）

type Overdraft struct{}

func (od Overdraft) Withdraw(account *Account, amount int) error {
	if amount <= 500 { // 假设允许透支500元
		account.balance -= amount
		fmt.Printf("Withdrawal successful with overdraft. New balance: %d\n", account.balance)
		return nil
	}
	return fmt.Errorf("exceeded overdraft limit")
}

func (od Overdraft) Deposit(account *Account, amount int) {
	account.balance += amount
	fmt.Printf("Deposit successful. New balance: %d\n", account.balance)
	// 当存款后余额足够时，切换回正常状态
	if account.balance > -500 {
		account.ChangeState(BalanceOK{})
	}
}

func main() {
	account := &Account{balance: 1000, state: BalanceOK{}}

	account.Withdraw(800) // 正常取款
	account.Withdraw(800) // 导致透支
	account.Deposit(900)  // 存款后余额足够，状态切换回正常

	// 再次尝试取款
	account.Withdraw(800)
}

//在这个例子中，Account类包含一个状态字段，并提供了根据当前状态进行取款和存款的方法。
//我们定义了两种状态：BalanceOK和Overdraft，每种状态都有不同的取款逻辑。
//当账户状态从正常变为透支时，取款规则会发生变化；
//而当存款后账户余额恢复正常时，状态会自动切换回正常状态
