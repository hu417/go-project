package utils

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// 生成随机字符串
func NanoId() string {
	// return gonanoid.Must(32)

	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 10)
	if err != nil {
		panic(err)
	}
	return id

}
