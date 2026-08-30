package cache

import (
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// ZAddAfterRemoveByScore 清理 now 前过期的成员，再以 now+expiration 为成员 score 添加，并刷新集合过期时间。
func ZAddAfterRemoveByScore(source, key string, now time.Time, expiration time.Duration, members ...string) error {
	if len(members) == 0 {
		return fmt.Errorf("no members")
	}
	if expiration <= 0 {
		return fmt.Errorf("invalid expiration")
	}
	rdb := RDB(source)
	if rdb == nil {
		return fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()

	expireAt := now.Add(expiration)
	zMembers := make([]redis.Z, 0, len(members))
	for _, member := range members {
		zMembers = append(zMembers, redis.Z{
			Score:  float64(expireAt.UnixMilli()),
			Member: member,
		})
	}

	tx := rdb.TxPipeline()
	tx.ZRemRangeByScore(ctx, key, "-inf", strconv.FormatInt(now.UnixMilli(), 10))
	tx.ZAdd(ctx, key, zMembers...)
	tx.Expire(ctx, key, expiration)
	_, err := tx.Exec(ctx)
	return err
}

// ZRangeByScore 获取分数区间内的成员，边界使用 Redis 区间语法。
func ZRangeByScore(source, key, min, max string) ([]string, error) {
	rdb := RDB(source)
	if rdb == nil {
		return nil, fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()

	return rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
}

// ZRem 从有序集合移除成员。
func ZRem(source, key string, members ...string) error {
	if len(members) == 0 {
		return nil
	}
	rdb := RDB(source)
	if rdb == nil {
		return fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()

	return rdb.ZRem(ctx, key, members).Err()
}
