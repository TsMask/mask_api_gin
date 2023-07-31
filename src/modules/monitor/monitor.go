package monitor

import (
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/middleware"
	"mask_api_gin/src/framework/middleware/operlog"
	"mask_api_gin/src/framework/middleware/repeat"
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

	// 服务器监控信息
	router.GET("/monitor/server",
		middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:server:info"}}),
		controller.NewServer.Info,
	)

	// 缓存监控信息
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
		sysCacheGroup.DELETE("/clearCacheName/:cacheName",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:remove"}}),
			controller.NewSysCache.ClearCacheName,
		)
		sysCacheGroup.DELETE("/clearCacheKey/:cacheName/:cacheKey",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:remove"}}),
			controller.NewSysCache.ClearCacheKey,
		)
		sysCacheGroup.DELETE("/clearCacheSafe",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:remove"}}),
			controller.NewSysCache.ClearCacheSafe,
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
			operlog.OperLog(operlog.OptionNew("调度任务日志信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.NewSysJobLog.Remove,
		)
		sysJobLogGroup.DELETE("/clean",
			repeat.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:remove"}}),
			operlog.OperLog(operlog.OptionNew("调度任务日志信息", operlog.BUSINESS_TYPE_CLEAN)),
			controller.NewSysJobLog.Clean,
		)
		sysJobLogGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:export"}}),
			operlog.OperLog(operlog.OptionNew("调度任务日志信息", operlog.BUSINESS_TYPE_EXPORT)),
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
			operlog.OperLog(operlog.OptionNew("调度任务信息", operlog.BUSINESS_TYPE_INSERT)),
			controller.NewSysJob.Add,
		)
		sysJobGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:edit"}}),
			operlog.OperLog(operlog.OptionNew("调度任务信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.NewSysJob.Edit,
		)
		sysJobGroup.DELETE("/:jobIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:remove"}}),
			operlog.OperLog(operlog.OptionNew("调度任务信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.NewSysJob.Remove,
		)
		sysJobGroup.PUT("/changeStatus",
			repeat.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:changeStatus"}}),
			operlog.OperLog(operlog.OptionNew("调度任务信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.NewSysJob.Status,
		)
		sysJobGroup.PUT("/run/:jobId",
			repeat.RepeatSubmit(10),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:changeStatus"}}),
			operlog.OperLog(operlog.OptionNew("调度任务信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.NewSysJob.Run,
		)
		sysJobGroup.PUT("/resetQueueJob",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:changeStatus"}}),
			operlog.OperLog(operlog.OptionNew("调度任务信息", operlog.BUSINESS_TYPE_CLEAN)),
			controller.NewSysJob.ResetQueueJob,
		)
		sysJobGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:job:export"}}),
			operlog.OperLog(operlog.OptionNew("调度任务信息", operlog.BUSINESS_TYPE_EXPORT)),
			controller.NewSysJob.Export,
		)
	}

	// 操作日志记录信息
	sysOperLogGroup := router.Group("/monitor/operlog")
	{
		sysOperLogGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:operlog:list"}}),
			controller.NewSysOperLog.List,
		)
		sysOperLogGroup.DELETE("/:operIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:operlog:remove"}}),
			operlog.OperLog(operlog.OptionNew("操作日志", operlog.BUSINESS_TYPE_DELETE)),
			controller.NewSysOperLog.Remove,
		)
		sysOperLogGroup.DELETE("/clean",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:operlog:remove"}}),
			operlog.OperLog(operlog.OptionNew("操作日志", operlog.BUSINESS_TYPE_CLEAN)),
			controller.NewSysOperLog.Clean,
		)
		sysOperLogGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:operlog:export"}}),
			operlog.OperLog(operlog.OptionNew("操作日志", operlog.BUSINESS_TYPE_EXPORT)),
			controller.NewSysOperLog.Export,
		)
	}

	// 登录访问信息
	sysLogininforGroup := router.Group("/monitor/logininfor")
	{
		sysLogininforGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:logininfor:list"}}),
			controller.NewSysLogininfor.List,
		)
		sysLogininforGroup.DELETE("/:infoIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:logininfor:remove"}}),
			operlog.OperLog(operlog.OptionNew("登录访问信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.NewSysLogininfor.Remove,
		)
		sysLogininforGroup.DELETE("/clean",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:logininfor:remove"}}),
			operlog.OperLog(operlog.OptionNew("登录访问信息", operlog.BUSINESS_TYPE_CLEAN)),
			controller.NewSysLogininfor.Clean,
		)
		sysLogininforGroup.PUT("/unlock/:userName",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:logininfor:unlock"}}),
			operlog.OperLog(operlog.OptionNew("登录访问信息", operlog.BUSINESS_TYPE_CLEAN)),
			controller.NewSysLogininfor.Unlock,
		)
		sysLogininforGroup.DELETE("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:logininfor:export"}}),
			operlog.OperLog(operlog.OptionNew("登录访问信息", operlog.BUSINESS_TYPE_EXPORT)),
			controller.NewSysLogininfor.Export,
		)
	}

	// 在线用户监控
	sysUserOnlineGroup := router.Group("/monitor/online")
	{
		sysUserOnlineGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:online:list"}}),
			controller.NewSysUserOnline.List,
		)
		sysUserOnlineGroup.DELETE("/:tokenId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:online:forceLogout"}}),
			controller.NewSysUserOnline.ForceLogout,
		)
	}
}

// InitLoad 初始参数
func InitLoad() {
	// 初始化定时任务处理
	processor.InitCronQueue()
	// 启动时，初始化调度任务
	service.NewSysJobImpl.ResetQueueJob()
}
