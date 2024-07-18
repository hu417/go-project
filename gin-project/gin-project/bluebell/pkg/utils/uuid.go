package utils

import "github.com/google/uuid"

// 随机id
func GetUuidInt() uint32 {  // 数字
	return uuid.New().ID()
}

func GetUuidString() string {  // 字符串
	return uuid.New().String()
}
