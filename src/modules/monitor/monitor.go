package monitor

import (
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/middleware"
	"mask_api_gin/src/modules/monitor/controller"
	"mask_api_gin/src/modules/monitor/processor"
	"mask_api_gin/src/modules/monitor/service"

	"github.com/gin-gonic/gin"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> monitor 模块路由")

	// 启动时需要的初始参数
	InitLoad()

	// 服务器服务信息
	router.GET("/monitor/system-info",
		middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:system:info"}}),
		controller.NewSystemInfo.Info,
	)

	// 缓存服务信息
	sysCacheGroup := router.Group("/monitor/cache")
	{
		sysCacheGroup.GET("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:info"}}),
			controller.NewSysCache.Info,
		)
		sysCacheGroup.GET("/getNames",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:list"}}),
			controller.NewSysCache.Names,
		)
		sysCacheGroup.GET("/getKeys/:cacheName",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:list"}}),
			controller.NewSysCache.Keys,
		)
		sysCacheGroup.GET("/getValue/:cacheName/:cacheKey",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:query"}}),
			controller.NewSysCache.Value,
		)
		sysCacheGroup.DELETE("/cleanCacheName/:cacheName",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:remove"}}),
			controller.NewSysCache.CleanCacheName,
		)
		sysCacheGroup.DELETE("/cleanCacheKey/:cacheName/:cacheKey",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:remove"}}),
			controller.NewSysCache.CleanCacheKey,
		)
		sysCacheGroup.DELETE("/cleanCacheSafe",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:remove"}}),
			controller.NewSysCache.CleanCacheSafe,
		)
	}

	// 调度任务日志信息
	sysJobLogGroup := router.Group("/monitor/jobLog")
	{
		sysJobLogGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:list"}}),
			controller.NewSysJobLog.List,
		)
		sysJobLogGroup.GET("/:jobLogId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:query"}}),
			controller.NewSysJobLog.Info,
		)
		sysJobLogGroup.DELETE("/:jobLogIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:remove"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务日志信息", middleware.BusinessTypeDelete)),
			controller.NewSysJobLog.Remove,
		)
		sysJobLogGroup.DELETE("/clean",
			middleware.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:remove"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务日志信息", middleware.BusinessTypeClean)),
			controller.NewSysJobLog.Clean,
		)
		sysJobLogGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:export"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务日志信息", middleware.BusinessTypeExport)),
			controller.NewSysJobLog.Export,
		)
	}

	// 调度任务信息
	sysJobGroup := router.Group("/monitor/job")
	{
		sysJobGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:list"}}),
			controller.NewSysJob.List,
		)
		sysJobGroup.GET("/:jobId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:query"}}),
			controller.NewSysJob.Info,
		)
		sysJobGroup.POST("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:add"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.BusinessTypeInsert)),
			controller.NewSysJob.Add,
		)
		sysJobGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:edit"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.BusinessTypeUpdate)),
			controller.NewSysJob.Edit,
		)
		sysJobGroup.DELETE("/:jobIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:remove"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.BusinessTypeDelete)),
			controller.NewSysJob.Remove,
		)
		sysJobGroup.PUT("/changeStatus",
			middleware.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:changeStatus"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.BusinessTypeUpdate)),
			controller.NewSysJob.Status,
		)
		sysJobGroup.PUT("/run/:jobId",
			middleware.RepeatSubmit(10),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:changeStatus"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.BusinessTypeUpdate)),
			controller.NewSysJob.Run,
		)
		sysJobGroup.PUT("/resetQueueJob",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:changeStatus"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.BusinessTypeClean)),
			controller.NewSysJob.ResetQueueJob,
		)
		sysJobGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:export"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.BusinessTypeExport)),
			controller.NewSysJob.Export,
		)
	}

	// 在线用户监控
	sysUserOnlineGroup := router.Group("/monitor/online")
	{
		sysUserOnlineGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:online:list"}}),
			controller.NewSysUserOnline.List,
		)
		sysUserOnlineGroup.DELETE("/logout",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:online:logout"}}),
			controller.NewSysUserOnline.Logout,
		)
	}
}

// InitLoad 初始参数
func InitLoad() {
	// 初始化定时任务处理
	processor.InitCronQueue()
	// 启动时，初始化调度任务
	service.NewSysJob.Reset()
}
