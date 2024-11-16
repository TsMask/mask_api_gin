package system

import (
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/middleware"
	"mask_api_gin/src/modules/system/controller"
	"mask_api_gin/src/modules/system/service"

	"github.com/gin-gonic/gin"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> system 模块路由")

	// 启动时需要的初始参数
	InitLoad()

	// 参数配置信息
	sysConfigGroup := router.Group("/system/config")
	{
		sysConfigGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:list"}}),
			controller.NewSysConfig.List,
		)
		sysConfigGroup.GET("/:configId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:query"}}),
			controller.NewSysConfig.Info,
		)
		sysConfigGroup.POST("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:add"}}),
			middleware.OperateLog(middleware.OptionNew("参数配置信息", middleware.BUSINESS_TYPE_INSERT)),
			controller.NewSysConfig.Add,
		)
		sysConfigGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:edit"}}),
			middleware.OperateLog(middleware.OptionNew("参数配置信息", middleware.BUSINESS_TYPE_UPDATE)),
			controller.NewSysConfig.Edit,
		)
		sysConfigGroup.DELETE("/:configId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:remove"}}),
			middleware.OperateLog(middleware.OptionNew("参数配置信息", middleware.BUSINESS_TYPE_DELETE)),
			controller.NewSysConfig.Remove,
		)
		sysConfigGroup.PUT("/refresh",
			middleware.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:remove"}}),
			middleware.OperateLog(middleware.OptionNew("参数配置信息", middleware.BUSINESS_TYPE_OTHER)),
			controller.NewSysConfig.Refresh,
		)
		sysConfigGroup.GET("/config-key/:configKey", controller.NewSysConfig.ConfigKey)
		sysConfigGroup.GET("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:export"}}),
			middleware.OperateLog(middleware.OptionNew("参数配置信息", middleware.BUSINESS_TYPE_EXPORT)),
			controller.NewSysConfig.Export,
		)
	}

	// 部门信息
	sysDeptGroup := router.Group("/system/dept")
	{
		sysDeptGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:list"}}),
			controller.NewSysDept.List,
		)
		sysDeptGroup.GET("/:deptId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:query"}}),
			controller.NewSysDept.Info,
		)
		sysDeptGroup.POST("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:add"}}),
			middleware.OperateLog(middleware.OptionNew("部门信息", middleware.BUSINESS_TYPE_INSERT)),
			controller.NewSysDept.Add,
		)
		sysDeptGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:edit"}}),
			middleware.OperateLog(middleware.OptionNew("部门信息", middleware.BUSINESS_TYPE_UPDATE)),
			controller.NewSysDept.Edit,
		)
		sysDeptGroup.DELETE("/:deptId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:remove"}}),
			middleware.OperateLog(middleware.OptionNew("部门信息", middleware.BUSINESS_TYPE_DELETE)),
			controller.NewSysDept.Remove,
		)
		sysDeptGroup.GET("/list/exclude/:deptId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:list"}}),
			controller.NewSysDept.ExcludeChild,
		)
		sysDeptGroup.GET("/tree",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:list", "system:user:list"}}),
			controller.NewSysDept.Tree,
		)
		sysDeptGroup.GET("/tree/role/:roleId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:query", "system:user:edit"}}),
			controller.NewSysDept.TreeRole,
		)
	}

	// 字典数据信息
	sysDictDataGroup := router.Group("/system/dict/data")
	{
		sysDictDataGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:list"}}),
			controller.NewSysDictData.List,
		)
		sysDictDataGroup.GET("/:dataId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:query"}}),
			controller.NewSysDictData.Info,
		)
		sysDictDataGroup.POST("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:add"}}),
			middleware.OperateLog(middleware.OptionNew("字典数据信息", middleware.BUSINESS_TYPE_INSERT)),
			controller.NewSysDictData.Add,
		)
		sysDictDataGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:edit"}}),
			middleware.OperateLog(middleware.OptionNew("字典数据信息", middleware.BUSINESS_TYPE_UPDATE)),
			controller.NewSysDictData.Edit,
		)
		sysDictDataGroup.DELETE("/:dataId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:remove"}}),
			middleware.OperateLog(middleware.OptionNew("字典数据信息", middleware.BUSINESS_TYPE_DELETE)),
			controller.NewSysDictData.Remove,
		)
		sysDictDataGroup.GET("/type/:dictType",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:query"}}),
			controller.NewSysDictData.DictType,
		)
		sysDictDataGroup.GET("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:export"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.BUSINESS_TYPE_EXPORT)),
			controller.NewSysDictData.Export,
		)
	}

	// 字典类型信息
	sysDictTypeGroup := router.Group("/system/dict/type")
	{
		sysDictTypeGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:list"}}),
			controller.NewSysDictType.List,
		)
		sysDictTypeGroup.GET("/:dictId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:query"}}),
			controller.NewSysDictType.Info,
		)
		sysDictTypeGroup.POST("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:add"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.BUSINESS_TYPE_INSERT)),
			controller.NewSysDictType.Add,
		)
		sysDictTypeGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:edit"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.BUSINESS_TYPE_UPDATE)),
			controller.NewSysDictType.Edit,
		)
		sysDictTypeGroup.DELETE("/:dictId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:remove"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.BUSINESS_TYPE_DELETE)),
			controller.NewSysDictType.Remove,
		)
		sysDictTypeGroup.PUT("/refresh",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:remove"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.BUSINESS_TYPE_OTHER)),
			controller.NewSysDictType.Refresh,
		)
		sysDictTypeGroup.GET("/options",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:query"}}),
			controller.NewSysDictType.Options,
		)
		sysDictTypeGroup.GET("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:export"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.BUSINESS_TYPE_EXPORT)),
			controller.NewSysDictType.Export,
		)
	}

	// 系统登录日志信息
	sysLogLoginGroup := router.Group("/system/log/login")
	{
		sysLogLoginGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:login:list"}}),
			controller.NewSysLogLogin.List,
		)
		sysLogLoginGroup.DELETE("/clean",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:login:remove"}}),
			middleware.OperateLog(middleware.OptionNew("系统登录信息", middleware.BUSINESS_TYPE_CLEAN)),
			controller.NewSysLogLogin.Clean,
		)
		sysLogLoginGroup.PUT("/unlock/:userId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:login:unlock"}}),
			middleware.OperateLog(middleware.OptionNew("系统登录信息", middleware.BUSINESS_TYPE_OTHER)),
			controller.NewSysLogLogin.Unlock,
		)
		sysLogLoginGroup.GET("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:login:export"}}),
			middleware.OperateLog(middleware.OptionNew("系统登录信息", middleware.BUSINESS_TYPE_EXPORT)),
			controller.NewSysLogLogin.Export,
		)
	}

	// 操作日志记录信息
	sysLogOperateGroup := router.Group("/system/log/operate")
	{
		sysLogOperateGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:operate:list"}}),
			controller.NewSysLogOperate.List,
		)
		sysLogOperateGroup.DELETE("/clean",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:operate:remove"}}),
			middleware.OperateLog(middleware.OptionNew("操作日志", middleware.BUSINESS_TYPE_CLEAN)),
			controller.NewSysLogOperate.Clean,
		)
		sysLogOperateGroup.GET("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:operate:export"}}),
			middleware.OperateLog(middleware.OptionNew("操作日志", middleware.BUSINESS_TYPE_EXPORT)),
			controller.NewSysLogOperate.Export,
		)
	}

	// 菜单信息
	sysMenuGroup := router.Group("/system/menu")
	{
		sysMenuGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:list"}}),
			controller.NewSysMenu.List,
		)
		sysMenuGroup.GET("/:menuId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:query"}}),
			controller.NewSysMenu.Info,
		)
		sysMenuGroup.POST("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:add"}}),
			middleware.OperateLog(middleware.OptionNew("菜单信息", middleware.BUSINESS_TYPE_INSERT)),
			controller.NewSysMenu.Add,
		)
		sysMenuGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:edit"}}),
			middleware.OperateLog(middleware.OptionNew("菜单信息", middleware.BUSINESS_TYPE_UPDATE)),
			controller.NewSysMenu.Edit,
		)
		sysMenuGroup.DELETE("/:menuId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:remove"}}),
			middleware.OperateLog(middleware.OptionNew("菜单信息", middleware.BUSINESS_TYPE_DELETE)),
			controller.NewSysMenu.Remove,
		)
		sysMenuGroup.GET("/tree",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:list", "system:role:query"}}),
			controller.NewSysMenu.Tree,
		)
		sysMenuGroup.GET("/tree/role/:roleId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:list", "system:role:query"}}),
			controller.NewSysMenu.TreeRole,
		)
	}

	// 通知公告信息
	sysNoticeGroup := router.Group("/system/notice")
	{
		sysNoticeGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:list"}}),
			controller.NewSysNotice.List,
		)
		sysNoticeGroup.GET("/:noticeId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:query"}}),
			controller.NewSysNotice.Info,
		)
		sysNoticeGroup.POST("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:add"}}),
			middleware.OperateLog(middleware.OptionNew("通知公告信息", middleware.BUSINESS_TYPE_INSERT)),
			controller.NewSysNotice.Add,
		)
		sysNoticeGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:edit"}}),
			middleware.OperateLog(middleware.OptionNew("通知公告信息", middleware.BUSINESS_TYPE_UPDATE)),
			controller.NewSysNotice.Edit,
		)
		sysNoticeGroup.DELETE("/:noticeId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:remove"}}),
			middleware.OperateLog(middleware.OptionNew("通知公告信息", middleware.BUSINESS_TYPE_DELETE)),
			controller.NewSysNotice.Remove,
		)
	}

	// 岗位信息
	sysPostGroup := router.Group("/system/post")
	{
		sysPostGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:list"}}),
			controller.NewSysPost.List,
		)
		sysPostGroup.GET("/:postId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:query"}}),
			controller.NewSysPost.Info,
		)
		sysPostGroup.POST("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:add"}}),
			middleware.OperateLog(middleware.OptionNew("岗位信息", middleware.BUSINESS_TYPE_INSERT)),
			controller.NewSysPost.Add,
		)
		sysPostGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:edit"}}),
			middleware.OperateLog(middleware.OptionNew("岗位信息", middleware.BUSINESS_TYPE_UPDATE)),
			controller.NewSysPost.Edit,
		)
		sysPostGroup.DELETE("/:postId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:remove"}}),
			middleware.OperateLog(middleware.OptionNew("岗位信息", middleware.BUSINESS_TYPE_DELETE)),
			controller.NewSysPost.Remove,
		)
		sysPostGroup.GET("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:export"}}),
			middleware.OperateLog(middleware.OptionNew("岗位信息", middleware.BUSINESS_TYPE_EXPORT)),
			controller.NewSysPost.Export,
		)
	}

	// 个人信息
	sysProfileGroup := router.Group("/system/user/profile")
	{
		sysProfileGroup.GET("",
			middleware.PreAuthorize(nil),
			controller.NewSysProfile.Info,
		)
		sysProfileGroup.PUT("",
			middleware.PreAuthorize(nil),
			controller.NewSysProfile.UpdateProfile,
		)
		sysProfileGroup.PUT("/passwd",
			middleware.PreAuthorize(nil),
			middleware.OperateLog(middleware.OptionNew("个人信息", middleware.BUSINESS_TYPE_UPDATE)),
			controller.NewSysProfile.UpdatePasswd,
		)
	}

	// 角色信息
	sysRoleGroup := router.Group("/system/role")
	{
		sysRoleGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:list"}}),
			controller.NewSysRole.List,
		)
		sysRoleGroup.GET("/:roleId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:query"}}),
			controller.NewSysRole.Info,
		)
		sysRoleGroup.POST("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:add"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.BUSINESS_TYPE_INSERT)),
			controller.NewSysRole.Add,
		)
		sysRoleGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:edit"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.BUSINESS_TYPE_UPDATE)),
			controller.NewSysRole.Edit,
		)
		sysRoleGroup.DELETE("/:roleId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:remove"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.BUSINESS_TYPE_DELETE)),
			controller.NewSysRole.Remove,
		)
		sysRoleGroup.PUT("/status",
			middleware.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:edit"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.BUSINESS_TYPE_UPDATE)),
			controller.NewSysRole.Status,
		)
		sysRoleGroup.PUT("/data-scope",
			middleware.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.BUSINESS_TYPE_UPDATE)),
			controller.NewSysRole.DataScope,
		)
		sysRoleGroup.GET("/user/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:list"}}),
			controller.NewSysRole.UserAuthList,
		)
		sysRoleGroup.PUT("/user/auth",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.BUSINESS_TYPE_GRANT)),
			controller.NewSysRole.UserAuthChecked,
		)
		sysRoleGroup.GET("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:export"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.BUSINESS_TYPE_EXPORT)),
			controller.NewSysRole.Export,
		)
	}

	// 用户信息
	sysUserGroup := router.Group("/system/user")
	{
		sysUserGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:list"}}),
			controller.NewSysUser.List,
		)
		sysUserGroup.GET("/:userId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:query"}}),
			controller.NewSysUser.Info,
		)
		sysUserGroup.POST("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:add"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.BUSINESS_TYPE_INSERT)),
			controller.NewSysUser.Add,
		)
		sysUserGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.BUSINESS_TYPE_UPDATE)),
			controller.NewSysUser.Edit,
		)
		sysUserGroup.DELETE("/:userId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:remove"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.BUSINESS_TYPE_DELETE)),
			controller.NewSysUser.Remove,
		)
		sysUserGroup.PUT("/passwd",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:resetPwd"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.BUSINESS_TYPE_UPDATE)),
			controller.NewSysUser.Passwd,
		)
		sysUserGroup.PUT("/status",
			middleware.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.BUSINESS_TYPE_UPDATE)),
			controller.NewSysUser.Status,
		)
		sysUserGroup.GET("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:export"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.BUSINESS_TYPE_EXPORT)),
			controller.NewSysUser.Export,
		)
		sysUserGroup.GET("/import/template",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:import"}}),
			controller.NewSysUser.Template,
		)
		sysUserGroup.POST("/import",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:import"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.BUSINESS_TYPE_IMPORT)),
			controller.NewSysUser.Import,
		)
	}
}

// InitLoad 初始参数
func InitLoad() {
	// 启动时，刷新缓存-参数配置
	service.NewSysConfig.CacheClean("*")
	service.NewSysConfig.CacheLoad("*")
	// 启动时，刷新缓存-字典类型数据
	service.NewSysDictType.CacheClean("*")
	service.NewSysDictType.CacheLoad("*")
}
