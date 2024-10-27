package keys

import (
	"context"
	"errors"
	"math"
	"time"

	"bluebell/global"

	"github.com/redis/go-redis/v9"
)

// 推荐阅读
// 基于用户投票的相关算法：http://www.ruanyifeng.com/blog/algorithm/

// 本项目使用简化版的投票分数
// 投一票就加432分   86400/200  --> 200张赞成票可以给你的帖子续一天

/* 投票的几种情况：
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
	oneWeekInSeconds = 7 * 24 * 3600 // 评分周期为: 一周
	scorePerVote     = 432           // 每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(ctx context.Context, postID string, communityID string) error {
	pipeline := global.RDS.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 更新：把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPF + communityID)
	pipeline.SAdd(ctx, cKey, postID)
	_, err := pipeline.Exec(ctx)
	return err
}

// VoteForPost 为帖子投票
// 投票值可以是 1（赞成）、-1（反对）或 0（取消投票）
func VoteForPost(ctx context.Context, userID, postID string, value float64) error {
	// 1. 判断投票限制
	// 去redis取帖子发布时间
	postTime := global.RDS.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Val()
	// 判断是否超过一周
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2和3需要放到一个pipeline事务中操作

	// 2. 更新贴子的分数
	// 先查当前用户给当前帖子的投票记录
	ov := global.RDS.ZScore(ctx, getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()

	// 更新：如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值
	pipeline := global.RDS.TxPipeline()
	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	// 3. 记录用户为该贴子投票的数据
	if value == 0 {
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value, // 赞成票还是反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec(ctx)
	return err
}
