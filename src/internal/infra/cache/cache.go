package cache

import (
	"context"
	"fmt"
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
	datasource := config.GetStringMapOr("cache.dataSource", nil)
	for k, v := range datasource {
		client := v.(map[string]any)
		address := fmt.Sprintf("%s:%d", client["host"], client["port"])
		rdb := redis.NewClient(&redis.Options{
			Addr:     address,
			Password: client["password"].(string),
			DB:       client["db"].(int),
		})

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
	if source == "" {
		source = config.GetStringOr("cache.defaultDataSourceName", "")
	}
	rdbMapMu.RLock()
	defer rdbMapMu.RUnlock()
	return rdbMap[source]
}
