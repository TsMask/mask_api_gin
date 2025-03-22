package oauth2

import (
	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/middleware"
	monitorController "mask_api_gin/src/modules/monitor/controller"
	"mask_api_gin/src/modules/oauth2/controller"
	oauth2Middleware "mask_api_gin/src/modules/oauth2/middleware"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> oauth2 模块路由")

	// 客户端授权管理
	oauth2ClientGroup := router.Group("/oauth2/client")
	{
		oauth2ClientGroup.GET("/list",
			middleware.PreAuthorize(nil),
			controller.NewOauth2Client.List,
		)
		oauth2ClientGroup.GET("/:clientId",
			middleware.PreAuthorize(nil),
			controller.NewOauth2Client.Info,
		)
		oauth2ClientGroup.POST("",
			middleware.PreAuthorize(nil),
			controller.NewOauth2Client.Add,
		)
		oauth2ClientGroup.PUT("",
			middleware.PreAuthorize(nil),
			controller.NewOauth2Client.Edit,
		)
		oauth2ClientGroup.DELETE("/:id",
			middleware.PreAuthorize(nil),
			controller.NewOauth2Client.Remove,
		)
	}

	// 授权认证
	oauth2Group := router.Group("/oauth2")
	{
		oauth2Group.GET("/authorize",
			middleware.RateLimit(middleware.LimitOption{
				Time:  60,
				Count: 30,
				Type:  middleware.LIMIT_IP,
			}),
			controller.NewOauth2.Authorize,
		)
		oauth2Group.POST("/token",
			middleware.RateLimit(middleware.LimitOption{
				Time:  180,
				Count: 15,
				Type:  middleware.LIMIT_IP,
			}),
			controller.NewOauth2.Tooken,
		)
		oauth2Group.POST("/refresh-token",
			middleware.RateLimit(middleware.LimitOption{
				Time:  60,
				Count: 5,
				Type:  middleware.LIMIT_IP,
			}),
			controller.NewOauth2.RefreshToken,
		)
	}

	// 系统监控
	{
		oauth2Group.GET("/monitor/system",
			oauth2Middleware.PreAuthorize(nil),
			monitorController.NewSystem.Info,
		)
	}

}
