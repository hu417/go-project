package common

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"rbac-v1/common/constants"
	"rbac-v1/model/po"
	"time"
)

var jwtKey = []byte(constants.JWT_SALT)

type CustomClaims struct {
	UserId uint `json:"user_id"`
	Name string `json:"name"`
	Username string `json:"username"`
	jwt.StandardClaims
}

//生成token
func GenerateToken(user *po.User) (tokenStr string, exp int64, err error) {
	expirationTime := time.Now().Add(constants.JWT_TOKEN_EXP)
	claims := &CustomClaims{
		UserId:         user.Id,
		Name:           user.Name,
		Username:       user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.UnixMilli(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expirationTime.UnixMilli(), nil
}

//校验token
func ParseToken(tokenStr string) (claims *CustomClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(*jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	err = errors.New("token is invalid")

	return nil, err
}