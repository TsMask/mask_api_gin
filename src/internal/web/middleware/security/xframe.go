package security

import "github.com/gin-gonic/gin"

// applyXFrame 在启用时下发 X-Frame-Options 头，控制页面是否允许被嵌入 frame/iframe。
func applyXFrame(c *gin.Context, cfg HeaderToggle) {
	if cfg.Enable {
		c.Header("x-frame-options", cfg.Value)
	}
}
