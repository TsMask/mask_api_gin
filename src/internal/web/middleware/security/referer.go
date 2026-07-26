package security

import (
	"net/url"
	"slices"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/web/resp"
)

// checkReferer 在启用且策略为 referer 校验时，比对请求 Referer 的 host 与白名单。
// 仅对可能修改状态的请求方法生效（GET/HEAD/OPTIONS/TRACE 直接放行）。
func checkReferer(c *gin.Context, cfg CSRFConfig) {
	if !cfg.Enable || !cfg.CheckReferer {
		return
	}

	switch c.Request.Method {
	case "GET", "HEAD", "OPTIONS", "TRACE":
		return
	}

	referer := c.GetHeader("Referer")
	if referer == "" {
		c.AbortWithStatusJSON(200, resp.ErrMsg("无效 Referer 未知"))
		return
	}

	u, err := url.Parse(referer)
	if err != nil {
		c.AbortWithStatusJSON(200, resp.ErrMsg("无效 Referer 未知"))
		return
	}

	if slices.Contains(cfg.AllowList, u.Host) {
		return
	}
	c.AbortWithStatusJSON(200, resp.ErrMsg("无效 Referer "+u.Host))
}
