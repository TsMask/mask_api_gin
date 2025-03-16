package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// HGetAllToStruct 获得缓存数据
func HGetAllToStruct(source, key string, dest any) error {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return fmt.Errorf("redis not client")
	}

	ctx := context.Background()
	return rdb.HGetAll(ctx, key).Scan(dest)
}

// HGetAllToMap 获得缓存数据
func HGetAllToMap(source, key string, field string) (map[string]string, error) {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return nil, fmt.Errorf("redis not client")
	}

	ctx := context.Background()
	v, err := rdb.HGetAll(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("no keys")
	}
	if err != nil {
		return nil, err
	}
	return v, nil
}

// HGet 获得缓存数据
func HGet(source, key string, field string) (string, error) {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return "", fmt.Errorf("redis not client")
	}

	ctx := context.Background()
	v, err := rdb.HGet(ctx, key, field).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("no keys")
	}
	if err != nil {
		return "", err
	}
	return v, nil
}

// HSet 设置缓存数据
func HSet(source, key string, values any) (int64, error) {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return 0, fmt.Errorf("redis not client")
	}

	ctx := context.Background()
	v, err := rdb.HSet(ctx, key, values).Result()
	if errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("no keys")
	}
	if err != nil {
		return 0, err
	}
	return v, nil
}

// HExists 缓存数据是否存在字段
func HExists(source, key string, field string) (bool, error) {
	// 数据源
	rdb := RDB(source)
	if rdb == nil {
		return false, fmt.Errorf("redis not client")
	}

	ctx := context.Background()
	v, err := rdb.HExists(ctx, key, field).Result()
	if errors.Is(err, redis.Nil) {
		return false, fmt.Errorf("no keys")
	}
	if err != nil {
		return false, err
	}
	return v, nil
}
