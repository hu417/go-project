package vote

import (
	"context"
	"errors"
	"math"
	"strconv"
	"time"

	"blue-bell/global"

	"github.com/redis/go-redis/v9"
)

// 推荐阅读
// 基于用户投票的相关算法：http://www.ruanyifeng.com/blog/algorithm/

// 本项目使用简化版的投票分数
// 投一票就加432分   86400/200  --> 200张赞成票可以给你的帖子续一天

/*
	 投票的几种情况：
	   direction=1时，有两种情况：
	   	1. 之前没有投过票，现在投赞成票    --> 更新分数和投票记录  差值的绝对值：1  +432
	   	2. 之前投反对票，现在改投赞成票    --> 更新分数和投票记录  差值的绝对值：2  +432*2
	   direction=0时，有两种情况：
	   	1. 之前投过反对票，现在要取消投票  --> 更新分数和投票记录  差值的绝对值：1  +432
		2. 之前投过赞成票，现在要取消投票  --> 更新分数和投票记录  差值的绝对值：1  -432
	   direction=-1时，有两种情况：
	   	1. 之前没有投过票，现在投反对票    --> 更新分数和投票记录  差值的绝对值：1  -432
	   	2. 之前投赞成票，现在改投反对票    --> 更新分数和投票记录  差值的绝对值：2  -432*2

	   投票的限制：
	   每个贴子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	   	1. 到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
	   	2. 到期之后删除那个 KeyPostVotedZSetPF
*/
const (
	oneWeekSecond = 7 * 24 * 3600
	oneVoteScore  = 432
)

var ctx = context.Background()

// CreatePostTime 记录帖子的创建时间
func CreatePostTime(postID int64) error {
	pipeline := global.RedisCli.Pipeline()
	pipeline.ZAdd(ctx, getKey(KeyPostTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	pipeline.ZAdd(ctx, getKey(KeyPostScore), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	_, err := pipeline.Exec(ctx)
	return err
}

// CheckVoteTime 检查投票时间
func CheckVoteTime(postID int64) error {
	// 拿到帖子创建时间，减去现在的时间，如果时间大于一周，那么则超过投票时间
	postIDStr := strconv.Itoa(int(postID))
	createPostTime, err := global.RedisCli.ZScore(ctx, getKey(KeyPostTime), postIDStr).Result()
	if err != nil {
		return err
	}
	timeNow := float64(time.Now().Unix())
	if timeNow-createPostTime < oneWeekSecond {
		return nil
	}
	return errors.New("帖子投票时间已过")
}

func UserPostVote(userID int64, postID int64, direction int) error {
	// 查询之前用户对该帖子投票记录
	userIDStr := strconv.Itoa(int(userID))
	postIDStr := strconv.Itoa(int(postID))
	userVotePost := global.RedisCli.ZScore(ctx, getKey(KeyPostVote+postIDStr), userIDStr).Val()
	value := float64(direction)
	var op float64
	if value > userVotePost {
		op = 1
	} else {
		op = -1
	}
	// 之前与现在的票数之差的绝对值
	diff := math.Abs(value - userVotePost)
	// 更新帖子的分数
	pipeline := global.RedisCli.Pipeline()
	pipeline.ZIncrBy(ctx, getKey(KeyPostScore), oneVoteScore*op*diff, postIDStr)
	// 更新用户对该帖子票数
	if value == 0 {
		pipeline.ZRem(ctx, getKey(KeyPostVote+postIDStr), userIDStr)
	} else {
		pipeline.ZAdd(ctx, getKey(KeyPostVote+postIDStr), redis.Z{
			Score:  value,
			Member: userIDStr,
		})
	}
	_, err := pipeline.Exec(ctx)
	return err
}
