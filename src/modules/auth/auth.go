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
			middleware.RateLimit(middleware.LimitOption{
				Time:  300,
				Count: 60,
				Type:  middleware.LIMIT_IP,
			}),
			captcha.Image,
		)
	}

	// 账号身份操作
	{
		account := controller.NewAccount
		router.POST("/auth/login",
			middleware.RateLimit(middleware.LimitOption{
				Time:  180,
				Count: 15,
				Type:  middleware.LIMIT_IP,
			}),
			account.Login,
		)
		router.POST("/auth/logout",
			middleware.RateLimit(middleware.LimitOption{
				Time:  120,
				Count: 15,
				Type:  middleware.LIMIT_IP,
			}),
			account.Logout,
		)
		router.POST("/auth/refresh-token",
			middleware.RateLimit(middleware.LimitOption{
				Time:  60,
				Count: 5,
				Type:  middleware.LIMIT_IP,
			}),
			account.RefreshToken,
		)
		router.GET("/me",
			middleware.AuthorizeUser(nil),
			account.Me,
		)
		router.GET("/router",
			middleware.AuthorizeUser(nil),
			account.Router,
		)
	}

	// 账号注册操作
	{
		register := controller.NewRegister
		router.POST("/auth/register",
			middleware.RateLimit(middleware.LimitOption{
				Time:  300,
				Count: 10,
				Type:  middleware.LIMIT_IP,
			}),
			register.Register,
		)
	}
}
