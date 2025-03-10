package auth

import (
	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/middleware"
	"mask_api_gin/src/modules/auth/controller"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> authorization 模块路由")

	// 验证码操作
	{
		router.GET("/captcha-image",
			middleware.RateLimit(middleware.LimitOption{
				Time:  300,
				Count: 60,
				Type:  middleware.LIMIT_IP,
			}),
			controller.NewCaptcha.Image,
		)
	}

	// 账号身份操作
	{
		router.POST("/login",
			middleware.RateLimit(middleware.LimitOption{
				Time:  180,
				Count: 15,
				Type:  middleware.LIMIT_IP,
			}),
			controller.NewAccount.Login,
		)
		router.POST("/logout",
			middleware.RateLimit(middleware.LimitOption{
				Time:  120,
				Count: 15,
				Type:  middleware.LIMIT_IP,
			}),
			controller.NewAccount.Logout,
		)
		router.GET("/me",
			middleware.PreAuthorize(nil),
			controller.NewAccount.Me,
		)
		router.GET("/router",
			middleware.PreAuthorize(nil),
			controller.NewAccount.Router,
		)
	}

	// 账号注册操作
	{
		router.POST("/register",
			middleware.RateLimit(middleware.LimitOption{
				Time:  300,
				Count: 10,
				Type:  middleware.LIMIT_IP,
			}),
			controller.NewRegister.Register,
		)
	}
}
