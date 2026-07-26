package auth

import (
	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/internal/web/middleware"
	"mask_api_gin/src/modules/auth/controller"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> authorization 模块路由")

	// 验证码操作
	{
		captcha := controller.NewCaptcha
		router.GET("/captcha-image",
			middleware.RateLimit(300, 60, middleware.RateLimitIP),
			captcha.Image,
		)
	}

	// 账号身份操作
	{
		account := controller.NewAccount
		router.POST("/auth/login",
			middleware.RateLimit(180, 15, middleware.RateLimitIP),
			account.Login,
		)
		router.POST("/auth/logout",
			middleware.RateLimit(120, 15, middleware.RateLimitIP),
			account.Logout,
		)
		router.POST("/auth/refresh-token",
			middleware.RateLimit(60, 5, middleware.RateLimitIP),
			account.RefreshToken,
		)
		router.GET("/me",
			middleware.AuthorizeLogin(),
			account.Me,
		)
		router.GET("/router",
			middleware.AuthorizeLogin(),
			account.Router,
		)
	}

	// 账号注册操作
	{
		register := controller.NewRegister
		router.POST("/auth/register",
			middleware.RateLimit(300, 10, middleware.RateLimitIP),
			register.Register,
		)
	}
}
