package monitor

import (
	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/internal/web/middleware"
	"mask_api_gin/src/modules/monitor/controller"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> monitor 模块路由")

	// 服务器信息
	system := controller.NewSystem
	router.GET("/monitor/system",
		middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:system:info"}}),
		system.Info,
	)

	// 在线用户监控
	sysUserOnlineGroup := router.Group("/monitor/user-online")
	sysUserOnline := controller.NewSysUserOnline
	{
		sysUserOnlineGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:online:list"}}),
			sysUserOnline.List,
		)
		sysUserOnlineGroup.DELETE("/logout/:device_id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:online:logout"}}),
			middleware.OperateLog("在线用户监控", middleware.OperateTypeForce),
			sysUserOnline.Logout,
		)
	}

	// 缓存信息
	sysCacheGroup := router.Group("/monitor/cache")
	sysCache := controller.NewSysCache
	{
		sysCacheGroup.GET("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:cache:info"}}),
			sysCache.Info,
		)
		sysCacheGroup.GET("/names",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:cache:list"}}),
			sysCache.Names,
		)
		sysCacheGroup.GET("/keys",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:cache:list"}}),
			sysCache.Keys,
		)
		sysCacheGroup.GET("/value",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:cache:query"}}),
			sysCache.Value,
		)
		sysCacheGroup.DELETE("/names",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:cache:remove"}}),
			sysCache.CleanNames,
		)
		sysCacheGroup.DELETE("/keys",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:cache:remove"}}),
			sysCache.CleanKeys,
		)
		sysCacheGroup.DELETE("/value",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"monitor:cache:remove"}}),
			sysCache.CleanValue,
		)
	}
}
