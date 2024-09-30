package uid

import "github.com/google/uuid"

// 适合分布式系统中需要生成全局唯一 ID 的场景
func Uuid() string {
	// 类似 "550e8400-e29b-41d4-a716-446655440000"
	return uuid.New().String()
}
