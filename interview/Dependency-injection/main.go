package main

import (
	"fmt"
)

//Go语言（Golang）在设计模型中也支持依赖注入（Dependency Injection，DI），这是一种将依赖关系从高层模块传递到低层模块的设计模式，有助于降低耦合度，提高代码的可测试性和可维护性。
//在Go语言中实现依赖注入通常通过接口和构造函数的方式进行：
//例如，假设我们有一个UserService需要依赖UserRepository来操作用户数据：

// 定义User结构体

type User struct {
	ID   int
	Name string
}

// 定义UserRepository接口

type UserRepository interface {
	Find(id int) (*User, error)
	Save(user *User) error
}

// 实现UserRepository接口

type userRepository struct{}

func (u userRepository) Find(id int) (*User, error) {
	if id == 1 {
		return &User{ID: 1, Name: "Alice"}, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (u userRepository) Save(user *User) error {
	fmt.Printf("Saving user: %+v\n", user)
	return nil
}

// UserService依赖于UserRepository

type UserService struct {
	repo UserRepository
}

// 通过构造函数注入UserRepository

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int) (*User, error) {
	return s.repo.Find(id)
}

func (s *UserService) SaveUser(user *User) error {
	return s.repo.Save(user)
}

func main() {
	// 创建一个UserRepository实例
	userRepo := userRepository{}

	// 通过构造函数注入到UserService
	userService := NewUserService(userRepo)

	// 使用UserService获取用户
	user, err := userService.GetUser(1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Found user: %+v\n", user)
	}

	// 使用UserService保存用户
	newUser := &User{ID: 2, Name: "Bob"}
	err = userService.SaveUser(newUser)
	if err != nil {
		fmt.Println(err)
	}
}
