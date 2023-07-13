package system

import (
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/controller"
	"mask_api_gin/src/modules/system/service"

	"github.com/gin-gonic/gin"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> system 模块路由")

	// 启动时需要的初始参数
	InitLoad()

	// 用户信息
	sysUserGroup := router.Group("/system/user")
	{
		sysUserGroup.GET("/list", controller.SysUser.List)
	}

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

// InitLoad 初始参数
func InitLoad() {
	// 启动时，刷新缓存-参数配置
	service.SysConfigImpl.ResetConfigCache()
}
