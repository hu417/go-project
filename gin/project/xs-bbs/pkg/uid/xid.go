package uid

import (
	"github.com/rs/xid"
)


// 适用于：分布式系统中需要紧凑唯一 ID，如数据库主键、消息队列等
func Xid() string {
	// 类似 "9m4e2mr0ui3e8a215n4g"
	return xid.New().String()

}
