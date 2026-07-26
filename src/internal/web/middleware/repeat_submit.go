package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/infra/cache"
	"mask_api_gin/src/internal/infra/ip2region"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/crypto"
	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
)

// repeatRecord 重复提交校验在 Redis 中存储的记录。
type repeatRecord struct {
	Time   int64  `json:"time"`   // 上次提交的 Unix 时间（秒）
	Params string `json:"params"` // 上次提交的参数 JSON
}

// RepeatSubmit 防止表单重复提交：在间隔时间内收到参数完全一致的请求视为重复。
// interval 单位秒，最小 5，默认 5。
func RepeatSubmit(interval int64) gin.HandlerFunc {
	if interval < 5 {
		interval = 5
	}

	return func(c *gin.Context) {
		params := reqctx.RequestParamsMap(c)
		paramsJSON, err := json.Marshal(params)
		if err != nil {
			logger.Errorf("repeat submit json marshal err: %v", err)
		}
		paramsStr := string(paramsJSON)

		funcName := c.HandlerName()
		clientIP := ip2region.ClientIP(c.ClientIP())
		repeatKey := constants.CACHE_REPEAT_SUBMIT + ":" + crypto.MD5(fmt.Sprintf("%s:%s", clientIP, funcName))

		// 比对上一次提交
		if raw, _ := cache.Get("", repeatKey); raw != "" {
			var prev repeatRecord
			if err := json.Unmarshal([]byte(raw), &prev); err != nil {
				logger.Errorf("repeat submit json unmarshal err: %v", err)
			}

			// 下次允许提交的时间戳（当前锁的过期时刻），单位毫秒
			resetAt := (prev.Time + interval) * 1000
			c.Header("X-RepeatSubmit-Reset", fmt.Sprintf("%d", resetAt))

			elapsed := time.Now().Unix() - prev.Time
			if elapsed < interval && prev.Params == paramsStr {
				c.JSON(200, resp.ErrMsg("不允许重复提交，请稍候再试"))
				c.Abort()
				return
			}
		}

		// 记录本次提交
		rec, err := json.Marshal(repeatRecord{
			Time:   time.Now().Unix(),
			Params: paramsStr,
		})
		if err != nil {
			logger.Errorf("repeat submit json marshal err: %v", err)
		}
		_ = cache.Set("", repeatKey, string(rec), time.Duration(interval)*time.Second)

		c.Next()
	}
}
