package demo

import (
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/demo/controller"

	"github.com/gin-gonic/gin"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> demo 模块路由")

	// 演示-GORM基本使用
	ormGroup := router.Group("/demo/orm")
	{
		ormGroup.GET("/list", controller.NewDemoORM.List)
		ormGroup.GET("/all", controller.NewDemoORM.All)
		ormGroup.GET("/:id", controller.NewDemoORM.Info)
		ormGroup.POST("", controller.NewDemoORM.Add)
		ormGroup.PUT("", controller.NewDemoORM.Edit)
		ormGroup.DELETE("/:id", controller.NewDemoORM.Remove)
		ormGroup.DELETE("/clean", controller.NewDemoORM.Clean)
	}
}
