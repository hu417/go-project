package uid

import (
	"github.com/segmentio/ksuid"
)

// 适用于：需要唯一且具备时间排序特性的场景，如日志、消息队列等
func KsUid() string {
	// 输出类似 "1CHzCEl82ZZe2r2KS1qEz3XF8Ve"
	return ksuid.New().String()
}
