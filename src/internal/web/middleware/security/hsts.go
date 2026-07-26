package security

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// applyHSTS 在启用时下发 Strict-Transport-Security 头，强制浏览器仅通过 HTTPS 访问。
func applyHSTS(c *gin.Context, cfg HSTSConfig) {
	if !cfg.Enable {
		return
	}
	value := fmt.Sprintf("max-age=%d", cfg.MaxAge)
	if cfg.IncludeSubdomains {
		value += "; includeSubdomains"
	}
	c.Header("strict-transport-security", value)
}
