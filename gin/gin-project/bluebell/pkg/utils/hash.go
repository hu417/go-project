package utils

import "golang.org/x/crypto/bcrypt"

// 生成加密密码
func GeneratePassword(password string) (string, error) {
	const cost = 10 //哈希成本
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}
 
// 密码效验
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
