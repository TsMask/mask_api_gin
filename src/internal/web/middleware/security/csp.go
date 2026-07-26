package security

import (
	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/generate"
)

// applyCSP 在启用时下发一次性 CSP nonce，经 x-csp-nonce 头透出供模板拼装 CSP 指令。
func applyCSP(c *gin.Context, cfg Toggle) {
	if cfg.Enable {
		c.Header("x-csp-nonce", generate.Code(8))
	}
}
