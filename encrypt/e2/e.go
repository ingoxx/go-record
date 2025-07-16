package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin" + "112233" + "gin-web"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil || len(hashedPassword) == 0 {
		return
	}

	fmt.Println(string(hashedPassword))

	err = decryptionPwd([]byte("$2a$10$qkDYMDcceCVOaUUmCjjMBu5MTqO5v6/O2MgF4cJUmS.pSQNN9TZ92"))
	fmt.Println("err >>> ", err)
}

func decryptionPwd(hashPwd []byte) (err error) {
	password := "二毛" + "123321" + "gin-web"
	err = bcrypt.CompareHashAndPassword(hashPwd, []byte(password))
	if err != nil {
		return
	}

	return
}
