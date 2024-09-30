package uid

import (

    "github.com/teris-io/shortid"
)

// 适合于：需要短链或验证码的场景
func Shortid() string{
	id, _ := shortid.Generate()
	return id
}