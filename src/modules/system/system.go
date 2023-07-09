package system

import (
	"mask_api_gin/src/modules/system/controller"

	"github.com/gin-gonic/gin"
)

// 模块路由注册
func Setup(router *gin.Engine) {
	// 参数配置信息
	sysConfigGroup := router.Group("/system/config")
	{
		// 导出参数配置信息
		sysConfigGroup.GET("/export", controller.SysConfig.Export)
	}
}
