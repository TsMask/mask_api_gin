package oauth2

import (
	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/middleware"
	"mask_api_gin/src/modules/oauth2/controller"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> oauth2 模块路由")

	// 客户端授权管理
	oauth2ClientGroup := router.Group("/oauth2/client")
	{
		oauth2ClientGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"matchRoles": {"admin"}}),
			controller.NewOauth2Client.List,
		)
		oauth2ClientGroup.GET("/:clientId",
			middleware.PreAuthorize(map[string][]string{"matchRoles": {"admin"}}),
			controller.NewOauth2Client.Info,
		)
		oauth2ClientGroup.POST("",
			middleware.PreAuthorize(map[string][]string{"matchRoles": {"admin"}}),
			middleware.OperateLog(middleware.OptionNew("客户端授权管理", middleware.BUSINESS_TYPE_INSERT)),
			controller.NewOauth2Client.Add,
		)
		oauth2ClientGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"matchRoles": {"admin"}}),
			middleware.OperateLog(middleware.OptionNew("客户端授权管理", middleware.BUSINESS_TYPE_UPDATE)),
			controller.NewOauth2Client.Edit,
		)
		oauth2ClientGroup.DELETE("/:id",
			middleware.PreAuthorize(map[string][]string{"matchRoles": {"admin"}}),
			middleware.OperateLog(middleware.OptionNew("客户端授权管理", middleware.BUSINESS_TYPE_DELETE)),
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
			controller.NewOauth2.Token,
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

}
