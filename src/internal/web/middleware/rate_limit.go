package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/infra/cache"
	"mask_api_gin/src/internal/infra/ip2region"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/crypto"
	"mask_api_gin/src/internal/web/reqctx"
	"mask_api_gin/src/internal/web/resp"
)

// RateLimitType 限流维度的类型。
type RateLimitType int

const (
	RateLimitGlobal RateLimitType = iota // 按处理函数全局限流（默认）
	RateLimitIP                          // 按客户端 IP 限流
	RateLimitUser                        // 按登录用户 ID 限流
)

// RateLimit 请求限流。
//
//	middleware.RateLimit(5, 10, middleware.RateLimitIP)
//
// 表示 5 秒内最多 10 次，按 IP 限流。seconds/count 低于下限(5/10)时自动取下限。
// 使用 RateLimitUser 时需在 AuthorizeUser 之后注册，此时若无登录用户信息则返回 401。
func RateLimit(seconds, count int64, limitType RateLimitType) gin.HandlerFunc {
	return func(c *gin.Context) {
		if seconds < 5 {
			seconds = 5
		}
		if count < 10 {
			count = 10
		}

		funcName := c.HandlerName()
		limitKey := buildLimitKey(limitType, funcName, c)
		if c.IsAborted() { // RateLimitUser 无登录用户时已写出 401
			return
		}

		// 记录并查询当前窗口内的请求次数
		rateCount, err := cache.RateLimit("", limitKey, time.Duration(seconds)*time.Second)
		if err != nil {
			rejectRateLimited(c)
			return
		}
		rateTTL, err := cache.GetExpire("", limitKey)
		if err != nil {
			rejectRateLimited(c)
			return
		}

		// 声明限流信息
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", count))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", count-rateCount))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Unix()+rateTTL))

		if rateCount >= count {
			rejectRateLimited(c)
			return
		}

		c.Next()
	}
}

// buildLimitKey 按限流维度生成 Redis key。
// RateLimitUser 要求已登录，否则直接写出 401 并返回空串（调用方应感知 c.Abort）。
func buildLimitKey(t RateLimitType, funcName string, c *gin.Context) string {
	switch t {
	case RateLimitUser:
		userId := reqctx.LoginUserToUserID(c)
		if userId <= 0 {
			c.JSON(401, resp.CodeMsg(resp.CODE_AUTH_INVALID, "invalid login user information"))
			c.Abort()
			return ""
		}
		return constants.CACHE_RATE_LIMIT + ":" + crypto.MD5(fmt.Sprintf("%d:%s", userId, funcName))
	case RateLimitIP:
		clientIP := ip2region.ClientIP(c.ClientIP())
		return constants.CACHE_RATE_LIMIT + ":" + crypto.MD5(fmt.Sprintf("%s:%s", clientIP, funcName))
	default: // RateLimitGlobal
		return constants.CACHE_RATE_LIMIT + ":" + crypto.MD5(funcName)
	}
}

// rejectRateLimited 统一写出限流响应。
func rejectRateLimited(c *gin.Context) {
	c.JSON(200, resp.CodeMsg(resp.CODE_RATELIMIT, resp.MSG_RATELIMIT))
	c.Abort()
}
