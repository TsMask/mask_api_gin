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

	// 个人信息
	systemUserProfileGroup := router.Group("/system/user/profile")
	{
		systemUserProfileGroup.GET("/",
			middleware.PreAuthorize(nil),
			controller.SysProfile.Info,
		)
		systemUserProfileGroup.PUT("/",
			middleware.PreAuthorize(nil),
			controller.SysProfile.UpdateProfile,
		)
		systemUserProfileGroup.PUT("/updatePwd",
			middleware.PreAuthorize(nil),
			operlog.OperLog(operlog.OptionNew("个人信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysProfile.UpdatePwd,
		)
		systemUserProfileGroup.POST("/avatar",
			middleware.PreAuthorize(nil),
			operlog.OperLog(operlog.OptionNew("用户头像", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysProfile.Avatar,
		)
	}

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
		// 用户信息列表导入模板下载 TODO
		// 用户信息列表导入 TODO
		// 用户信息列表导出 TODO
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
		systemConfigGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:export"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_EXPORT)),
			controller.SysConfig.Export,
		)
	}

	// 通知公告信息
	systemNoticeGroup := router.Group("/system/notice")
	{
		systemNoticeGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:list"}}),
			controller.SysNotice.List,
		)
		systemNoticeGroup.GET("/:noticeId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:query"}}),
			controller.SysNotice.Info,
		)
		systemNoticeGroup.POST("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:add"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_INSERT)),
			controller.SysNotice.Add,
		)
		systemNoticeGroup.PUT("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:edit"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysNotice.Edit,
		)
		systemNoticeGroup.DELETE("/:noticeIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:remove"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.SysNotice.Remove,
		)
	}

	// 菜单信息
	systemMenuGroup := router.Group("/system/menu")
	{
		systemMenuGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:list"}}),
			controller.SysMenu.List,
		)
		systemMenuGroup.GET("/:menuId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:query"}}),
			controller.SysMenu.Info,
		)
		systemMenuGroup.POST("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:add"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_INSERT)),
			controller.SysMenu.Add,
		)
		systemMenuGroup.PUT("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:edit"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysMenu.Edit,
		)
		systemMenuGroup.DELETE("/:menuId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:remove"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.SysMenu.Remove,
		)
		systemMenuGroup.GET("/treeSelect",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:list"}}),
			controller.SysMenu.TreeSelect,
		)
		systemMenuGroup.GET("/roleMenuTreeSelect/:roleId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:list"}}),
			controller.SysMenu.RoleMenuTreeSelect,
		)
	}
}

// InitLoad 初始参数
func InitLoad() {
	// 启动时，刷新缓存-参数配置
	service.SysConfigImpl.ResetConfigCache()
}
