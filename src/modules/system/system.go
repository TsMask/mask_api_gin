package system

import (
	"github.com/gin-gonic/gin"

	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/internal/web/middleware"
	"mask_api_gin/src/modules/system/controller"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
)

// Setup 模块路由注册
func Setup(router *gin.Engine) {
	logger.Infof("开始加载 ====> system 模块路由")

	// 启动时需要的初始参数
	initLoad()

	// 参数配置信息
	sysConfigGroup := router.Group("/system/config")
	sysConfig := controller.NewSysConfig
	{
		sysConfigGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:config:list"}}),
			sysConfig.List,
		)
		sysConfigGroup.GET("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:config:query"}}),
			sysConfig.Info,
		)
		sysConfigGroup.POST("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:config:add"}}),
			middleware.OperateLog(middleware.OptionNew("参数配置信息", middleware.OPERATE_TYPE_INSERT)),
			sysConfig.Add,
		)
		sysConfigGroup.PUT("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:config:edit"}}),
			middleware.OperateLog(middleware.OptionNew("参数配置信息", middleware.OPERATE_TYPE_UPDATE)),
			sysConfig.Edit,
		)
		sysConfigGroup.DELETE("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:config:remove"}}),
			middleware.OperateLog(middleware.OptionNew("参数配置信息", middleware.OPERATE_TYPE_DELETE)),
			sysConfig.Remove,
		)
		sysConfigGroup.PUT("/refresh",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:config:remove"}}),
			middleware.OperateLog(middleware.OptionNew("参数配置信息", middleware.OPERATE_TYPE_OTHER)),
			sysConfig.Refresh,
		)
		sysConfigGroup.GET("/config-key/:configKey",
			middleware.RateLimit(middleware.LimitOption{
				Time:  120,
				Count: 15,
				Type:  middleware.LIMIT_IP,
			}),
			sysConfig.ConfigKey,
		)
		sysConfigGroup.GET("/export",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:config:export"}}),
			middleware.OperateLog(middleware.OptionNew("参数配置信息", middleware.OPERATE_TYPE_EXPORT)),
			sysConfig.Export,
		)
	}

	// 部门信息
	sysDeptGroup := router.Group("/system/dept")
	sysDept := controller.NewSysDept
	{
		sysDeptGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dept:list"}}),
			sysDept.List,
		)
		sysDeptGroup.GET("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dept:query"}}),
			sysDept.Info,
		)
		sysDeptGroup.POST("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dept:add"}}),
			middleware.OperateLog(middleware.OptionNew("部门信息", middleware.OPERATE_TYPE_INSERT)),
			sysDept.Add,
		)
		sysDeptGroup.PUT("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dept:edit"}}),
			middleware.OperateLog(middleware.OptionNew("部门信息", middleware.OPERATE_TYPE_UPDATE)),
			sysDept.Edit,
		)
		sysDeptGroup.DELETE("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dept:remove"}}),
			middleware.OperateLog(middleware.OptionNew("部门信息", middleware.OPERATE_TYPE_DELETE)),
			sysDept.Remove,
		)
		sysDeptGroup.GET("/list/exclude/:deptId",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dept:list"}}),
			sysDept.ExcludeChild,
		)
		sysDeptGroup.GET("/tree",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dept:list", "system:user:list"}}),
			sysDept.Tree,
		)
		sysDeptGroup.GET("/tree/role/:roleId",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dept:query", "system:user:edit"}}),
			sysDept.TreeRole,
		)
	}

	// 字典数据信息
	sysDictDataGroup := router.Group("/system/dict/data")
	sysDictData := controller.NewSysDictData
	{
		sysDictDataGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dict:list"}}),
			sysDictData.List,
		)
		sysDictDataGroup.GET(":id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dict:query"}}),
			sysDictData.Info,
		)
		sysDictDataGroup.POST("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dict:add"}}),
			middleware.OperateLog(middleware.OptionNew("字典数据信息", middleware.OPERATE_TYPE_INSERT)),
			sysDictData.Add,
		)
		sysDictDataGroup.PUT("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dict:edit"}}),
			middleware.OperateLog(middleware.OptionNew("字典数据信息", middleware.OPERATE_TYPE_UPDATE)),
			sysDictData.Edit,
		)
		sysDictDataGroup.DELETE("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dict:remove"}}),
			middleware.OperateLog(middleware.OptionNew("字典数据信息", middleware.OPERATE_TYPE_DELETE)),
			sysDictData.Remove,
		)
		sysDictDataGroup.GET("/type/:dictType",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dict:query"}}),
			sysDictData.DictType,
		)
		sysDictDataGroup.GET("/export",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dict:export"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.OPERATE_TYPE_EXPORT)),
			sysDictData.Export,
		)
	}

	// 字典类型信息
	sysDictTypeGroup := router.Group("/system/dict/type")
	sysDictType := controller.NewSysDictType
	{
		sysDictTypeGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dict:list"}}),
			sysDictType.List,
		)
		sysDictTypeGroup.GET("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dict:query"}}),
			sysDictType.Info,
		)
		sysDictTypeGroup.POST("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dict:add"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.OPERATE_TYPE_INSERT)),
			sysDictType.Add,
		)
		sysDictTypeGroup.PUT("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dict:edit"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.OPERATE_TYPE_UPDATE)),
			sysDictType.Edit,
		)
		sysDictTypeGroup.DELETE("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dict:remove"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.OPERATE_TYPE_DELETE)),
			sysDictType.Remove,
		)
		sysDictTypeGroup.PUT("/refresh",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dict:remove"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.OPERATE_TYPE_OTHER)),
			sysDictType.Refresh,
		)
		sysDictTypeGroup.GET("/options",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dict:query"}}),
			sysDictType.Options,
		)
		sysDictTypeGroup.GET("/export",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:dict:export"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.OPERATE_TYPE_EXPORT)),
			sysDictType.Export,
		)
	}

	// 系统登录日志信息
	sysLogLoginGroup := router.Group("/system/log/login")
	sysLogLogin := controller.NewSysLogLogin
	{
		sysLogLoginGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:log:login:list"}}),
			sysLogLogin.List,
		)
		sysLogLoginGroup.DELETE("/clean",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:log:login:remove"}}),
			middleware.OperateLog(middleware.OptionNew("系统登录信息", middleware.OPERATE_TYPE_CLEAN)),
			sysLogLogin.Clean,
		)
		sysLogLoginGroup.GET("/export",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:log:login:export"}}),
			middleware.OperateLog(middleware.OptionNew("系统登录信息", middleware.OPERATE_TYPE_EXPORT)),
			sysLogLogin.Export,
		)
	}

	// 操作日志记录信息
	sysLogOperateGroup := router.Group("/system/log/operate")
	sysLogOperate := controller.NewSysLogOperate
	{
		sysLogOperateGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:log:operate:list"}}),
			sysLogOperate.List,
		)
		sysLogOperateGroup.DELETE("/clean",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:log:operate:remove"}}),
			middleware.OperateLog(middleware.OptionNew("操作日志", middleware.OPERATE_TYPE_CLEAN)),
			sysLogOperate.Clean,
		)
		sysLogOperateGroup.GET("/export",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:log:operate:export"}}),
			middleware.OperateLog(middleware.OptionNew("操作日志", middleware.OPERATE_TYPE_EXPORT)),
			sysLogOperate.Export,
		)
	}

	// 菜单信息
	sysMenuGroup := router.Group("/system/menu")
	sysMenu := controller.NewSysMenu
	{
		sysMenuGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:menu:list"}}),
			sysMenu.List,
		)
		sysMenuGroup.GET("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:menu:query"}}),
			sysMenu.Info,
		)
		sysMenuGroup.POST("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:menu:add"}}),
			middleware.OperateLog(middleware.OptionNew("菜单信息", middleware.OPERATE_TYPE_INSERT)),
			sysMenu.Add,
		)
		sysMenuGroup.PUT("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:menu:edit"}}),
			middleware.OperateLog(middleware.OptionNew("菜单信息", middleware.OPERATE_TYPE_UPDATE)),
			sysMenu.Edit,
		)
		sysMenuGroup.DELETE("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:menu:remove"}}),
			middleware.OperateLog(middleware.OptionNew("菜单信息", middleware.OPERATE_TYPE_DELETE)),
			sysMenu.Remove,
		)
		sysMenuGroup.GET("/tree",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:menu:list", "system:role:query"}}),
			sysMenu.Tree,
		)
		sysMenuGroup.GET("/tree/role/:roleId",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:menu:list", "system:role:query"}}),
			sysMenu.TreeRole,
		)
	}

	// 通知公告信息
	sysNoticeGroup := router.Group("/system/notice")
	sysNotice := controller.NewSysNotice
	{
		sysNoticeGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:notice:list"}}),
			sysNotice.List,
		)
		sysNoticeGroup.GET("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:notice:query"}}),
			sysNotice.Info,
		)
		sysNoticeGroup.POST("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:notice:add"}}),
			middleware.OperateLog(middleware.OptionNew("通知公告信息", middleware.OPERATE_TYPE_INSERT)),
			sysNotice.Add,
		)
		sysNoticeGroup.PUT("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:notice:edit"}}),
			middleware.OperateLog(middleware.OptionNew("通知公告信息", middleware.OPERATE_TYPE_UPDATE)),
			sysNotice.Edit,
		)
		sysNoticeGroup.DELETE("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:notice:remove"}}),
			middleware.OperateLog(middleware.OptionNew("通知公告信息", middleware.OPERATE_TYPE_DELETE)),
			sysNotice.Remove,
		)
	}

	// 岗位信息
	sysPostGroup := router.Group("/system/post")
	sysPost := controller.NewSysPost
	{
		sysPostGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:post:list"}}),
			sysPost.List,
		)
		sysPostGroup.GET("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:post:query"}}),
			sysPost.Info,
		)
		sysPostGroup.POST("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:post:add"}}),
			middleware.OperateLog(middleware.OptionNew("岗位信息", middleware.OPERATE_TYPE_INSERT)),
			sysPost.Add,
		)
		sysPostGroup.PUT("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:post:edit"}}),
			middleware.OperateLog(middleware.OptionNew("岗位信息", middleware.OPERATE_TYPE_UPDATE)),
			sysPost.Edit,
		)
		sysPostGroup.DELETE("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:post:remove"}}),
			middleware.OperateLog(middleware.OptionNew("岗位信息", middleware.OPERATE_TYPE_DELETE)),
			sysPost.Remove,
		)
		sysPostGroup.GET("/export",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:post:export"}}),
			middleware.OperateLog(middleware.OptionNew("岗位信息", middleware.OPERATE_TYPE_EXPORT)),
			sysPost.Export,
		)
	}

	// 个人信息
	sysProfileGroup := router.Group("/system/user/profile")
	sysProfile := controller.NewSysProfile
	{
		sysProfileGroup.GET("",
			middleware.AuthorizeUser(nil),
			sysProfile.Info,
		)
		sysProfileGroup.PUT("",
			middleware.AuthorizeUser(nil),
			middleware.OperateLog(middleware.OptionNew("个人信息", middleware.OPERATE_TYPE_UPDATE)),
			sysProfile.UpdateProfile,
		)
		sysProfileGroup.PUT("/password",
			middleware.AuthorizeUser(nil),
			middleware.OperateLog(middleware.OptionNew("个人信息", middleware.OPERATE_TYPE_UPDATE)),
			sysProfile.UpdatePassword,
		)
	}

	// 角色信息
	sysRoleGroup := router.Group("/system/role")
	sysRole := controller.NewSysRole
	{
		sysRoleGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:role:list"}}),
			sysRole.List,
		)
		sysRoleGroup.GET("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:role:query"}}),
			sysRole.Info,
		)
		sysRoleGroup.POST("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:role:add"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.OPERATE_TYPE_INSERT)),
			sysRole.Add,
		)
		sysRoleGroup.PUT("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:role:edit"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.OPERATE_TYPE_UPDATE)),
			sysRole.Edit,
		)
		sysRoleGroup.DELETE("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:role:remove"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.OPERATE_TYPE_DELETE)),
			sysRole.Remove,
		)
		sysRoleGroup.PUT("/status",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:role:edit"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.OPERATE_TYPE_UPDATE)),
			sysRole.Status,
		)
		sysRoleGroup.PUT("/data-scope",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:user:edit"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.OPERATE_TYPE_UPDATE)),
			sysRole.DataScope,
		)
		sysRoleGroup.GET("/user/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:user:list"}}),
			sysRole.UserAuthList,
		)
		sysRoleGroup.PUT("/user/auth",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:user:edit"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.OPERATE_TYPE_GRANT)),
			sysRole.UserAuthChecked,
		)
		sysRoleGroup.GET("/export",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:user:export"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.OPERATE_TYPE_EXPORT)),
			sysRole.Export,
		)
	}

	// 用户信息
	sysUserGroup := router.Group("/system/user")
	sysUser := controller.NewSysUser
	{
		sysUserGroup.GET("/list",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:user:list"}}),
			sysUser.List,
		)
		sysUserGroup.GET("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:user:query"}}),
			sysUser.Info,
		)
		sysUserGroup.POST("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:user:add"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.OPERATE_TYPE_INSERT)),
			sysUser.Add,
		)
		sysUserGroup.PUT("",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:user:edit"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.OPERATE_TYPE_UPDATE)),
			sysUser.Edit,
		)
		sysUserGroup.DELETE("/:id",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:user:remove"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.OPERATE_TYPE_DELETE)),
			sysUser.Remove,
		)
		sysUserGroup.PUT("/password",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:user:resetPwd"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.OPERATE_TYPE_UPDATE)),
			sysUser.Password,
		)
		sysUserGroup.PUT("/status",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:user:edit"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.OPERATE_TYPE_UPDATE)),
			sysUser.Status,
		)
		sysUserGroup.PUT("/unlock/:userName",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:user:unlock"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.OPERATE_TYPE_CLEAN)),
			sysUser.Unlock,
		)
		sysUserGroup.GET("/export",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:user:export"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.OPERATE_TYPE_EXPORT)),
			sysUser.Export,
		)
		sysUserGroup.GET("/import/template",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:user:import"}}),
			sysUser.Template,
		)
		sysUserGroup.POST("/import",
			middleware.AuthorizeUser(map[string][]string{"hasPerms": {"system:user:import"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.OPERATE_TYPE_IMPORT)),
			sysUser.Import,
		)
	}
}

// InitLoad 初始参数
func initLoad() {
	// 启动时，刷新缓存-参数配置
	service.NewSysConfig.CacheClean("*")
	service.NewSysConfig.CacheLoad("*")
	// 启动时，刷新缓存-字典类型数据
	service.NewSysDictType.CacheClean("*")
	service.NewSysDictType.CacheLoad("*")
}

// RegisterOnOperateLog 注入操作日志持久化钩子。
//
// 将中间件发出的操作事件映射为操作日志并落库；中间件本身不依赖任何业务模型。
// 需在路由注册（触发操作日志）之前调用。
func RegisterOnOperateLog() {
	middleware.OnOperateLog = func(rec middleware.OperateLogRecord) {
		service.NewSysLogOperate.Insert(model.SysLogOperate{
			Title:           rec.Title,
			OperateType:     rec.OperateType,
			OperateFunction: rec.OperateFunction,
			OperateUrl:      rec.OperateUrl,
			OperateMethod:   rec.OperateMethod,
			OperateIp:       rec.OperateIp,
			OperateLocation: rec.OperateLocation,
			CreateBy:        rec.CreateBy,
			OperateParam:    rec.OperateParam,
			OperateStatus:   rec.OperateStatus,
			OperateMsg:      rec.OperateMsg,
			OperateCostTime: rec.OperateCostTime,
		})
	}
}
