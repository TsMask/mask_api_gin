package cache

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/shared/logger"
)

// Redis连接实例
var (
	rdbMapMu sync.RWMutex
	rdbMap   = make(map[string]*redis.Client)

	// defaultTimeout Redis 操作默认超时时间
	defaultTimeout = 5 * time.Second
)

// DefaultCtx 创建带默认超时的 context
func DefaultCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultTimeout)
}

// Connect 连接Redis实例
func Connect() {
	ctx := context.Background()
	// 读取数据源配置
	datasource := config.GetStringMapOr("redis.dataSource", nil)
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
		rdbMapMu.Lock()
		rdbMap[k] = rdb
		rdbMapMu.Unlock()
	}
}

// Close 关闭Redis实例
func Close() {
	rdbMapMu.Lock()
	defer rdbMapMu.Unlock()
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
		source = config.GetStringOr("redis.defaultDataSourceName", "")
	}
	rdbMapMu.RLock()
	defer rdbMapMu.RUnlock()
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

	ctx, cancel := DefaultCtx()
	defer cancel()
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

	ctx, cancel := DefaultCtx()
	defer cancel()
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

	ctx, cancel := DefaultCtx()
	defer cancel()
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

// Has 判断是否存在
func Has(source string, keys ...string) (int64, error) {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return 0, fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()
	exists, err := rdb.Exists(ctx, keys...).Result()
	if err != nil {
		return 0, err
	}
	return exists, nil
}

// SetExpire 设置过期时间
func SetExpire(source, key string, expiration time.Duration) error {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()
	err := rdb.Expire(ctx, key, expiration).Err()
	if err != nil {
		logger.Errorf("redis Expire err %v", err)
		return err
	}
	return nil
}

// GetExpire 获取键的剩余有效时间（秒）
func GetExpire(source string, key string) (int64, error) {
	// 数据源
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
	ctx, cancel := DefaultCtx()
	defer cancel()
	// 循环遍历获取匹配的键
	for {
		// 使用 SCAN 命令获取匹配的键
		batchKeys, nextCursor, err := rdb.Scan(ctx, cursor, pattern, count).Result()
		if err != nil {
			logger.Errorf("failed to scan keys: %v", err)
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
	ctx, cancel := DefaultCtx()
	defer cancel()
	v, err := rdb.MGet(ctx, keys...).Result()
	if err != nil || errors.Is(err, redis.Nil) {
		logger.Errorf("failed to get batch data: %v", err)
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

	ctx, cancel := DefaultCtx()
	defer cancel()
	v, err := rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("no keys")
	}
	if err != nil {
		return "", err
	}
	return v, nil
}

// Set 设置缓存数据
func Set(source, key string, value any, expiration time.Duration) error {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		logger.Errorf("redis Set err %v", err)
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

	// 数据源
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

// RateLimit 限流计数（MULTI/EXEC 原子执行 INCR + EXPIRE）
// 返回值: currentCount 当前计数, err 错误
func RateLimit(source, limitKey string, expiration time.Duration) (int64, error) {
	rdb := RDB(source)
	if rdb == nil {
		return 0, fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()

	// MULTI/EXEC 原子执行 INCR + EXPIRE
	tx := rdb.TxPipeline()
	incrCmd := tx.Incr(ctx, limitKey)
	tx.Expire(ctx, limitKey, expiration)
	if _, err := tx.Exec(ctx); err != nil {
		logger.Errorf("redis RateLimit tx err %v", err)
		return 0, err
	}

	return incrCmd.Val(), nil
}

// GetAndDel 原子获取并删除 key 的值（MULTI/EXEC 事务）
func GetAndDel(source, key string) (string, error) {
	rdb := RDB(source)
	if rdb == nil {
		return "", fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()

	// MULTI/EXEC 原子执行 GET + DEL，防止并发重复消费
	tx := rdb.TxPipeline()
	getCmd := tx.Get(ctx, key)
	tx.Del(ctx, key)
	if _, err := tx.Exec(ctx); err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("no keys")
		}
		return "", err
	}

	return getCmd.Val(), nil
}

// SetWithExistingTTL 仅当 key 存在且有过期时间时更新值，保留原 TTL（MULTI/EXEC 事务）
func SetWithExistingTTL(source, key string, value any) (int64, error) {
	rdb := RDB(source)
	if rdb == nil {
		return 0, fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()

	// 先读取 TTL 判断 key 状态
	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// key 不存在
	if ttl == -2*time.Second {
		return -2, nil
	}

	// key 存在但无过期时间
	if ttl == -1*time.Second {
		return -1, nil
	}

	// MULTI/EXEC 原子执行 SET + 保留原 TTL，避免 TTL 竞争
	ttlMs := int64(ttl.Milliseconds())
	tx := rdb.TxPipeline()
	tx.Set(ctx, key, value, time.Duration(ttlMs)*time.Millisecond)
	if _, err := tx.Exec(ctx); err != nil {
		return 0, err
	}

	return int64(ttl.Seconds()), nil
}

// IncrWithExpire 原子递增，首次创建时设置过期时间（MULTI/EXEC 事务）
func IncrWithExpire(source, key string, expireSeconds int64) (int64, error) {
	rdb := RDB(source)
	if rdb == nil {
		return 0, fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()

	// MULTI/EXEC 原子执行 INCR + EXPIRE，防止 INCR 成功但 EXPIRE 失败导致 key 永不过期
	tx := rdb.TxPipeline()
	incrCmd := tx.Incr(ctx, key)
	tx.Expire(ctx, key, time.Duration(expireSeconds)*time.Second)
	if _, err := tx.Exec(ctx); err != nil {
		logger.Errorf("redis IncrWithExpire tx err %v", err)
		return 0, err
	}

	return incrCmd.Val(), nil
}
