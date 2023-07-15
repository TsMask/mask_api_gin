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

// RateLimit 请求限流
//
// 限流时间,单位秒 time
//
// 限流次数 count
//
// 限流条件类型 type
//
// 使用 USER 时，请在用户身份授权认证校验后使用
// 以便获取登录用户信息，无用户信息时默认为 GLOBAL
//
// 示例参数：map[string]int64{"time":5,"count":10,"type":IP}
//
// 参数表示：5秒内，最多请求10次，类型记录IP
func RateLimit(options map[string]int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 初始可选参数数据
		var limitTime int64 = 5
		var limitCount int64 = 10
		var limitType int64 = LIMIT_GLOBAL

		funcName := c.HandlerName()
		lastDotIndex := strings.LastIndex(funcName, "/")
		funcName = funcName[lastDotIndex+1:]
		var combinedKey string = cachekey.RATE_LIMIT_KEY + funcName

		if v, ok := options["time"]; ok {
			limitTime = v
		}
		if v, ok := options["count"]; ok {
			limitCount = v
		}
		if v, ok := options["type"]; ok {
			limitType = v
		}

		// 用户
		if limitType == LIMIT_USER {
			loginUser, err := ctxUtils.LoginUser(c)
			if err != nil {
				c.JSON(401, result.Err(map[string]interface{}{
					"code": 401,
					"msg":  err.Error(),
				}))
				c.Abort() // 停止执行后续的处理函数
				return
			}
			combinedKey = cachekey.RATE_LIMIT_KEY + loginUser.UserID + ":" + funcName
		}

		// IP
		if limitType == LIMIT_IP {
			combinedKey = cachekey.RATE_LIMIT_KEY + c.ClientIP() + ":" + funcName
		}

		// 在Redis查询并记录请求次数
		rateCount := redis.RateLimit(combinedKey, limitTime, limitCount)
		rateTime := redis.GetExpire(combinedKey)

		// 设置响应头中的限流声明字段
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limitCount))                        // 总请求数限制
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", limitCount-rateCount))          // 剩余可用请求数
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Unix()+int64(rateTime))) // 重置时间戳

		if rateCount >= limitCount {
			c.JSON(200, result.ErrMsg("访问过于频繁，请稍候再试"))
			c.Abort() // 停止执行后续的处理函数
			return
		}

		// 调用下一个处理程序
		c.Next()
	}
}
