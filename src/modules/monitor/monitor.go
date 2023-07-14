package monitor

import (
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/monitor/controller"

	"github.com/gin-gonic/gin"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> monitor 模块路由")

	// 启动时需要的初始参数
	InitLoad()

	// 服务器监控
	router.GET("/monitor/server", controller.Server.Info)

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

// InitLoad 初始参数
func InitLoad() {
}
