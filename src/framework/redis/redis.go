package redis

import (
	"context"
	"errors"
	"fmt"
	"mask_api_gin/src/framework/config"
	"mask_api_gin/src/framework/logger"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis连接实例
var rdbMap = make(map[string]*redis.Client)

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

// Connect 连接Redis实例
func Connect() {
	ctx := context.Background()
	// 读取数据源配置
	datasource := config.Get("redis.dataSource").(map[string]any)
	for k, v := range datasource {
		client := v.(map[string]any)
		// 创建连接
		address := fmt.Sprintf("%s:%d", client["host"], client["port"])
		rdb := redis.NewClient(&redis.Options{
			Addr:     address,
			Password: client["password"].(string),
			DB:       client["db"].(int),
		})
		// 测试数据库连接
		pong, err := rdb.Ping(ctx).Result()
		if err != nil {
			logger.Fatalf("Ping redis %s is %v", k, err)
		}
		logger.Infof("redis %s %d %s connection is successful.", k, client["db"].(int), pong)
		rdbMap[k] = rdb
	}
}

// Close 关闭Redis实例
func Close() {
	for _, rdb := range rdbMap {
		if err := rdb.Close(); err != nil {
			logger.Errorf("redis db close: %s", err)
		}
	}
}

// RDB 获取实例
func RDB(source string) *redis.Client {
	// 不指定时获取默认实例
	if source == "" {
		source = config.Get("redis.defaultDataSourceName").(string)
	}
	return rdbMap[source]
}

// Info 获取redis服务信息
func Info(source string) map[string]map[string]string {
	infoObj := make(map[string]map[string]string)
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return infoObj
	}

	ctx := context.Background()
	info, err := rdb.Info(ctx).Result()
	if err != nil {
		return infoObj
	}

	lines := strings.Split(info, "\r\n")
	label := ""
	for _, line := range lines {
		if strings.Contains(line, "#") {
			label = strings.Fields(line)[len(strings.Fields(line))-1]
			label = strings.ToLower(label)
			infoObj[label] = make(map[string]string)
			continue
		}
		kvArr := strings.Split(line, ":")
		if len(kvArr) >= 2 {
			key := strings.TrimSpace(kvArr[0])
			value := strings.TrimSpace(kvArr[len(kvArr)-1])
			infoObj[label][key] = value
		}
	}
	return infoObj
}

// KeySize 获取redis当前连接可用键Key总数信息
func KeySize(source string) int64 {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return 0
	}

	ctx := context.Background()
	size, err := rdb.DBSize(ctx).Result()
	if err != nil {
		return 0
	}
	return size
}

// CommandStats 获取redis命令状态信息
func CommandStats(source string) []map[string]string {
	statsObjArr := make([]map[string]string, 0)
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return statsObjArr
	}

	ctx := context.Background()
	commandstats, err := rdb.Info(ctx, "commandstats").Result()
	if err != nil {
		return statsObjArr
	}

	lines := strings.Split(commandstats, "\r\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "cmdstat_") {
			continue
		}
		kvArr := strings.Split(line, ":")
		key := kvArr[0]
		valueStr := kvArr[len(kvArr)-1]
		statsObj := make(map[string]string)
		statsObj["name"] = key[8:]
		statsObj["value"] = valueStr[6:strings.Index(valueStr, ",usec=")]
		statsObjArr = append(statsObjArr, statsObj)
	}
	return statsObjArr
}

// GetExpire 获取键的剩余有效时间（秒）
func GetExpire(source string, key string) (int64, error) {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return 0, fmt.Errorf("redis not client")
	}

	ctx := context.Background()
	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int64(ttl.Seconds()), nil
}

// GetKeys 获得缓存数据的key列表
func GetKeys(source string, pattern string) ([]string, error) {
	keys := make([]string, 0)
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return keys, fmt.Errorf("redis not client")
	}

	// 游标
	var cursor uint64 = 0
	var count int64 = 100
	ctx := context.Background()
	// 循环遍历获取匹配的键
	for {
		// 使用 SCAN 命令获取匹配的键
		batchKeys, nextCursor, err := rdb.Scan(ctx, cursor, pattern, count).Result()
		if err != nil {
			logger.Errorf("Failed to scan keys: %v", err)
			return keys, err
		}
		cursor = nextCursor
		keys = append(keys, batchKeys...)
		// 当 cursor 为 0，表示遍历完成
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

// GetBatch 批量获得缓存数据
func GetBatch(source string, keys []string) ([]any, error) {
	result := make([]any, 0)
	if len(keys) == 0 {
		return result, fmt.Errorf("not keys")
	}

	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return result, fmt.Errorf("redis not client")
	}

	// 获取缓存数据
	v, err := rdb.MGet(context.Background(), keys...).Result()
	if err != nil || errors.Is(err, redis.Nil) {
		logger.Errorf("Failed to get batch data: %v", err)
		return result, err
	}
	return v, nil
}

// Get 获得缓存数据
func Get(source, key string) (string, error) {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return "", fmt.Errorf("redis not client")
	}

	ctx := context.Background()
	v, err := rdb.Get(ctx, key).Result()
	if err != nil || errors.Is(err, redis.Nil) {
		return "", err
	}
	return v, nil
}

// GetHash 获得缓存数据Hash
func GetHash(source, key string) (map[string]string, error) {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return map[string]string{}, fmt.Errorf("redis not client")
	}

	ctx := context.Background()
	value, err := rdb.HGetAll(ctx, key).Result()
	if err != nil || errors.Is(err, redis.Nil) {
		return map[string]string{}, err
	}
	return value, nil
}

// Has 判断是否存在
func Has(source string, keys ...string) (int64, error) {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return 0, fmt.Errorf("redis not client")
	}

	ctx := context.Background()
	exists, err := rdb.Exists(ctx, keys...).Result()
	if err != nil {
		return 0, err
	}
	return exists, nil
}

// Set 设置缓存数据
func Set(source, key string, value any) error {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return fmt.Errorf("redis not client")
	}

	ctx := context.Background()
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		logger.Errorf("redis Set err %v", err)
		return err
	}
	return nil
}

// SetByExpire 设置缓存数据与过期时间
func SetByExpire(source, key string, value any, expiration time.Duration) error {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return fmt.Errorf("redis not client")
	}

	ctx := context.Background()
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		logger.Errorf("redis SetByExpire err %v", err)
		return err
	}
	return nil
}

// Del 删除单个
func Del(source string, key string) error {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return fmt.Errorf("redis not client")
	}

	ctx := context.Background()
	err := rdb.Del(ctx, key).Err()
	if err != nil {
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

	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return fmt.Errorf("redis not client")
	}

	ctx := context.Background()
	err := rdb.Del(ctx, keys...).Err()
	if err != nil {
		logger.Errorf("redis DelKeys err %v", err)
		return err
	}
	return nil
}

// RateLimit 限流查询并记录
func RateLimit(source, limitKey string, time, count int64) (int64, error) {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return 0, fmt.Errorf("redis not client")
	}

	ctx := context.Background()
	result, err := rateLimitCommand.Run(ctx, rdb, []string{limitKey}, time, count).Result()
	if err != nil {
		logger.Errorf("redis lua script err %v", err)
		return 0, err
	}
	return result.(int64), err
}
