package monitor

import (
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/middleware"
	"mask_api_gin/src/framework/middleware/operlog"
	"mask_api_gin/src/modules/monitor/controller"

	"github.com/gin-gonic/gin"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> monitor 模块路由")

	// 启动时需要的初始参数
	InitLoad()

	// 缓存监控信息
	sysCacheGroup := router.Group("/monitor/cache")
	{
		sysCacheGroup.GET("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:info"}}),
			controller.SysCache.Info,
		)
		sysCacheGroup.GET("/getNames",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:list"}}),
			controller.SysCache.Names,
		)
		sysCacheGroup.GET("/getKeys/:cacheName",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:list"}}),
			controller.SysCache.Keys,
		)
		sysCacheGroup.GET("/getValue/:cacheName/:cacheKey",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:query"}}),
			controller.SysCache.Value,
		)
		sysCacheGroup.DELETE("/clearCacheName/:cacheName",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:remove"}}),
			controller.SysCache.ClearCacheName,
		)
		sysCacheGroup.DELETE("/clearCacheKey/:cacheName/:cacheKey",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:remove"}}),
			controller.SysCache.ClearCacheKey,
		)
		sysCacheGroup.DELETE("/clearCacheSafe",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:cache:remove"}}),
			controller.SysCache.ClearCacheSafe,
		)
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

	// 操作日志记录信息
	sysOperLogGroup := router.Group("/monitor/operlog")
	{
		sysOperLogGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:operlog:list"}}),
			controller.SysOperLog.List,
		)
		sysOperLogGroup.DELETE("/:operIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:operlog:remove"}}),
			operlog.OperLog(operlog.OptionNew("操作日志", operlog.BUSINESS_TYPE_DELETE)),
			controller.SysOperLog.Remove,
		)
		sysOperLogGroup.DELETE("/clean",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:operlog:remove"}}),
			operlog.OperLog(operlog.OptionNew("操作日志", operlog.BUSINESS_TYPE_CLEAN)),
			controller.SysOperLog.Clean,
		)
		sysOperLogGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:operlog:export"}}),
			operlog.OperLog(operlog.OptionNew("操作日志", operlog.BUSINESS_TYPE_EXPORT)),
			controller.SysOperLog.Export,
		)
	}

	// 登录访问信息
	sysLogininforGroup := router.Group("/monitor/logininfor")
	{
		sysLogininforGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:logininfor:list"}}),
			controller.SysLogininfor.List,
		)
		sysLogininforGroup.DELETE("/:infoIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:logininfor:remove"}}),
			operlog.OperLog(operlog.OptionNew("登录访问信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.SysLogininfor.Remove,
		)
		sysLogininforGroup.DELETE("/clean",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:logininfor:remove"}}),
			operlog.OperLog(operlog.OptionNew("登录访问信息", operlog.BUSINESS_TYPE_CLEAN)),
			controller.SysLogininfor.Clean,
		)
		sysLogininforGroup.PUT("/unlock/:userName",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:logininfor:unlock"}}),
			operlog.OperLog(operlog.OptionNew("登录访问信息", operlog.BUSINESS_TYPE_CLEAN)),
			controller.SysLogininfor.Unlock,
		)
		sysLogininforGroup.DELETE("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:logininfor:export"}}),
			operlog.OperLog(operlog.OptionNew("登录访问信息", operlog.BUSINESS_TYPE_EXPORT)),
			controller.SysLogininfor.Export,
		)
	}

	// 在线用户监控
	sysUserOnlineGroup := router.Group("/monitor/online")
	{
		sysUserOnlineGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:online:list"}}),
			controller.SysUserOnline.List,
		)
		sysUserOnlineGroup.DELETE("/:tokenId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:online:forceLogout"}}),
			controller.SysUserOnline.ForceLogout,
		)
	}

	// 服务器监控信息
	router.GET("/monitor/server",
		middleware.PreAuthorize(map[string][]string{"hasPerms": {"monitor:server:info"}}),
		controller.ServerController.Info,
	)
}

// InitLoad 初始参数
func InitLoad() {
}
