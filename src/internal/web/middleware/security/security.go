package security

import "github.com/gin-gonic/gin"

// Security 聚合安全中间件：CSRF(Referer) 校验 + 一组安全响应头。
// 配置在注册时读取一次，运行期不再访问配置源。
func Security() gin.HandlerFunc {
	cfg := loadConfig()
	return func(c *gin.Context) {
		checkReferer(c, cfg.CSRF) // 命中策略时可能终止请求

		applyXFrame(c, cfg.Frame)
		applyCSP(c, cfg.CSP)
		applyHSTS(c, cfg.HSTS)
		applyNoOpen(c, cfg.NoOpen)
		applyNoSniff(c, cfg.NoSniff)
		applyXSS(c, cfg.XSS)

		c.Next()
	}
}
