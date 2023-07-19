package middleware

import (
	"mask_api_gin/src/framework/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// Report 请求响应日志
func Report() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 调用下一个处理程序
		c.Next()

		// 计算请求处理时间，并打印日志
		duration := time.Since(start)
		logger.Infof("%s %s report end=> %v", c.Request.Method, c.Request.RequestURI, duration)
	}
}
