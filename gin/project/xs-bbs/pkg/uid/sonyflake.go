package uid

import (
	"github.com/sony/sonyflake"
)


// 适用于：高性能分布式环境，如订单号、日志序列。
func SonyflakeUid() uint64 {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, _ := sf.NextID()
	// 输出类似 "173301874793540608"
	return id

}
