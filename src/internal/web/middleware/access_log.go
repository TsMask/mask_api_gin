package middleware

import (
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/logger"
)

// AccessLog 记录每次请求的方法、路径、耗时与当前协程数，
// 仅在非生产环境注册，用于开发期观测。
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		logger.Infof("\n访问接口: %s %s\n总耗时: %v\n当前活跃的Goroutine数量: %d",
			c.Request.Method, c.Request.RequestURI, duration, runtime.NumGoroutine())
	}
}
