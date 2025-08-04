package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 假设这是用户注册时输入的密码或需要哈希的字符串
	originalStr := "user_1123sadd"

	// 1. 生成哈希值
	// 第二个参数是 cost，值越高，哈希计算越慢，但也越安全。
	// bcrypt.DefaultCost (值为10) 是一个很好的默认值。
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(originalStr), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("无法生成哈希: %v", err)
	}
	hashedStr := string(hashedBytes)

	fmt.Printf("原始字符串: %s\n", originalStr)
	fmt.Printf("生成的哈希值: %s\n", hashedStr)

	// 注意：每次运行，生成的哈希值都不同，因为 bcrypt 自动加入了随机的 "盐" (salt)。

	// 2. 验证哈希值 (模拟用户登录)
	// 假设这是用户登录时输入的字符串
	loginAttemptStr := "user_1123sadd"

	// 使用 CompareHashAndPassword 来比较原始字符串和哈希值
	err = bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(loginAttemptStr))
	if err != nil {
		// 如果 err 不为 nil，说明密码不匹配
		log.Printf("验证失败: 密码不匹配 (%v)", err)
	} else {
		fmt.Println("验证成功: 密码匹配!")
	}

	// 测试一个错误的密码
	wrongLoginAttemptStr := "wrong_password"
	err = bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(wrongLoginAttemptStr))
	if err != nil {
		fmt.Printf("对错误密码 '%s' 的验证失败，这是预期的行为。\n", wrongLoginAttemptStr)
	}
}
