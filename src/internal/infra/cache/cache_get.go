package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"

	"mask_api_gin/src/internal/shared/logger"
)

// Info 获取redis服务信息
func Info(source string) map[string]map[string]string {
	infoObj := make(map[string]map[string]string)
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
		fields := strings.Fields(line)
		if strings.HasPrefix(line, "#") && len(fields) > 0 {
			label = strings.ToLower(fields[len(fields)-1])
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

// GetExpire 获取键的剩余有效时间（秒）
func GetExpire(source string, key string) (int64, error) {
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
	rdb := RDB(source)
	if rdb == nil {
		return keys, fmt.Errorf("redis not client")
	}

	var cursor uint64 = 0
	var count int64 = 100
	ctx, cancel := DefaultCtx()
	defer cancel()
	for {
		batchKeys, nextCursor, err := rdb.Scan(ctx, cursor, pattern, count).Result()
		if err != nil {
			logger.Errorf("failed to scan keys: %v", err)
			return keys, err
		}
		cursor = nextCursor
		keys = append(keys, batchKeys...)
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

	rdb := RDB(source)
	if rdb == nil {
		return result, fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()
	v, err := rdb.MGet(ctx, keys...).Result()
	if err != nil {
		logger.Errorf("failed to get batch data: %v", err)
		return result, err
	}
	return v, nil
}

// Get 获得缓存数据
func Get(source, key string) (string, error) {
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

// GetValue 按键的数据类型读取缓存值，并序列化为 JSON 字符串。
func GetValue(source, key string) (string, error) {
	rdb := RDB(source)
	if rdb == nil {
		return "", fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()

	keyType, err := rdb.Type(ctx, key).Result()
	if err != nil {
		return "", err
	}

	var value any
	switch keyType {
	case "string":
		value, err = rdb.Get(ctx, key).Result()
	case "set":
		value, err = rdb.SMembers(ctx, key).Result()
	case "list":
		value, err = rdb.LRange(ctx, key, 0, -1).Result()
	case "hash":
		value, err = rdb.HGetAll(ctx, key).Result()
	case "zset":
		value, err = rdb.ZRangeWithScores(ctx, key, 0, -1).Result()
	case "none":
		return "", fmt.Errorf("no keys")
	default:
		return "", fmt.Errorf("unsupported cache value type: %s", keyType)
	}
	if err != nil {
		return "", err
	}
	if text, ok := value.(string); ok {
		return text, nil
	}

	valueBytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(valueBytes), nil
}

// GetAndDel 原子获取并删除 key 的值（MULTI/EXEC 事务）
func GetAndDel(source, key string) (string, error) {
	rdb := RDB(source)
	if rdb == nil {
		return "", fmt.Errorf("redis not client")
	}

	ctx, cancel := DefaultCtx()
	defer cancel()
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
