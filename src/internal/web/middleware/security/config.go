package security

import (
	"mask_api_gin/src/internal/config"
)

// Config 安全相关配置，在 Security() 注册时一次性解析。
type Config struct {
	CSRF    CSRFConfig
	Frame   HeaderToggle
	CSP     Toggle
	HSTS    HSTSConfig
	NoOpen  Toggle
	NoSniff Toggle
	XSS     HeaderToggle
}

// Toggle 仅含开关的配置。
type Toggle struct {
	Enable bool
}

// HeaderToggle 带可配置取值的开关型响应头。
type HeaderToggle struct {
	Enable bool
	Value  string
}

// HSTSConfig Strict-Transport-Security 配置。
type HSTSConfig struct {
	Enable            bool
	MaxAge            int
	IncludeSubdomains bool
}

// CSRFConfig 基于 Referer 的 CSRF 校验配置。
type CSRFConfig struct {
	Enable       bool
	CheckReferer bool // type ∈ {all, any, referer} 时启用
	AllowList    []string
}

// loadConfig 在注册期一次性读取 security.* 配置。
func loadConfig() Config {
	csrfType := config.GetStringOr("security.csrf.type", "")
	return Config{
		CSRF: CSRFConfig{
			Enable:       config.GetBoolOr("security.csrf.enable", false),
			CheckReferer: csrfType == "all" || csrfType == "any" || csrfType == "referer",
			AllowList:    config.GetStringSliceOr("security.csrf.refererWhiteList", nil),
		},
		Frame: HeaderToggle{
			Enable: config.GetBoolOr("security.xframe.enable", false),
			Value:  config.GetStringOr("security.xframe.value", "SAMEORIGIN"),
		},
		CSP: Toggle{Enable: config.GetBoolOr("security.csp.enable", false)},
		HSTS: HSTSConfig{
			Enable:            config.GetBoolOr("security.hsts.enable", false),
			MaxAge:            config.GetIntOr("security.hsts.maxAge", 31536000),
			IncludeSubdomains: config.GetBoolOr("security.hsts.includeSubdomains", false),
		},
		NoOpen:  Toggle{Enable: config.GetBoolOr("security.noopen.enable", false)},
		NoSniff: Toggle{Enable: config.GetBoolOr("security.nosniff.enable", false)},
		XSS: HeaderToggle{
			Enable: config.GetBoolOr("security.xssProtection.enable", false),
			Value:  config.GetStringOr("security.xssProtection.value", "1; mode=block"),
		},
	}
}
