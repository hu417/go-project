package uid

import (
	"github.com/bwmarrin/snowflake"
)

// 适用于：分布式系统中需要有序的唯一 ID，如订单号、用户 ID 等
func SnowflakeUuid() int64 {
	node, _ := snowflake.NewNode(1)
	// 类似 "1234567890123456789"
	return node.Generate().Int64()

}
