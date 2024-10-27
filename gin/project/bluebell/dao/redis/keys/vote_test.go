package keys

import (
	"bluebell/global"
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// 模拟 global.RDS
var mockRDS *redis.Client

// TestCreatePost 测试创建帖子
func TestCreatePost(t *testing.T) {
	// 模拟上下文
	ctx := context.Background()

	// 测试异常情况，模拟 pipeline.Exec 出错
	mockRDS = redis.NewClient(&redis.Options{
		// 模拟连接错误
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		DB:       0,
		Password: "qaz123",
	})

	global.RDS = mockRDS

	// 测试正常情况
	err := CreatePost(ctx, "post2", "community1")
	if err != nil {
		t.Errorf("CreatePost 正常情况出错: %v", err)
	}
}

// TestVoteForPost 测试投票
func TestVoteForPost(t *testing.T) {
	// 模拟上下文
	ctx := context.Background()

	// err
	var err error

	// 模拟连接
	mockRDS = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		DB:       0,
		Password: "qaz123",
	})

	global.RDS = mockRDS

	// 模拟用户 ID 和帖子 ID
	userID := "user1"
	postID := "post2"

	// 测试投票时间未过期
	err = VoteForPost(ctx, userID, postID, 1)
	assert.NoError(t, err)

	// // 测试重复投票
	// err = VoteForPost(ctx, userID, postID, 1)
	// assert.Equal(t, ErrVoteRepeated, err)

	// // 测试投票时间过期
	// // 将帖子发布时间设置为一周前
	// global.RDS.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), redis.Z{
	// 	Score:  float64(time.Now().Add(-7 * 24 * time.Hour).Unix()),
	// 	Member: postID,
	// })
	// err = VoteForPost(ctx, userID, postID, 1)
	// assert.Equal(t, ErrVoteTimeExpire, err)

	// // 测试取消投票
	// err = VoteForPost(ctx, userID, postID, 0.0)
	// assert.NoError(t, err)
}
