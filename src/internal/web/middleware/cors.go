package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/config"
)

// corsConfig 在注册期一次性解析的 CORS 配置，避免每个请求重复读取。
type corsConfig struct {
	Origin        string // 配置的允许来源；为空时回显请求 Origin
	Credentials   bool
	MaxAge        int
	AllowMethods  []string
	AllowHeaders  []string
	ExposeHeaders []string
}

func loadCorsConfig() corsConfig {
	return corsConfig{
		Origin:        config.GetStringOr("cors.origin", ""),
		Credentials:   config.GetBoolOr("cors.credentials", false),
		MaxAge:        config.GetIntOr("cors.maxAge", 0),
		AllowMethods:  config.GetStringSliceOr("cors.allowMethods", nil),
		AllowHeaders:  config.GetStringSliceOr("cors.allowHeaders", nil),
		ExposeHeaders: config.GetStringSliceOr("cors.exposeHeaders", nil),
	}
}

// resolveOrigin 计算写入 Access-Control-Allow-Origin 的值。
// 配置了 origin 时按配置返回（含 "*"）；否则回显请求来源。
func (cfg corsConfig) resolveOrigin(requestOrigin string) string {
	if cfg.Origin != "" {
		return cfg.Origin
	}
	return requestOrigin
}

// Cors 跨域中间件。配置在注册时读取一次。
func Cors() gin.HandlerFunc {
	cfg := loadCorsConfig()

	return func(c *gin.Context) {
		c.Header("Vary", "Origin")

		requestOrigin := c.GetHeader("Origin")
		if requestOrigin == "" {
			c.Next()
			return
		}

		c.Header("Access-Control-Allow-Origin", cfg.resolveOrigin(requestOrigin))
		if cfg.Credentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 预检请求
		if c.Request.Method == "OPTIONS" {
			if c.GetHeader("Access-Control-Request-Method") == "" {
				c.Next()
				return
			}

			if cfg.MaxAge > 10000 {
				c.Header("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))
			}

			methods := cfg.AllowMethods
			if len(methods) == 0 {
				methods = []string{"GET", "HEAD", "PUT", "POST", "DELETE", "PATCH"}
			}
			c.Header("Access-Control-Allow-Methods", strings.Join(methods, ","))

			if len(cfg.AllowHeaders) > 0 {
				c.Header("Access-Control-Allow-Headers", strings.Join(cfg.AllowHeaders, ","))
			}

			c.AbortWithStatus(204)
			return
		}

		// 暴露给前端可读的响应头
		if len(cfg.ExposeHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(cfg.ExposeHeaders, ","))
		}

		c.Next()
	}
}
