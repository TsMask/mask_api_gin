package security

import "github.com/gin-gonic/gin"

// applyNoOpen 在启用时下发 X-Download-Options: noopen，阻止 IE 直接打开下载文件。
func applyNoOpen(c *gin.Context, cfg Toggle) {
	if cfg.Enable {
		c.Header("x-download-options", "noopen")
	}
}
