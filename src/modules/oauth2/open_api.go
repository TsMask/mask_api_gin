package oauth2

import (
	"github.com/gin-gonic/gin"

	"mask_api_gin/src/framework/middleware"
	monitorController "mask_api_gin/src/modules/monitor/controller"
)

// openAPI 客户端授权开放接口
func openAPI(router *gin.Engine) {
	openApiGroup := router.Group("/open-api")

	// 监控
	{
		openApiGroup.GET("/monitor/system",
			middleware.AuthorizeOauth2(nil),
			monitorController.NewSystem.Info,
		)
	}

}
