package system

import (
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/middleware"
	"mask_api_gin/src/framework/middleware/operlog"
	"mask_api_gin/src/framework/middleware/repeat"
	"mask_api_gin/src/modules/system/controller"
	"mask_api_gin/src/modules/system/service"

	"github.com/gin-gonic/gin"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> system 模块路由")

	// 启动时需要的初始参数
	InitLoad()

	// 用户信息
	systemUserGroup := router.Group("/system/user")
	{
		systemUserGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:list"}}),
			controller.SysUser.List,
		)
		systemUserGroup.GET("/:userId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:query"}}),
			controller.SysUser.Info,
		)
		systemUserGroup.POST("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:add"}}),
			operlog.OperLog(operlog.OptionNew("用户信息", operlog.BUSINESS_TYPE_INSERT)),
			controller.SysUser.Add,
		)
		systemUserGroup.PUT("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			operlog.OperLog(operlog.OptionNew("用户信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysUser.Edit,
		)
		systemUserGroup.DELETE("/:userIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:remove"}}),
			operlog.OperLog(operlog.OptionNew("用户信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.SysUser.Remove,
		)
		systemUserGroup.PUT("/resetPwd",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:resetPwd"}}),
			operlog.OperLog(operlog.OptionNew("用户信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysUser.ResetPwd,
		)
		systemUserGroup.PUT("/changeStatus",
			repeat.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			operlog.OperLog(operlog.OptionNew("用户信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysUser.Status,
		)
	}

	// 参数配置信息
	systemConfigGroup := router.Group("/system/config")
	{
		systemConfigGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:list"}}),
			controller.SysConfig.List,
		)
		systemConfigGroup.GET("/:configId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:query"}}),
			controller.SysConfig.Info,
		)
		systemConfigGroup.POST("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:add"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_INSERT)),
			controller.SysConfig.Add,
		)
		systemConfigGroup.PUT("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:edit"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysConfig.Edit,
		)
		systemConfigGroup.DELETE("/:configIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:remove"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.SysConfig.Remove,
		)
		systemConfigGroup.PUT("/refreshCache",
			repeat.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:remove"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_CLEAN)),
			controller.SysConfig.RefreshCache,
		)
		systemConfigGroup.GET("/configKey/:configKey", controller.SysConfig.ConfigKey)
		systemConfigGroup.POST("/export", controller.SysConfig.Export)
	}
}

// InitLoad 初始参数
func InitLoad() {
	// 启动时，刷新缓存-参数配置
	service.SysConfigImpl.ResetConfigCache()
}
