package job

import (
	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/infra/cron"
	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/internal/web/middleware"
	"mask_api_gin/src/modules/job/controller"
	"mask_api_gin/src/modules/job/processors/examples"
	"mask_api_gin/src/modules/job/service"
)

// Setup 模块注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> job 模块路由")

	// 启动时加载
	initLoad()

	// 调度任务信息
	sysJobGroup := router.Group("/api/v1/jobs")
	sysJob := controller.NewSysJob
	{
		sysJobGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:job:list"}}),
			sysJob.List,
		)
		sysJobGroup.GET("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:job:query"}}),
			sysJob.Info,
		)
		sysJobGroup.POST("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:job:add"}}),
			middleware.OperationLog("调度任务信息", middleware.OperationTypeInsert),
			sysJob.Add,
		)
		sysJobGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:job:edit"}}),
			middleware.OperationLog("调度任务信息", middleware.OperationTypeUpdate),
			sysJob.Edit,
		)
		sysJobGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:job:remove"}}),
			middleware.OperationLog("调度任务信息", middleware.OperationTypeDelete),
			sysJob.Remove,
		)
		sysJobGroup.PUT("/status",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:job:status"}}),
			middleware.OperationLog("调度任务信息", middleware.OperationTypeUpdate),
			sysJob.Status,
		)
		sysJobGroup.PUT("/run/:id",
			middleware.RepeatSubmit(10),
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:job:status"}}),
			middleware.OperationLog("调度任务信息", middleware.OperationTypeUpdate),
			sysJob.Run,
		)
		sysJobGroup.PUT("/reset",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:job:status"}}),
			middleware.OperationLog("调度任务信息", middleware.OperationTypeOther),
			sysJob.Reset,
		)
		sysJobGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:job:export"}}),
			middleware.OperationLog("调度任务信息", middleware.OperationTypeExport),
			sysJob.Export,
		)
	}

	// 调度任务日志信息
	sysJobLogGroup := router.Group("/api/v1/job/logs")
	sysJobLog := controller.NewSysJobLog
	{
		sysJobLogGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:job:log:list"}}),
			sysJobLog.List,
		)
		sysJobLogGroup.GET("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:job:log:query"}}),
			sysJobLog.Info,
		)
		sysJobLogGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:job:log:remove"}}),
			middleware.OperationLog("调度任务日志信息", middleware.OperationTypeDelete),
			sysJobLog.Remove,
		)
		sysJobLogGroup.DELETE("/clean",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:job:log:remove"}}),
			middleware.OperationLog("调度任务日志信息", middleware.OperationTypeClean),
			sysJobLog.Clean,
		)
		sysJobLogGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:job:log:export"}}),
			middleware.OperationLog("调度任务日志信息", middleware.OperationTypeExport),
			sysJobLog.Export,
		)
	}
}

// initLoad 初始加载
func initLoad() {
	// 注册任务执行日志持久化钩子
	// 将引擎发出的执行事件映射为调度任务日志并落库；引擎本身不依赖任何业务模型。
	// 需在调度任务初始化（Reset）之前调用。
	cron.SetOnJobLog(func(record cron.JobLogRecord) {
		service.NewSysJob.SaveJobLog(record)
	})
	cron.SetOnEngineError(func(event cron.EngineError) {
		logger.Errorf("cron 引擎错误 queue=%s id=%s err=%v", event.Queue, event.ID, event.Err)
	})

	// 先注册处理器，QueueNames 才能用于恢复数据库中的任务。
	cron.CreateQueue("test_job", examples.NewBasicProcessor)
	cron.CreateQueue("test_progress", examples.NewProgressProcessor)
	cron.CreateQueue("test_error", examples.NewErrorProcessor)

	// 启动时，初始化调度任务
	service.NewSysJob.Reset()
}
