package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptMake 生成密码hash
func BcryptMake(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

// BcryptMakeCheck 校验密码hash
func BcryptMakeCheck(pwd, hashedPwd string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd))

	return err == nil
}
