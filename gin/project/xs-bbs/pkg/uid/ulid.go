package uid

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// 适用于：需要排序且唯一的 ID，如消息队列、日志 ID 等
func Ulid() string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	// 类似 "01FHZB1E8B9B8YZW6CVDBKMT6T"
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

}
