package middleware

import (
	"mask_api_gin/src/framework/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 用于记录请求处理时间
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 调用下一个处理程序
		c.Next()

		// 计算请求处理时间，并打印日志
		duration := time.Since(start)
		logger.Infof("Request processed in %v\n", duration)
	}
}
