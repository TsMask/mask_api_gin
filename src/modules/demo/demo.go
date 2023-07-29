package demo

import (
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/demo/controller"

	"github.com/gin-gonic/gin"
)

// 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> demo 模块路由")

	// 演示-GORM基本使用
	zzormGroup := router.Group("/zzorm")
	{
		zzormGroup.GET("/list", controller.ZzOrm.List)
		zzormGroup.GET("/all", controller.ZzOrm.All)
		zzormGroup.GET("/:id", controller.ZzOrm.Info)
		zzormGroup.POST("", controller.ZzOrm.Add)
		zzormGroup.PUT("", controller.ZzOrm.Edit)
		zzormGroup.DELETE("/:ids", controller.ZzOrm.Remove)
		zzormGroup.DELETE("/clean", controller.ZzOrm.Clean)
	}
}
