package keys

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"bluebell/controller/request"
	"bluebell/global"

	"github.com/redis/go-redis/v9"
)

func getIDsFormKey(ctx context.Context, key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	// 3. ZREVRANGE 按分数从大到小的顺序查询指定数量的元素
	return global.RDS.ZRevRange(ctx, key, start, end).Result()
}

// GetPostIDsInOrder 根据给定的参数从redis查询ids
func GetPostIDsInOrder(ctx context.Context, p *request.ParamPostList) ([]string, error) {
	// 从redis获取id
	// 1. 根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == global.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 2. 确定查询的索引起始点
	return getIDsFormKey(ctx, key, p.Page, p.Size)
}

// GetPostVoteData 根据ids（post id）查询每篇帖子的投赞成票的数据
func GetPostVoteData(ctx context.Context, ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPF + id)
	//	// 查找key中分数是1的元素的数量->统计每篇帖子的赞成票的数量
	//	v := global.RDS.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}
	// 使用pipeline一次发送多条命令,减少RTT
	pipeline := global.RDS.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(ctx, key, "1", "1") // 统计每篇帖子的赞成票的数量
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

func getOrderKey(order string) string {
	switch order {
	case global.OrderScore:
		return getRedisKey(KeyPostScoreZSet)
	default:
		return getRedisKey(KeyPostTimeZSet)
	}
}

// GetCommunityPostIDsInOrder 按社区查询ids
func GetCommunityPostIDsInOrder(ctx context.Context, p *request.ParamPostList) ([]string, error) {
	if p == nil || p.CommunityID <= 0 || (p.Order != global.OrderTime && p.Order != global.OrderScore) {
		return nil, errors.New("无效的参数")
	}

	if p.Page <= 0 || p.Size <= 0 {
		return nil, errors.New("无效的分页参数")
	}

	log.Printf("开始处理社区帖子列表请求，社区ID: %d, 排序方式: %s, 分页: %d, 每页大小: %d", p.CommunityID, p.Order, p.Page, p.Size)

	// 需要排序的key
	orderKey := getOrderKey(p.Order)

	// 使用 zinterstore 把分区的帖子set与帖子分数的 zset 生成一个新的zset
	// 针对新的zset 按之前的逻辑取数据

	// 社区的key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	if cKey == "" {
		// 没有对应的社区id
		return nil, errors.New("该社区不存在")
	}
	// 利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	log.Printf("社区键: %s, 排序键: %s, 合成键: %s", cKey, orderKey, key)

	if global.RDS.Exists(ctx, key).Val() < 1 {
		// 不存在，需要计算
		pipeline := global.RDS.Pipeline()
		// 计算帖子分数
		pipeline.ZInterStore(
			ctx,
			key,
			&redis.ZStore{
				Keys:      []string{cKey, orderKey},
				Aggregate: "MAX",
			},
		) // zinterstore 计算

		pipeline.Expire(ctx, key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec(ctx)
		if err != nil {
			log.Printf("执行 pipeline 失败: %v", err)
			return nil, err
		}
	}
	// 存在的话就直接根据key查询ids
	return getIDsFormKey(ctx, key, p.Page, p.Size)
}
