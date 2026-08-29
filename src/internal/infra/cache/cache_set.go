package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"mask_api_gin/src/internal/shared/logger"
)

// SetExpire 设置过期时间
func SetExpire(source, key string, expiration time.Duration) error {
	rdb := RDB(source)
	if rdb == nil {
		return fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()
	if err := rdb.Expire(ctx, key, expiration).Err(); err != nil {
		logger.Errorf("redis Expire err %v", err)
		return err
	}
	return nil
}

// Set 设置缓存数据
func Set(source, key string, value any, expiration time.Duration) error {
	rdb := RDB(source)
	if rdb == nil {
		return fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()
	if err := rdb.Set(ctx, key, value, expiration).Err(); err != nil {
		logger.Errorf("redis Set err %v", err)
		return err
	}
	return nil
}

// Del 删除单个
func Del(source string, key string) error {
	rdb := RDB(source)
	if rdb == nil {
		return fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()
	if err := rdb.Del(ctx, key).Err(); err != nil {
		logger.Errorf("redis Del err %v", err)
		return err
	}
	return nil
}

// DelKeys 删除多个
func DelKeys(source string, keys []string) error {
	if len(keys) == 0 {
		return fmt.Errorf("no keys")
	}

	rdb := RDB(source)
	if rdb == nil {
		return fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()
	if err := rdb.Del(ctx, keys...).Err(); err != nil {
		logger.Errorf("redis DelKeys err %v", err)
		return err
	}
	return nil
}

// SAdd 向集合添加成员，并刷新集合过期时间。
func SAdd(source, key string, expiration time.Duration, members ...string) error {
	if len(members) == 0 {
		return fmt.Errorf("no members")
	}
	rdb := RDB(source)
	if rdb == nil {
		return fmt.Errorf("redis not client")
	}
	ctx, cancel := DefaultCtx()
	defer cancel()
	tx := rdb.TxPipeline()
	tx.SAdd(ctx, key, members)
	tx.Expire(ctx, key, expiration)
	_, err := tx.Exec(ctx)
	return err
}

// SRem 从集合移除成员。
func SRem(source, key string, members ...string) error {
	if len(members) == 0 {
		return nil
	}
	rdb := RDB(source)
	if rdb == nil {
		return fmt.Errorf("redis not client")
	}
	ctx, cancel := DefaultCtx()
	defer cancel()
	return rdb.SRem(ctx, key, members).Err()
}

// RateLimit 限流计数；窗口创建时设置过期时间，后续请求不续期。
func RateLimit(source, limitKey string, expiration time.Duration) (int64, error) {
	if expiration <= 0 {
		return 0, fmt.Errorf("expiration must be positive")
	}
	rdb := RDB(source)
	if rdb == nil {
		return 0, fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()
	count, err := incrementWithInitialExpire(ctx, rdb, limitKey, expiration)
	if err != nil {
		logger.Errorf("redis RateLimit err %v", err)
		return 0, err
	}
	return count, nil
}

// SetWithExistingTTL 仅当 key 存在且有过期时间时更新值，保留原 TTL。
func SetWithExistingTTL(source, key string, value any) (int64, error) {
	rdb := RDB(source)
	if rdb == nil {
		return 0, fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()
	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if ttl == -2*time.Second {
		return -2, nil
	}
	if ttl == -1*time.Second {
		return -1, nil
	}

	tx := rdb.TxPipeline()
	tx.Set(ctx, key, value, time.Duration(ttl.Milliseconds())*time.Millisecond)
	if _, err := tx.Exec(ctx); err != nil {
		return 0, err
	}
	return int64(ttl.Seconds()), nil
}

// IncrWithExpire 递增计数；首次创建或发现无过期时间时补齐过期时间。
func IncrWithExpire(source, key string, expireSeconds int64) (int64, error) {
	if expireSeconds <= 0 {
		return 0, fmt.Errorf("expireSeconds must be positive")
	}
	rdb := RDB(source)
	if rdb == nil {
		return 0, fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()
	count, err := incrementWithInitialExpire(ctx, rdb, key, time.Duration(expireSeconds)*time.Second)
	if err != nil {
		logger.Errorf("redis IncrWithExpire err %v", err)
		return 0, err
	}
	return count, nil
}

// incrementWithInitialExpire 使用 WATCH/MULTI/EXEC 与 Go 函数实现“只初始化 TTL”的递增。
// 已有 TTL 的键执行 INCR，不刷新窗口；并发修改导致事务冲突时自动重试。
func incrementWithInitialExpire(ctx context.Context, rdb *redis.Client, key string, expiration time.Duration) (int64, error) {
	const maxRetries = 5
	for attempt := 0; attempt < maxRetries; attempt++ {
		var count int64
		err := rdb.Watch(ctx, func(tx *redis.Tx) error {
			current, getErr := tx.Get(ctx, key).Result()
			missing := errors.Is(getErr, redis.Nil)
			if getErr != nil && !missing {
				return getErr
			}

			ttl, ttlErr := tx.TTL(ctx, key).Result()
			if ttlErr != nil {
				return ttlErr
			}

			next := int64(1)
			if !missing {
				currentCount, parseErr := strconv.ParseInt(current, 10, 64)
				if parseErr != nil {
					return fmt.Errorf("redis key %s is not an integer", key)
				}
				next = currentCount + 1
			}

			_, execErr := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
				if ttl > 0 {
					pipe.Incr(ctx, key)
					return nil
				}
				pipe.Set(ctx, key, next, expiration)
				return nil
			})
			if execErr == nil {
				count = next
			}
			return execErr
		}, key)

		if errors.Is(err, redis.TxFailedErr) {
			continue
		}
		return count, err
	}
	return 0, fmt.Errorf("redis increment transaction conflict: %s", key)
}
