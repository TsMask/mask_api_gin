package security

import "github.com/gin-gonic/gin"

// applyXSS 在启用时下发 X-XSS-Protection 头，开启浏览器内置 XSS 过滤。
func applyXSS(c *gin.Context, cfg HeaderToggle) {
	if cfg.Enable {
		c.Header("x-xss-protection", cfg.Value)
	}
}
