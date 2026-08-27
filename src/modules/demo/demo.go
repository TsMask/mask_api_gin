package demo

import (
	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/modules/demo/controller"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> demo 模块路由")

	// 演示-GORM基本使用
	ormGroup := router.Group("/api/v1/demo/orm")
	demoORM := controller.NewDemoORM
	{
		ormGroup.GET("/list", demoORM.List)
		ormGroup.GET("/all", demoORM.All)
		ormGroup.GET("/:id", demoORM.Info)
		ormGroup.POST("", demoORM.Add)
		ormGroup.PUT("", demoORM.Edit)
		ormGroup.DELETE("/:id", demoORM.Remove)
		ormGroup.DELETE("/clean", demoORM.Clean)
	}
}
