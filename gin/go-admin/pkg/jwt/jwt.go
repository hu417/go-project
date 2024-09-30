package jwt

import (
	"errors"
	"go-admin/global"

	"github.com/dgrijalva/jwt-go"
)


type UserClaim struct {
	Id      uint
	Name    string
	IsAdmin bool // 是否是超管
	RoleId  uint // 角色唯一标识
	jwt.StandardClaims
}

// GenerateToken 生成token
func GenerateToken(id uint, roleId uint, name string, expireAt int64) (string, error) {
	uc := &UserClaim{
		Id:     id,
		Name:   name,
		RoleId: roleId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString([]byte(global.JwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyzeToken Token 解析
func AnalyzeToken(token string) (*UserClaim, error) {
	uc := new(UserClaim)
	claims, err := jwt.ParseWithClaims(token, uc, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return uc, errors.New("token is invalid")
	}
	return uc, err
}
