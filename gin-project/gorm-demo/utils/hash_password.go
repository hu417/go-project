package utils

import "golang.org/x/crypto/bcrypt"

// 密码加密
func HashPw(password string) string {
	const cost = 10
	HashPw, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		panic(err)
	}
	return string(HashPw)
}
