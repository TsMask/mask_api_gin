package monitor

import (
	"mask_api_gin/src/modules/monitor/controller"

	"github.com/gin-gonic/gin"
)

// 模块路由注册
func Setup(router *gin.Engine) {
	// 服务器监控
	serverGroup := router.Group("/monitor/server")
	{
		serverGroup.GET("/", controller.Server.Info)
	}

	// 调度任务日志信息
	jobLogGroup := router.Group("/monitor/jobLog")
	{
		// 导出调度任务日志信息
		jobLogGroup.POST("/export", controller.SysJobLog.Export)
		// 调度任务日志列表
		jobLogGroup.GET("/list", controller.SysJobLog.List)
		// 调度任务日志信息
		jobLogGroup.GET("/:jobLogId", controller.SysJobLog.Info)
		// 调度任务日志删除
		jobLogGroup.DELETE("/:jobLogIds", controller.SysJobLog.Remove)
		// 调度任务日志清空
		jobLogGroup.DELETE("/clean", controller.SysJobLog.Clean)
	}
}
