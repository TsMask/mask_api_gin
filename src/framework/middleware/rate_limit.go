package middleware

import (
	"fmt"
	"mask_api_gin/src/framework/constants/cachekey"
	"mask_api_gin/src/framework/redis"
	ctxUtils "mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/vo/result"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// 默认策略全局限流
	LIMIT_GLOBAL = 1

	// 根据请求者IP进行限流
	LIMIT_IP = 2

	// 根据用户ID进行限流
	LIMIT_USER = 3
)

// LimitOption 请求限流参数
type LimitOption struct {
	Time  int64 `json:"time"`  // 限流时间,单位秒
	Count int64 `json:"count"` // 限流次数
	Type  int64 `json:"type"`  // 限流条件类型,默认LIMIT_GLOBAL
}

// RateLimit 请求限流
//
// 示例参数：middleware.LimitOption{ Time:  5, Count: 10, Type:  middleware.LIMIT_IP }
//
// 参数表示：5秒内，最多请求10次，限制类型为 IP
//
// 使用 USER 时，请在用户身份授权认证校验后使用
// 以便获取登录用户信息，无用户信息时默认为 GLOBAL
func RateLimit(option LimitOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 初始可选参数数据
		if option.Time < 5 {
			option.Time = 5
		}
		if option.Count < 10 {
			option.Count = 10
		}
		if option.Type == 0 {
			option.Type = LIMIT_GLOBAL
		}

		// 获取执行函数名称
		funcName := c.HandlerName()
		lastDotIndex := strings.LastIndex(funcName, "/")
		funcName = funcName[lastDotIndex+1:]
		// 生成限流key
		var limitKey string = cachekey.RATE_LIMIT_KEY + funcName

		// 用户
		if option.Type == LIMIT_USER {
			loginUser, err := ctxUtils.LoginUser(c)
			if err != nil {
				c.JSON(401, result.Err(map[string]interface{}{
					"code": 401,
					"msg":  err.Error(),
				}))
				c.Abort() // 停止执行后续的处理函数
				return
			}
			limitKey = cachekey.RATE_LIMIT_KEY + loginUser.UserID + ":" + funcName
		}

		// IP
		if option.Type == LIMIT_IP {
			limitKey = cachekey.RATE_LIMIT_KEY + c.ClientIP() + ":" + funcName
		}

		// 在Redis查询并记录请求次数
		rateCount := redis.RateLimit(limitKey, option.Time, option.Count)
		rateTime := redis.GetExpire(limitKey)

		// 设置响应头中的限流声明字段
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", option.Count))                      // 总请求数限制
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", option.Count-rateCount))        // 剩余可用请求数
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Unix()+int64(rateTime))) // 重置时间戳

		if rateCount >= option.Count {
			c.JSON(200, result.ErrMsg("访问过于频繁，请稍候再试"))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		// 调用下一个处理程序
		c.Next()
	}
}
