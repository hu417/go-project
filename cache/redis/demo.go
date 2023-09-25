package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// 参考: https://redis.uptrace.dev/zh/guide/go-redis.html

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379",
		Password:     "123456", // 没有密码, 默认值
		DB:           0,        // 默认DB 0
		ReadTimeout:  -1,       // 从网络连接中读取数据超时时间
		WriteTimeout: -1,       // 把数据写入网络连接的超时时间
		PoolSize:     1000,     // 连接池最大连接数量
		MinIdleConns: 10,       // 连接池保持的最小空闲连接数
		MaxIdleConns: 100,      // 连接池保持的最大空闲连接数
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("连接redis出错, 错误信息: %v", err)
		return
	} else {
		fmt.Println("成功连接redis,", pong)
	}

	// ===== key val get set =====
	// RdbString(rdb)

	// // ===== list ======
	// RdbList(rdb)

	// // ===== 集合set =====
	// RdbSet(rdb)

	// // ===== 哈希hash操作 =====
	// RdbHash(rdb)

	// ===== 有序集合sort set ======
	RdbZset(rdb)

}

// 字符串
func RdbString(rdb *redis.Client) {
	// set
	err := rdb.Set(ctx, "guid", "89_1909979_2232118", 0).Err()
	if err != nil {
		fmt.Printf("redis set失败, 错误信息: %v", err)
	}

	// get
	val, err := rdb.Get(ctx, "guid").Result()
	if err != nil {
		fmt.Printf("redis get失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis get成功, key: %s, val:%s\n", "guid", val)
	}

	// del
	err = rdb.Del(ctx, "guid").Err()
	if err != nil {
		fmt.Printf("redis Del失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis Del成功, key: %s\n", "guid")
	}

	val, err = rdb.Get(ctx, "guidxx").Result()
	if err != nil {
		fmt.Printf("redis get失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis get成功, key: %s, val:%s\n", "guidxx", val)
	}

	// 过期时间
	// set key val NX, key不存在才set, 并设置指定过期时间
	err = rdb.SetNX(ctx, "task", "123", 10*time.Second).Err()
	if err != nil {
		fmt.Printf("redis SetNX失败, 错误信息: %v\n", err)
	}

	time.Sleep(time.Duration(1) * time.Second)

	// 获取过期剩余时间
	tm, err := rdb.TTL(ctx, "task").Result()
	if err != nil {
		fmt.Printf("redis TTL失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis TTL成功, key: %s, tm:%v\n", "task", tm)
	}

	val, err = rdb.Get(ctx, "task").Result()
	if err != nil {
		fmt.Printf("redis get失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis get成功, key: %s, val:%s\n", "task", val)
	}

	time.Sleep(time.Duration(11) * time.Second)

	val, err = rdb.Get(ctx, "task").Result()
	if err != nil {
		fmt.Printf("redis get失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis get成功, key: %s, val:%s\n", "task", val)
	}

	// 对key设置newValue这个值, 并且返回key原来的旧值。原子操作。
	oldVal, err := rdb.GetSet(ctx, "task", "456").Result()
	if err != nil {
		fmt.Printf("redis GetSet失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis GetSet成功, key: %s, oldVal:%s\n", "task", oldVal)
	}

	// 批量get set
	err = rdb.MSet(ctx, "key1", "val1", "key2", "val2", "key3", "val3").Err()
	if err != nil {
		fmt.Printf("redis MSet失败, 错误信息: %v\n", err)
	}

	vals, err := rdb.MGet(ctx, "key1", "key2", "key3").Result()
	if err != nil {
		fmt.Printf("redis MGet失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis MGet成功, vals: %v\n", vals)
	}

	// 自增自减, 原子操作
	err = rdb.Incr(ctx, "age").Err()
	if err != nil {
		fmt.Printf("redis Incr失败, 错误信息: %v\n", err)
	}
	// +5
	err = rdb.IncrBy(ctx, "age", 5).Err()
	if err != nil {
		fmt.Printf("redis IncrBy失败, 错误信息: %v\n", err)
	}

	err = rdb.Decr(ctx, "age").Err()
	if err != nil {
		fmt.Printf("redisDecr失败, 错误信息: %v\n", err)
	}
	// -3
	err = rdb.DecrBy(ctx, "age", 3).Err()
	if err != nil {
		fmt.Printf("redis DecrBy失败, 错误信息: %v\n", err)
	}

	// 2
	val, err = rdb.Get(ctx, "age").Result()
	if err != nil {
		fmt.Printf("redis get失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis get成功, key: %s, val:%s\n", "age", val)
	}

	// 对key设置过期时间
	rdb.Expire(ctx, "key1", 3*time.Second)

	tm, err = rdb.TTL(ctx, "key1").Result()
	if err != nil {
		fmt.Printf("redis TTL失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis TTL成功, key: %s, tm:%v\n", "key1", tm)
	}

	//批量删除
	err = rdb.Del(ctx, "key1", "key2", "key3").Err()
	if err != nil {
		fmt.Printf("redis Del失败, 错误信息: %v\n", err)
	}
}

// 列表
func RdbList(rdb *redis.Client) {
	// LPushX 从左往右插入
	// //仅当列表存在的时候才插入数据,此时列表不存在, 无法插入
	err := rdb.LPushX(ctx, "uids", 120).Err()
	if err != nil {
		fmt.Printf("redis LPushX失败, 错误信息: %v\n", err)
	}

	// /此时列表不存在, 依然可以插入
	err = rdb.LPush(ctx, "uids", 130).Err()
	if err != nil {
		fmt.Printf("redis LPush失败, 错误信息: %v\n", err)
	}

	// 批量插入
	err = rdb.LPushX(ctx, "uids", 130, 140, 154, 132).Err()
	if err != nil {
		fmt.Printf("redis LPush失败, 错误信息: %v\n", err)
	}

	// 返回全部数据
	vals2, err := rdb.LRange(ctx, "uids", 0, -1).Result()
	if err != nil {
		fmt.Printf("redis LRange失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis LRange成功, key: %s, vals:%v\n", "uids", vals2) // redis LRange成功, key: uids, vals:[132 154 140 130 130]
	}

	// 取队列[2,4]的数据, 即第3, 4, 5位置数据
	vals2, err = rdb.LRange(ctx, "uids", 2, 4).Result()
	if err != nil {
		fmt.Printf("redis LRange失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis LRange成功, key: %s, vals:%v\n", "uids", vals2) // redis LRange成功, key: uids, vals:[140 130 130]
	}

	// 返回队列长度
	llen, err := rdb.LLen(ctx, "uids").Result()
	if err != nil {
		fmt.Printf("redis LLen失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis LLen成功, key: %s, llen:%v\n", "uids", llen) // redis LLen成功, key: uids, llen:5
	}

	// 修改list中指定位置的值
	err = rdb.LSet(ctx, "uids", 2, 1000).Err()
	if err != nil {
		fmt.Printf("redis LSet失败, 错误信息: %v\n", err)
	}

	// 返回全部数据
	vals2, err = rdb.LRange(ctx, "uids", 0, -1).Result()
	if err != nil {
		fmt.Printf("redis LRange失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis LRange成功, key: %s, vals:%v\n", "uids", vals2) //redis LRange成功, key: uids, vals:[132 154 1000 130 130]
	}

	// 出队列, 左进右出
	err = rdb.RPop(ctx, "uids").Err()
	if err != nil {
		fmt.Printf("redis RPop失败, 错误信息: %v\n", err)
	}

	vals2, err = rdb.LRange(ctx, "uids", 0, -1).Result()
	if err != nil {
		fmt.Printf("redis LRange失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis LRange成功, key: %s, vals:%v\n", "uids", vals2)
	}

	// 出队列, 左进右出。没有就会阻塞, 可以设置阻塞超时值
	err = rdb.BRPop(ctx, time.Duration(1)*time.Second, "uids").Err()
	if err != nil {
		fmt.Printf("redis BRPop失败, 错误信息: %v\n", err)
	}

	// 删除一定位置范围内的值。删除count个key的list中值为value 的元素。如果出现重复元素, 仅删除1次, 也就是删除第一个
	err = rdb.LRem(ctx, "uids", 3, 130).Err()
	if err != nil {
		fmt.Printf("redis LRem失败, 错误信息: %v\n", err)
	}

	// 返回全部数据
	vals2, err = rdb.LRange(ctx, "uids", 0, -1).Result()
	if err != nil {
		fmt.Printf("redis LRange失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis LRange成功, key: %s, vals:%v\n", "uids", vals2) //redis LRange成功, key: uids, vals:[132 154 1000]
	}
}

// 集合
func RdbSet(rdb *redis.Client) {
	//redis集合特性: 元素无序且唯一

	// 批量入集合
	err := rdb.SAdd(ctx, "students", "Alice", "James", "James").Err()
	if err != nil {
		fmt.Printf("redis SAdd失败, 错误信息: %v\n", err)
	}

	// 获取集合大小
	size, err := rdb.SCard(ctx, "students").Result()
	if err != nil {
		fmt.Printf("redis SCard失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis SCard成功, key: %s, size:%v\n", "students", size) //redis SCard成功, key: students, size:2
	}

	// 返回集合所有元素
	sMem, err := rdb.SMembers(ctx, "students").Result()
	if err != nil {
		fmt.Printf("redis SMembers失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis SMembers成功, key: %s, size:%v\n", "students", sMem) //redis SMembers成功, key: students, size:[James Alice]
	}

	// 判断元素是否在集合中
	flag, err := rdb.SIsMember(ctx, "students", "James").Result()
	if err != nil {
		fmt.Printf("redis SIsMember失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis SIsMember成功, key: %s, size:%v\n", "students", flag) //redis SIsMember成功, key: students, size:true
	}

	// 删除集合元素
	err = rdb.SRem(ctx, "students", "Alice").Err()
	if err != nil {
		fmt.Printf("redis SRem失败, 错误信息: %v\n", err)
	}
}

// 哈希
func RdbHash(rdb *redis.Client) {
	// 多级嵌套HASH  China 是hash Guangdong 是字段名, Tencent是字段值
	err := rdb.HSet(ctx, "China", "Guangdong", "Tencent").Err()
	if err != nil {
		fmt.Printf("redisHSet失败, 错误信息: %v\n", err)
	}

	hvar, err := rdb.HGet(ctx, "China", "Guangdong").Result()
	if err != nil {
		fmt.Printf("redis HGet失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis HGet成功, hvar:%v\n", hvar)
	}

	err = rdb.HSet(ctx, "China", "Hanzhou", "Alibaba").Err()
	if err != nil {
		fmt.Printf("redis HSet失败, 错误信息: %v\n", err)
	}

	// 返回的是个map
	hvarAll, err := rdb.HGetAll(ctx, "China").Result()
	if err != nil {
		fmt.Printf("redis HGetAll失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis HGetAll成功, hvarAll:%v\n", hvarAll) //redis HGetAll成功, hvarAll:map[Guangdong:Tencent Hanzhou:Alibaba]
	}

	// 将map塞进hash
	batchData := make(map[string]interface{})
	batchData["username"] = "test"
	batchData["password"] = 123456
	err = rdb.HMSet(ctx, "users", batchData).Err()
	if err != nil {
		fmt.Printf("redis HMSet失败, 错误信息: %v\n", err)
	}

	hvarAll, err = rdb.HGetAll(ctx, "users").Result()
	if err != nil {
		fmt.Printf("redis HGetAll失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis HGetAll成功, hvarAll:%v\n", hvarAll) //redis HGetAll成功, hvarAll:map[password:123456 username:test]
	}

	// "Hanzhou"字段不存在才Set
	err = rdb.HSetNX(ctx, "China", "Hanzhou", "Netease").Err()
	if err != nil {
		fmt.Printf("redis HSetNX失败, 错误信息: %v\n", err)
	}

	// 对key的值+n
	count, err := rdb.HIncrBy(ctx, "users", "password", 10).Result()
	if err != nil {
		fmt.Printf("redis HIncrBy失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis HIncrBy成功, count:%v\n", count) // redis HIncrBy成功, count:123466
	}

	// 返回所有key,返回值是个string数组
	keys, err := rdb.HKeys(ctx, "China").Result()
	if err != nil {
		fmt.Printf("redisHKeys失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis HKeys成功, keys:%v\n", keys) //redis HKeys成功, keys:[Guangdong Hanzhou]
	}

	// 根据key, 查询hash的字段数量
	hlen, err := rdb.HLen(ctx, "China").Result()
	if err != nil {
		fmt.Printf("redis HLen失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis HLen成功, hlen:%v\n", hlen) // redis HLen成功, hlen:2
	}

	//删除多个key
	err = rdb.HDel(ctx, "China", "Hanzhou", "Guangdong").Err()
	if err != nil {
		fmt.Printf("redis HDel失败, 错误信息: %v\n", err)
	}
}

// 集合
func RdbZset(rdb *redis.Client) {
	/*
		// Z 表示已排序的集合成员
		type Z struct {
			Score  float64  // 分数
			Member interface{} // 元素名
		}
	*/
	zsetKey := "companys_rank"
	companys := []redis.Z{
		{Score: 100.0, Member: "Apple"},
		{Score: 90.0, Member: "MicroSoft"},
		{Score: 70.0, Member: "Amazon"},
		{Score: 87.0, Member: "Google"},
		{Score: 77.0, Member: "Facebook"},
		{Score: 67.0, Member: "Tesla"},
	}

	// 设置集合
	err := rdb.ZAdd(ctx, zsetKey, companys...).Err()
	if err != nil {
		fmt.Printf("redis ZAdd失败, 错误信息: %v\n", err)
	}
	// 指定元素自增 2
	err = rdb.ZIncrBy(ctx, zsetKey, 1, "Google").Err()
	if err != nil {
		fmt.Printf("redis ZIncrBy失败, 错误信息: %v\n", err)
	}

	//返回从0到-1位置的集合元素,  元素按分数从小到大排序
	rank, err := rdb.ZRange(ctx, zsetKey, 0, -1).Result()
	if err != nil {
		fmt.Printf("redis ZRange失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis ZRange成功, rank => %v\n", rank) //redis ZRange成功, rank:[Tesla Amazon Facebook Google MicroSoft Apple]
	}

	//返回top3
	rank, err = rdb.ZRange(ctx, zsetKey, -3, -1).Result()
	if err != nil {
		fmt.Printf("redis ZRange失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis ZRange成功, rank top3 => %+v\n", rank) //redis ZRange成功, rank:[Google MicroSoft Apple]
	}

	op := &redis.ZRangeBy{
		Min:    "50",  // 最小分数
		Max:    "200", // 最大分数
		Offset: 0,     // 类似sql的limit, 表示开始偏移量
		Count:  10,    // 一次返回多少数据
	}

	///根据分数范围返回集合元素, 实现排行榜, 取top n
	rank, err = rdb.ZRangeByScore(ctx, zsetKey, op).Result()
	if err != nil {
		fmt.Printf("redis ZRangeByScore失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis ZRangeByScore成功, 200 > rank > 50 => %v\n", rank) //redis ZRangeByScore成功, rank:[Google MicroSoft]
	}

	//根据元素名, 查询集合元素在集合中的排名, 从0开始算, 集合元素按分数从小到大排序
	rk, err := rdb.ZRank(ctx, zsetKey, "Apple").Result()
	if err != nil {
		fmt.Printf("redis ZRank失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis ZRank成功, rk[v] => %v\n", rk) // redis ZRank成功, rk:5
	}

	// 一次删除多个key
	err = rdb.ZRem(ctx, zsetKey, "Apple", "Amazon").Err()
	if err != nil {
		fmt.Printf("redis ZRem失败, 错误信息: %v\n", err)
	}

	//集合元素按分数排序, 从最低分到高分, 删除第0个元素到第1个元素。 这里相当于删除最低分的2个元素
	err = rdb.ZRemRangeByRank(ctx, zsetKey, 0, 1).Err()
	if err != nil {
		fmt.Printf("redis ZRemRangeByRank失败, 错误信息: %v\n", err)
	}

	rank, err = rdb.ZRange(ctx, zsetKey, 0, -1).Result()
	if err != nil {
		fmt.Printf("redis ZRange失败, 错误信息: %v\n", err)
	} else {
		fmt.Printf("redis ZRange成功, rank => %v\n", rank) //redis ZRange成功, rank:[Google MicroSoft]
	}
}
