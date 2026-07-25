package monitor

import (
	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/internal/web/middleware"
	"mask_api_gin/src/modules/monitor/controller"
	"mask_api_gin/src/modules/monitor/service"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> monitor 模块路由")

	// 启动时需要的初始参数
	InitLoad()

	// 服务器信息
	system := controller.NewSystem
	router.GET("/monitor/system",
		middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:system:info"}}),
		system.Info,
	)

	// 在线用户监控
	sysUserOnlineGroup := router.Group("/monitor/user-online")
	sysUserOnline := controller.NewSysUserOnline
	{
		sysUserOnlineGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:online:list"}}),
			sysUserOnline.List,
		)
		sysUserOnlineGroup.DELETE("/logout/:tokenId",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:online:logout"}}),
			middleware.OperateLog(middleware.OptionNew("在线用户监控", middleware.OPERATE_TYPE_FORCE)),
			sysUserOnline.Logout,
		)
	}

	// 缓存信息
	sysCacheGroup := router.Group("/monitor/cache")
	sysCache := controller.NewSysCache
	{
		sysCacheGroup.GET("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:cache:info"}}),
			sysCache.Info,
		)
		sysCacheGroup.GET("/names",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:cache:list"}}),
			sysCache.Names,
		)
		sysCacheGroup.GET("/keys",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:cache:list"}}),
			sysCache.Keys,
		)
		sysCacheGroup.GET("/value",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:cache:query"}}),
			sysCache.Value,
		)
		sysCacheGroup.DELETE("/names",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:cache:remove"}}),
			sysCache.CleanNames,
		)
		sysCacheGroup.DELETE("/keys",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:cache:remove"}}),
			sysCache.CleanKeys,
		)
		sysCacheGroup.DELETE("/value",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:cache:remove"}}),
			sysCache.CleanValue,
		)
	}

	// 调度任务信息
	sysJobGroup := router.Group("/monitor/job")
	sysJob := controller.NewSysJob
	{
		sysJobGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:job:list"}}),
			sysJob.List,
		)
		sysJobGroup.GET("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:job:query"}}),
			sysJob.Info,
		)
		sysJobGroup.POST("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:job:add"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.OPERATE_TYPE_INSERT)),
			sysJob.Add,
		)
		sysJobGroup.PUT("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:job:edit"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.OPERATE_TYPE_UPDATE)),
			sysJob.Edit,
		)
		sysJobGroup.DELETE("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:job:remove"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.OPERATE_TYPE_DELETE)),
			sysJob.Remove,
		)
		sysJobGroup.PUT("/status",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:job:status"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.OPERATE_TYPE_UPDATE)),
			sysJob.Status,
		)
		sysJobGroup.PUT("/run/:id",
			middleware.RepeatSubmit(10),
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:job:status"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.OPERATE_TYPE_UPDATE)),
			sysJob.Run,
		)
		sysJobGroup.PUT("/reset",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:job:status"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.OPERATE_TYPE_OTHER)),
			sysJob.Reset,
		)
		sysJobGroup.GET("/export",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:job:export"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务信息", middleware.OPERATE_TYPE_EXPORT)),
			sysJob.Export,
		)
	}

	// 调度任务日志信息
	sysJobLogGroup := router.Group("/monitor/job/log")
	sysJobLog := controller.NewSysJobLog
	{
		sysJobLogGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:job:list"}}),
			sysJobLog.List,
		)
		sysJobLogGroup.GET("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:job:query"}}),
			sysJobLog.Info,
		)
		sysJobLogGroup.DELETE("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:job:remove"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务日志信息", middleware.OPERATE_TYPE_DELETE)),
			sysJobLog.Remove,
		)
		sysJobLogGroup.DELETE("/clean",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:job:remove"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务日志信息", middleware.OPERATE_TYPE_CLEAN)),
			sysJobLog.Clean,
		)
		sysJobLogGroup.GET("/export",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"monitor:job:export"}}),
			middleware.OperateLog(middleware.OptionNew("调度任务日志信息", middleware.OPERATE_TYPE_EXPORT)),
			sysJobLog.Export,
		)
	}

}

// InitLoad 初始参数
func InitLoad() {
	// 启动时，初始化调度任务
	service.NewSysJob.Reset()
}
