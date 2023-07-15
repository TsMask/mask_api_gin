package redis

import (
	"context"
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis连接实例
var rdb *redis.Client

// 声明定义限流脚本命令
var rateLimitCommand = redis.NewScript(`
local key = KEYS[1]
local time = tonumber(ARGV[1])
local count = tonumber(ARGV[2])
local current = redis.call('get', key);
if current and tonumber(current) >= count then
	return tonumber(current);
end
current = redis.call('incr', key)
if tonumber(current) == 1 then 
	redis.call('expire', key, time)
end
return tonumber(current);`)

// 连接Redis实例
func Connect() {
	ctx := context.Background()
	client := config.Get("redis.client").(map[string]interface{})
	address := fmt.Sprintf("%s:%d", client["host"], client["port"])
	// 创建连接
	rdb = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: client["password"].(string),
		DB:       client["db"].(int),
	})
	// 测试数据库连接
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Panicf("Failed to ping redis: %v", err)
	}
	logger.Infof("Redis %s connection is successful.", pong)
}

// 关闭Redis实例
func Close() {
	if err := rdb.Close(); err != nil {
		logger.Panicf("fatal error db close: %s", err)
	}
}

// 获取键的剩余有效时间（秒）
func GetExpire(key string) float64 {
	ctx := context.Background()
	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		return 0
	}
	return ttl.Seconds()
}

// 获得缓存数据的key列表
func GetKeys(pattern string) []string {
	// 初始化变量
	var keys []string
	var cursor uint64 = 0
	ctx := context.Background()
	// 循环遍历获取匹配的键
	for {
		// 使用 SCAN 命令获取匹配的键
		batchKeys, nextCursor, err := rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			logger.Errorf("Failed to scan keys: %v", err)
			break
		}
		cursor = nextCursor
		keys = append(keys, batchKeys...)
		// 当 cursor 为 0，表示遍历完成
		if cursor == 0 {
			break
		}
	}
	return keys
}

// 批量获得缓存数据
func GetBatch(keys []string) []interface{} {
	if len(keys) == 0 {
		return []interface{}{}
	}
	// 获取缓存数据
	result, err := rdb.MGet(context.Background(), keys...).Result()
	if err != nil {
		logger.Errorf("Failed to get batch data: %v", err)
		return []interface{}{}
	}
	return result
}

// 获得缓存数据
func Get(key string) string {
	ctx := context.Background()
	value, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return ""
	}
	return value
}

// 判断是否存在
func Has(keys ...string) bool {
	ctx := context.Background()
	exists, err := rdb.Exists(ctx, keys...).Result()
	if err != nil {
		return false
	}
	return exists >= 1
}

// 设置缓存数据
func Set(key string, value interface{}) bool {
	ctx := context.Background()
	err := rdb.Set(ctx, key, value, 0).Err()
	return err == nil
}

// 设置缓存数据与过期时间
func SetByExpire(key string, value interface{}, expiration time.Duration) bool {
	ctx := context.Background()
	err := rdb.Set(ctx, key, value, expiration).Err()
	return err == nil
}

// 删除单个
func Del(key string) bool {
	ctx := context.Background()
	err := rdb.Del(ctx, key).Err()
	return err == nil
}

// 删除多个
func DelKeys(keys []string) bool {
	if len(keys) == 0 {
		return false
	}
	ctx := context.Background()
	err := rdb.Del(ctx, keys...).Err()
	return err == nil
}

// 限流查询并记录
func RateLimit(limitKey string, time, count int64) int64 {
	ctx := context.Background()
	result, err := rateLimitCommand.Run(ctx, rdb, []string{limitKey}, time, count).Result()
	if err != nil {
		logger.Errorf("redis lua script err %v", err)
		return 0
	}
	return result.(int64)
}