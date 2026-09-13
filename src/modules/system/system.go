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

	// 启动时加载
	initLoad()

	// 参数配置信息
	sysConfigGroup := router.Group("/api/v1/system/configs")
	sysConfig := controller.NewSysConfig
	{
		sysConfigGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:config:list"}}),
			sysConfig.List,
		)
		sysConfigGroup.GET("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:config:query"}}),
			sysConfig.Info,
		)
		sysConfigGroup.POST("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:config:add"}}),
			middleware.OperationLog("参数配置信息", middleware.OperationTypeInsert),
			sysConfig.Add,
		)
		sysConfigGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:config:edit"}}),
			middleware.OperationLog("参数配置信息", middleware.OperationTypeUpdate),
			sysConfig.Edit,
		)
		sysConfigGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:config:remove"}}),
			middleware.OperationLog("参数配置信息", middleware.OperationTypeDelete),
			sysConfig.Remove,
		)
		sysConfigGroup.PUT("/refresh",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:config:remove"}}),
			middleware.OperationLog("参数配置信息", middleware.OperationTypeOther),
			sysConfig.Refresh,
		)
		sysConfigGroup.GET("/config-key/:config_key",
			middleware.RateLimit(120, 15, middleware.RateLimitIP),
			sysConfig.ConfigKey,
		)
		sysConfigGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:config:export"}}),
			middleware.OperationLog("参数配置信息", middleware.OperationTypeExport),
			sysConfig.Export,
		)
	}

	// 部门信息
	sysDeptGroup := router.Group("/api/v1/system/depts")
	sysDept := controller.NewSysDept
	{
		sysDeptGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dept:list"}}),
			sysDept.List,
		)
		sysDeptGroup.GET("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dept:query"}}),
			sysDept.Info,
		)
		sysDeptGroup.POST("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dept:add"}}),
			middleware.OperationLog("部门信息", middleware.OperationTypeInsert),
			sysDept.Add,
		)
		sysDeptGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dept:edit"}}),
			middleware.OperationLog("部门信息", middleware.OperationTypeUpdate),
			sysDept.Edit,
		)
		sysDeptGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dept:remove"}}),
			middleware.OperationLog("部门信息", middleware.OperationTypeDelete),
			sysDept.Remove,
		)
		sysDeptGroup.GET("/list/exclude/:dept_id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dept:list"}}),
			sysDept.ExcludeChild,
		)
		sysDeptGroup.GET("/tree",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dept:list", "system:user:list"}}),
			sysDept.Tree,
		)
		sysDeptGroup.GET("/tree/role/:role_id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dept:query", "system:user:edit"}}),
			sysDept.TreeRole,
		)
	}

	// 字典数据信息
	sysDictDataGroup := router.Group("/api/v1/system/dict/datas")
	sysDictData := controller.NewSysDictData
	{
		sysDictDataGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:data:list"}}),
			sysDictData.List,
		)
		sysDictDataGroup.GET(":id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:data:query"}}),
			sysDictData.Info,
		)
		sysDictDataGroup.POST("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:data:add"}}),
			middleware.OperationLog("字典数据信息", middleware.OperationTypeInsert),
			sysDictData.Add,
		)
		sysDictDataGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:data:edit"}}),
			middleware.OperationLog("字典数据信息", middleware.OperationTypeUpdate),
			sysDictData.Edit,
		)
		sysDictDataGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:data:remove"}}),
			middleware.OperationLog("字典数据信息", middleware.OperationTypeDelete),
			sysDictData.Remove,
		)
		sysDictDataGroup.GET("/type/:dict_type",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:data:query"}}),
			sysDictData.DictType,
		)
		sysDictDataGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:data:export"}}),
			middleware.OperationLog("字典类型信息", middleware.OperationTypeExport),
			sysDictData.Export,
		)
	}

	// 字典类型信息
	sysDictTypeGroup := router.Group("/api/v1/system/dict/types")
	sysDictType := controller.NewSysDictType
	{
		sysDictTypeGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:type:list"}}),
			sysDictType.List,
		)
		sysDictTypeGroup.GET("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:type:query"}}),
			sysDictType.Info,
		)
		sysDictTypeGroup.POST("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:type:add"}}),
			middleware.OperationLog("字典类型信息", middleware.OperationTypeInsert),
			sysDictType.Add,
		)
		sysDictTypeGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:type:edit"}}),
			middleware.OperationLog("字典类型信息", middleware.OperationTypeUpdate),
			sysDictType.Edit,
		)
		sysDictTypeGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:type:remove"}}),
			middleware.OperationLog("字典类型信息", middleware.OperationTypeDelete),
			sysDictType.Remove,
		)
		sysDictTypeGroup.PUT("/refresh",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:type:remove"}}),
			middleware.OperationLog("字典类型信息", middleware.OperationTypeOther),
			sysDictType.Refresh,
		)
		sysDictTypeGroup.GET("/options",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:type:query"}}),
			sysDictType.Options,
		)
		sysDictTypeGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:type:export"}}),
			middleware.OperationLog("字典类型信息", middleware.OperationTypeExport),
			sysDictType.Export,
		)
	}

	// 系统登录日志信息
	sysLoginLogGroup := router.Group("/api/v1/system/logs/login")
	sysLoginLog := controller.NewSysLoginLog
	{
		sysLoginLogGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:login:log:list"}}),
			sysLoginLog.List,
		)
		sysLoginLogGroup.DELETE("/clean",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:login:log:remove"}}),
			middleware.OperationLog("系统登录信息", middleware.OperationTypeClean),
			sysLoginLog.Clean,
		)
		sysLoginLogGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:login:log:export"}}),
			middleware.OperationLog("系统登录信息", middleware.OperationTypeExport),
			sysLoginLog.Export,
		)
	}

	// 操作日志记录信息
	sysOperationLogGroup := router.Group("/api/v1/system/logs/operation")
	sysOperationLog := controller.NewSysOperationLog
	{
		sysOperationLogGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:operation:log:list"}}),
			sysOperationLog.List,
		)
		sysOperationLogGroup.DELETE("/clean",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:operation:log:remove"}}),
			middleware.OperationLog("操作日志", middleware.OperationTypeClean),
			sysOperationLog.Clean,
		)
		sysOperationLogGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:operation:log:export"}}),
			middleware.OperationLog("操作日志", middleware.OperationTypeExport),
			sysOperationLog.Export,
		)
	}

	// 菜单信息
	sysMenuGroup := router.Group("/api/v1/system/menus")
	sysMenu := controller.NewSysMenu
	{
		sysMenuGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:menu:list"}}),
			sysMenu.List,
		)
		sysMenuGroup.GET("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:menu:query"}}),
			sysMenu.Info,
		)
		sysMenuGroup.POST("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:menu:add"}}),
			middleware.OperationLog("菜单信息", middleware.OperationTypeInsert),
			sysMenu.Add,
		)
		sysMenuGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:menu:edit"}}),
			middleware.OperationLog("菜单信息", middleware.OperationTypeUpdate),
			sysMenu.Edit,
		)
		sysMenuGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:menu:remove"}}),
			middleware.OperationLog("菜单信息", middleware.OperationTypeDelete),
			sysMenu.Remove,
		)
		sysMenuGroup.GET("/tree",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:menu:list", "system:role:query"}}),
			sysMenu.Tree,
		)
		sysMenuGroup.GET("/tree/role/:role_id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:menu:list", "system:role:query"}}),
			sysMenu.TreeRole,
		)
	}

	// 通知公告信息
	sysNoticeGroup := router.Group("/api/v1/system/notices")
	sysNotice := controller.NewSysNotice
	{
		sysNoticeGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:notice:list"}}),
			sysNotice.List,
		)
		sysNoticeGroup.GET("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:notice:query"}}),
			sysNotice.Info,
		)
		sysNoticeGroup.POST("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:notice:add"}}),
			middleware.OperationLog("通知公告信息", middleware.OperationTypeInsert),
			sysNotice.Add,
		)
		sysNoticeGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:notice:edit"}}),
			middleware.OperationLog("通知公告信息", middleware.OperationTypeUpdate),
			sysNotice.Edit,
		)
		sysNoticeGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:notice:remove"}}),
			middleware.OperationLog("通知公告信息", middleware.OperationTypeDelete),
			sysNotice.Remove,
		)
	}

	// 岗位信息
	sysPostGroup := router.Group("/api/v1/system/posts")
	sysPost := controller.NewSysPost
	{
		sysPostGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:post:list"}}),
			sysPost.List,
		)
		sysPostGroup.GET("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:post:query"}}),
			sysPost.Info,
		)
		sysPostGroup.POST("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:post:add"}}),
			middleware.OperationLog("岗位信息", middleware.OperationTypeInsert),
			sysPost.Add,
		)
		sysPostGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:post:edit"}}),
			middleware.OperationLog("岗位信息", middleware.OperationTypeUpdate),
			sysPost.Edit,
		)
		sysPostGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:post:remove"}}),
			middleware.OperationLog("岗位信息", middleware.OperationTypeDelete),
			sysPost.Remove,
		)
		sysPostGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:post:export"}}),
			middleware.OperationLog("岗位信息", middleware.OperationTypeExport),
			sysPost.Export,
		)
	}

	// 个人信息
	sysProfileGroup := router.Group("/api/v1/system/users/profile")
	sysProfile := controller.NewSysProfile
	{
		sysProfileGroup.GET("",
			middleware.AuthorizeLogin(),
			sysProfile.Info,
		)
		sysProfileGroup.PUT("",
			middleware.AuthorizeLogin(),
			middleware.OperationLog("个人信息", middleware.OperationTypeUpdate),
			sysProfile.UpdateProfile,
		)
		sysProfileGroup.PUT("/password",
			middleware.AuthorizeLogin(),
			middleware.OperationLog("个人信息", middleware.OperationTypeUpdate),
			sysProfile.UpdatePassword,
		)
	}

	// 角色信息
	sysRoleGroup := router.Group("/api/v1/system/roles")
	sysRole := controller.NewSysRole
	{
		sysRoleGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:role:list"}}),
			sysRole.List,
		)
		sysRoleGroup.GET("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:role:query"}}),
			sysRole.Info,
		)
		sysRoleGroup.POST("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:role:add"}}),
			middleware.OperationLog("角色信息", middleware.OperationTypeInsert),
			sysRole.Add,
		)
		sysRoleGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:role:edit"}}),
			middleware.OperationLog("角色信息", middleware.OperationTypeUpdate),
			sysRole.Edit,
		)
		sysRoleGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:role:remove"}}),
			middleware.OperationLog("角色信息", middleware.OperationTypeDelete),
			sysRole.Remove,
		)
		sysRoleGroup.PUT("/status",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:role:edit"}}),
			middleware.OperationLog("角色信息", middleware.OperationTypeUpdate),
			sysRole.Status,
		)
		sysRoleGroup.PUT("/data-scope",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:role:edit"}}),
			middleware.OperationLog("角色信息", middleware.OperationTypeUpdate),
			sysRole.DataScope,
		)
		sysRoleGroup.GET("/user/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:role:auth"}}),
			sysRole.UserAuthList,
		)
		sysRoleGroup.PUT("/user/auth",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:role:auth"}}),
			middleware.OperationLog("角色信息", middleware.OperationTypeGrant),
			sysRole.UserAuthChecked,
		)
		sysRoleGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:role:export"}}),
			middleware.OperationLog("角色信息", middleware.OperationTypeExport),
			sysRole.Export,
		)
	}

	// 用户信息
	sysUserGroup := router.Group("/api/v1/system/users")
	sysUser := controller.NewSysUser
	{
		sysUserGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:list"}}),
			sysUser.List,
		)
		sysUserGroup.GET("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:query"}}),
			sysUser.Info,
		)
		sysUserGroup.POST("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:add"}}),
			middleware.OperationLog("用户信息", middleware.OperationTypeInsert),
			sysUser.Add,
		)
		sysUserGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:edit"}}),
			middleware.OperationLog("用户信息", middleware.OperationTypeUpdate),
			sysUser.Edit,
		)
		sysUserGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:remove"}}),
			middleware.OperationLog("用户信息", middleware.OperationTypeDelete),
			sysUser.Remove,
		)
		sysUserGroup.PUT("/password",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:passwd"}}),
			middleware.OperationLog("用户信息", middleware.OperationTypeUpdate),
			sysUser.Password,
		)
		sysUserGroup.PUT("/status",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:edit"}}),
			middleware.OperationLog("用户信息", middleware.OperationTypeUpdate),
			sysUser.Status,
		)
		sysUserGroup.PUT("/unlock/:user_name",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:unlock"}}),
			middleware.OperationLog("用户信息", middleware.OperationTypeClean),
			sysUser.Unlock,
		)
		sysUserGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:export"}}),
			middleware.OperationLog("用户信息", middleware.OperationTypeExport),
			sysUser.Export,
		)
		sysUserGroup.GET("/import/template",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:import"}}),
			sysUser.Template,
		)
		sysUserGroup.POST("/import",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:import"}}),
			middleware.OperationLog("用户信息", middleware.OperationTypeImport),
			sysUser.Import,
		)
	}
}

// initLoad 初始加载
func initLoad() {
	// 启动时，刷新缓存-参数配置
	service.NewSysConfig.CacheClean("*")
	service.NewSysConfig.CacheLoad("*")
	// 启动时，刷新缓存-字典类型数据
	service.NewSysDictType.CacheClean("*")
	service.NewSysDictType.CacheLoad("*")
}

// RegisterOnOperationLog 注入操作日志持久化钩子。
//
// 将中间件发出的操作事件映射为操作日志并落库；中间件本身不依赖任何业务模型。
// 需在路由注册（触发操作日志）之前调用。
func RegisterOnOperationLog() {
	middleware.OnOperationLog = func(rec middleware.OperationLogRecord) {
		service.NewSysOperationLog.Insert(model.SysOperationLog{
			Title:             rec.Title,
			OperationType:     rec.OperationType,
			OperationFunction: rec.OperationFunction,
			OperationUrl:      rec.OperationUrl,
			OperationMethod:   rec.OperationMethod,
			OperationIp:       rec.OperationIp,
			OperationLocation: rec.OperationLocation,
			CreateBy:          rec.CreateBy,
			OperationParam:    rec.OperationParam,
			OperationStatus:   rec.OperationStatus,
			OperationMsg:      rec.OperationMsg,
			OperationCostTime: rec.OperationCostTime,
		})
	}
}
