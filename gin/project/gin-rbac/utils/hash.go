package utils

import "golang.org/x/crypto/bcrypt"

const (
	passwordCost = 12 // 密码加密次数 密码哈希的难度级别
)

// HashPassword 生成密码的哈希值
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePasswordHash 验证密码
func ComparePasswordHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
