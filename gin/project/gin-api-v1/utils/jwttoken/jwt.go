package jwttoken

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 自定义Claim结构体
type MyCustomClaims struct {
	User_Id  string
	Username string
	Role     string
	jwt.RegisteredClaims
}

// 随机字符串
var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(str_len int) string {
	rand_bytes := make([]rune, str_len)
	for i := range rand_bytes {
		rand_bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(rand_bytes)
}

type jwtStr struct{}

func Newjwt() *jwtStr {
	return &jwtStr{}
}

// 生成token
func (*jwtStr) GenerateJwt(userid string, username, role string, sign_key string, expires int64) (string, error) {
	// 初始化claim
	claim := &MyCustomClaims{
		User_Id:  userid,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                                            // 签发者
			Subject:   "Jwt",                                                                    // 签发对象
			Audience:  jwt.ClaimStrings{"Android_APP", "IOS_APP"},                               // 签发受众
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                           // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)),                          // 生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expires))), // 过期时间
			ID:        randStr(10),                                                              // jwt ID, 类似于盐值
		},
	}

	// 生成token
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(sign_key))

}

/*
自定义的 ParseJwt 函数，它负责解析 JWT 字符串，并根据验证结果返回 Claims 数据和一个可能的存在的错误。ParseJwt 函数内部利用 jwt.Parse 解析 JWT 字符串。解析后，函数检查得到的 token 对象的 Valid 属性以确认 Claims 是否有效。有效性检查包括但不限于验证签名、检查 token 是否过期。如果 token 通过所有验证，函数返回 Claims 数据；如果验证失败（如签名不匹配或 token 已过期），则返回错误
*/
func (*jwtStr) ParseJwt(sign_key string, tokenstr string, options ...jwt.ParserOption) (*MyCustomClaims, error) {
	//解析、验证并返回token。
	str, err := jwt.ParseWithClaims(tokenstr, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(sign_key), nil
	})

	// 校验 Claims 对象是否有效，基于 exp（过期时间），nbf（不早于），iat（签发时间）等进行判断（如果有这些声明的话）。
	if err == nil && str.Valid {
		if claims, ok := str.Claims.(*MyCustomClaims); ok {
			fmt.Printf("%v %v\n", claims.Username, claims.RegisteredClaims)
			return claims, nil
		}
		return nil, errors.New("claims 解析失败")
	}
	return nil, err

}
