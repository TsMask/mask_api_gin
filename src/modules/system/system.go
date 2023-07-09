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
		sysConfigGroup.GET("/list", controller.SysConfig.List)
		sysConfigGroup.GET("/:configId", controller.SysConfig.Info)
		sysConfigGroup.POST("/", controller.SysConfig.Add)
		sysConfigGroup.PUT("/", controller.SysConfig.Edit)
		sysConfigGroup.DELETE("/", controller.SysConfig.Remove)
		sysConfigGroup.POST("/export", controller.SysConfig.Export)
		sysConfigGroup.GET("/configKey/:configKey", controller.SysConfig.ConfigKey)
		sysConfigGroup.PUT("/refreshCache", controller.SysConfig.RefreshCache)
	}
}
