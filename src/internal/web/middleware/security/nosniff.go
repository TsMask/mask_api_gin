package security

import "github.com/gin-gonic/gin"

// applyNoSniff 在启用时下发 X-Content-Type-Options: nosniff。
// 对 3xx 重定向响应跳过，避免与 Location 头相互干扰。
func applyNoSniff(c *gin.Context, cfg Toggle) {
	if !cfg.Enable {
		return
	}
	if status := c.Writer.Status(); status >= 300 && status <= 308 {
		return
	}
	c.Header("x-content-type-options", "nosniff")
}
