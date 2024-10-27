package keys

import (
	"context"
	"testing"

	"bluebell/controller/request"
	"bluebell/global"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// TestGetPostIDsInOrder 测试 GetPostIDsInOrder 函数
// 查询帖子 ID 相关信息的列表（如分数，创建时间）
func TestGetPostIDsInOrder(t *testing.T) {
	// err
	// var err error

	// 模拟连接
	global.RDS = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		DB:       0,
		Password: "qaz123",
	})

	// // 模拟用户 ID 和帖子 ID
	// userID := "user2"
	// postID := "post2"

	// 测试正常情况，order 为 models.OrderTime创建时间
	p1 := &request.ParamPostList{
		CommunityID: 1,
		Order:       global.OrderTime,
		Page:        1,
		Size:        10,
	}
	ids1, err1 := GetPostIDsInOrder(context.Background(), p1)
	if err1 != nil {
		t.Errorf("Expected no error, got %v", err1)
	}
	if len(ids1) == 0 {
		t.Errorf("Expected some IDs, got none")
	}
	t.Logf("ids1: %v", ids1)

	// 测试正常情况，order 为 models.OrderScore投票分数
	p2 := &request.ParamPostList{
		Order: global.OrderScore,
		Page:  1,
		Size:  10,
	}
	ids2, err2 := GetPostIDsInOrder(context.Background(), p2)
	if err2 != nil {
		t.Errorf("Expected no error, got %v", err2)
	}
	if len(ids2) == 0 {
		t.Errorf("Expected some IDs, got none")
	}
	t.Logf("ids2: %v", ids2)

	// // 测试边界情况，page 为 0
	// p3 := &models.ParamPostList{
	// 	Order: models.OrderTime,
	// 	Page:  0,
	// 	Size:  10,
	// }
	// ids3, err3 := GetPostIDsInOrder(context.Background(), p3)
	// if err3 == nil {
	// 	t.Errorf("Expected an error for page 0, got none")
	// }
	// t.Logf("ids3: %v", ids3)

	// // 测试边界情况，size 为 0
	// p4 := &models.ParamPostList{
	// 	Order: models.OrderTime,
	// 	Page:  1,
	// 	Size:  0,
	// }
	// ids4, err4 := GetPostIDsInOrder(context.Background(), p4)
	// if err4 == nil {
	// 	t.Errorf("Expected an error for size 0, got none")
	// }
	// t.Logf("ids4: %v", ids4)
}

// TestGetPostVoteData 测试 GetPostVoteData 函数
// 查询每篇帖子的投赞成票的数据
func TestGetPostVoteData(t *testing.T) {
	// 模拟 Redis 客户端
	global.RDS = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		DB:       0,
		Password: "qaz123",
	})

	// 模拟上下文
	ctx := context.Background()

	// 测试用例 1：正常情况
	ids := []string{"post1", "post2"}
	// 执行函数
	data, err := GetPostVoteData(ctx, ids)

	// 断言结果
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	t.Logf("data: %v", data)

	// // 测试用例 2：空 ID 列表
	// ids = []string{}
	// expectedData = []int64{}

	// // 执行函数
	// data, err = GetPostVoteData(ctx, ids)

	// // 断言结果
	// if err!= nil {
	//     t.Errorf("Expected no error, got %v", err)
	// }
	// if len(data)!= len(expectedData) {
	//     t.Errorf("Expected data length %d, got %d", len(expectedData), len(data))
	// }

	// // 测试用例 3：错误的 ID
	// ids = []string{"post1", "post2", "invalid_id"}

	// // 执行函数
	// data, err = GetPostVoteData(ctx, ids)

	// // 断言结果
	// if err == nil {
	//     t.Errorf("Expected an error, got none")
	// }
}

// TestGetCommunityPostIDsInOrder 测试 GetCommunityPostIDsInOrder 函数
// 查询某社区的帖子 ID 列表
func TestGetCommunityPostIDsInOrder(t *testing.T) {
	// 模拟上下文
	ctx := context.Background()

	// 模拟 Redis 客户端
	global.RDS = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		DB:       0,
		Password: "qaz123",
	})
	defer global.RDS.Close()

	// 模拟参数
	p := &request.ParamPostList{
		CommunityID: 1,
		Order:       global.OrderScore,
		Page:        1,
		Size:        10,
	}

	// 测试正常情况
	ids, err := GetCommunityPostIDsInOrder(ctx, p)
	assert.NoError(t, err)
	// assert.NotEmpty(t, ids)
	t.Logf("ids: %v", ids)

	// // 测试不存在的社区 ID
	// p.CommunityID = 999
	// ids, err = GetCommunityPostIDsInOrder(ctx, p)
	// assert.Error(t, err)
	// assert.Empty(t, ids)

	// // 测试错误的排序方式
	// p.Order = "invalidOrder"
	// ids, err = GetCommunityPostIDsInOrder(ctx, p)
	// assert.Error(t, err)
	// assert.Empty(t, ids)

	// 测试分页和大小参数
	// p.Page = 2
	// p.Size = 5
	// ids, err = GetCommunityPostIDsInOrder(ctx, p)
	// assert.NoError(t, err)
	// assert.NotEmpty(t, ids)
}
