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
	sysConfigGroup := router.Group("/api/v1/system/config")
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
			middleware.OperateLog("参数配置信息", middleware.OperateTypeInsert),
			sysConfig.Add,
		)
		sysConfigGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:config:edit"}}),
			middleware.OperateLog("参数配置信息", middleware.OperateTypeUpdate),
			sysConfig.Edit,
		)
		sysConfigGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:config:remove"}}),
			middleware.OperateLog("参数配置信息", middleware.OperateTypeDelete),
			sysConfig.Remove,
		)
		sysConfigGroup.PUT("/refresh",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:config:remove"}}),
			middleware.OperateLog("参数配置信息", middleware.OperateTypeOther),
			sysConfig.Refresh,
		)
		sysConfigGroup.GET("/config-key/:config_key",
			middleware.RateLimit(120, 15, middleware.RateLimitIP),
			sysConfig.ConfigKey,
		)
		sysConfigGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:config:export"}}),
			middleware.OperateLog("参数配置信息", middleware.OperateTypeExport),
			sysConfig.Export,
		)
	}

	// 部门信息
	sysDeptGroup := router.Group("/api/v1/system/dept")
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
			middleware.OperateLog("部门信息", middleware.OperateTypeInsert),
			sysDept.Add,
		)
		sysDeptGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dept:edit"}}),
			middleware.OperateLog("部门信息", middleware.OperateTypeUpdate),
			sysDept.Edit,
		)
		sysDeptGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dept:remove"}}),
			middleware.OperateLog("部门信息", middleware.OperateTypeDelete),
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
	sysDictDataGroup := router.Group("/api/v1/system/dict/data")
	sysDictData := controller.NewSysDictData
	{
		sysDictDataGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:list"}}),
			sysDictData.List,
		)
		sysDictDataGroup.GET(":id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:query"}}),
			sysDictData.Info,
		)
		sysDictDataGroup.POST("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:add"}}),
			middleware.OperateLog("字典数据信息", middleware.OperateTypeInsert),
			sysDictData.Add,
		)
		sysDictDataGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:edit"}}),
			middleware.OperateLog("字典数据信息", middleware.OperateTypeUpdate),
			sysDictData.Edit,
		)
		sysDictDataGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:remove"}}),
			middleware.OperateLog("字典数据信息", middleware.OperateTypeDelete),
			sysDictData.Remove,
		)
		sysDictDataGroup.GET("/type/:dict_type",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:query"}}),
			sysDictData.DictType,
		)
		sysDictDataGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:export"}}),
			middleware.OperateLog("字典类型信息", middleware.OperateTypeExport),
			sysDictData.Export,
		)
	}

	// 字典类型信息
	sysDictTypeGroup := router.Group("/api/v1/system/dict/type")
	sysDictType := controller.NewSysDictType
	{
		sysDictTypeGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:list"}}),
			sysDictType.List,
		)
		sysDictTypeGroup.GET("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:query"}}),
			sysDictType.Info,
		)
		sysDictTypeGroup.POST("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:add"}}),
			middleware.OperateLog("字典类型信息", middleware.OperateTypeInsert),
			sysDictType.Add,
		)
		sysDictTypeGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:edit"}}),
			middleware.OperateLog("字典类型信息", middleware.OperateTypeUpdate),
			sysDictType.Edit,
		)
		sysDictTypeGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:remove"}}),
			middleware.OperateLog("字典类型信息", middleware.OperateTypeDelete),
			sysDictType.Remove,
		)
		sysDictTypeGroup.PUT("/refresh",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:remove"}}),
			middleware.OperateLog("字典类型信息", middleware.OperateTypeOther),
			sysDictType.Refresh,
		)
		sysDictTypeGroup.GET("/options",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:query"}}),
			sysDictType.Options,
		)
		sysDictTypeGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:dict:export"}}),
			middleware.OperateLog("字典类型信息", middleware.OperateTypeExport),
			sysDictType.Export,
		)
	}

	// 系统登录日志信息
	sysLogLoginGroup := router.Group("/api/v1/system/log/login")
	sysLogLogin := controller.NewSysLogLogin
	{
		sysLogLoginGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:log:login:list"}}),
			sysLogLogin.List,
		)
		sysLogLoginGroup.DELETE("/clean",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:log:login:remove"}}),
			middleware.OperateLog("系统登录信息", middleware.OperateTypeClean),
			sysLogLogin.Clean,
		)
		sysLogLoginGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:log:login:export"}}),
			middleware.OperateLog("系统登录信息", middleware.OperateTypeExport),
			sysLogLogin.Export,
		)
	}

	// 操作日志记录信息
	sysLogOperateGroup := router.Group("/api/v1/system/log/operate")
	sysLogOperate := controller.NewSysLogOperate
	{
		sysLogOperateGroup.GET("/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:log:operate:list"}}),
			sysLogOperate.List,
		)
		sysLogOperateGroup.DELETE("/clean",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:log:operate:remove"}}),
			middleware.OperateLog("操作日志", middleware.OperateTypeClean),
			sysLogOperate.Clean,
		)
		sysLogOperateGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:log:operate:export"}}),
			middleware.OperateLog("操作日志", middleware.OperateTypeExport),
			sysLogOperate.Export,
		)
	}

	// 菜单信息
	sysMenuGroup := router.Group("/api/v1/system/menu")
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
			middleware.OperateLog("菜单信息", middleware.OperateTypeInsert),
			sysMenu.Add,
		)
		sysMenuGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:menu:edit"}}),
			middleware.OperateLog("菜单信息", middleware.OperateTypeUpdate),
			sysMenu.Edit,
		)
		sysMenuGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:menu:remove"}}),
			middleware.OperateLog("菜单信息", middleware.OperateTypeDelete),
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
	sysNoticeGroup := router.Group("/api/v1/system/notice")
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
			middleware.OperateLog("通知公告信息", middleware.OperateTypeInsert),
			sysNotice.Add,
		)
		sysNoticeGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:notice:edit"}}),
			middleware.OperateLog("通知公告信息", middleware.OperateTypeUpdate),
			sysNotice.Edit,
		)
		sysNoticeGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:notice:remove"}}),
			middleware.OperateLog("通知公告信息", middleware.OperateTypeDelete),
			sysNotice.Remove,
		)
	}

	// 岗位信息
	sysPostGroup := router.Group("/api/v1/system/post")
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
			middleware.OperateLog("岗位信息", middleware.OperateTypeInsert),
			sysPost.Add,
		)
		sysPostGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:post:edit"}}),
			middleware.OperateLog("岗位信息", middleware.OperateTypeUpdate),
			sysPost.Edit,
		)
		sysPostGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:post:remove"}}),
			middleware.OperateLog("岗位信息", middleware.OperateTypeDelete),
			sysPost.Remove,
		)
		sysPostGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:post:export"}}),
			middleware.OperateLog("岗位信息", middleware.OperateTypeExport),
			sysPost.Export,
		)
	}

	// 个人信息
	sysProfileGroup := router.Group("/api/v1/system/user/profile")
	sysProfile := controller.NewSysProfile
	{
		sysProfileGroup.GET("",
			middleware.AuthorizeLogin(),
			sysProfile.Info,
		)
		sysProfileGroup.PUT("",
			middleware.AuthorizeLogin(),
			middleware.OperateLog("个人信息", middleware.OperateTypeUpdate),
			sysProfile.UpdateProfile,
		)
		sysProfileGroup.PUT("/password",
			middleware.AuthorizeLogin(),
			middleware.OperateLog("个人信息", middleware.OperateTypeUpdate),
			sysProfile.UpdatePassword,
		)
	}

	// 角色信息
	sysRoleGroup := router.Group("/api/v1/system/role")
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
			middleware.OperateLog("角色信息", middleware.OperateTypeInsert),
			sysRole.Add,
		)
		sysRoleGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:role:edit"}}),
			middleware.OperateLog("角色信息", middleware.OperateTypeUpdate),
			sysRole.Edit,
		)
		sysRoleGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:role:remove"}}),
			middleware.OperateLog("角色信息", middleware.OperateTypeDelete),
			sysRole.Remove,
		)
		sysRoleGroup.PUT("/status",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:role:edit"}}),
			middleware.OperateLog("角色信息", middleware.OperateTypeUpdate),
			sysRole.Status,
		)
		sysRoleGroup.PUT("/data-scope",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:edit"}}),
			middleware.OperateLog("角色信息", middleware.OperateTypeUpdate),
			sysRole.DataScope,
		)
		sysRoleGroup.GET("/user/list",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:list"}}),
			sysRole.UserAuthList,
		)
		sysRoleGroup.PUT("/user/auth",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:edit"}}),
			middleware.OperateLog("角色信息", middleware.OperateTypeGrant),
			sysRole.UserAuthChecked,
		)
		sysRoleGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:export"}}),
			middleware.OperateLog("角色信息", middleware.OperateTypeExport),
			sysRole.Export,
		)
	}

	// 用户信息
	sysUserGroup := router.Group("/api/v1/system/user")
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
			middleware.OperateLog("用户信息", middleware.OperateTypeInsert),
			sysUser.Add,
		)
		sysUserGroup.PUT("",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:edit"}}),
			middleware.OperateLog("用户信息", middleware.OperateTypeUpdate),
			sysUser.Edit,
		)
		sysUserGroup.DELETE("/:id",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:remove"}}),
			middleware.OperateLog("用户信息", middleware.OperateTypeDelete),
			sysUser.Remove,
		)
		sysUserGroup.PUT("/password",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:passwd"}}),
			middleware.OperateLog("用户信息", middleware.OperateTypeUpdate),
			sysUser.Password,
		)
		sysUserGroup.PUT("/status",
			middleware.RepeatSubmit(5),
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:edit"}}),
			middleware.OperateLog("用户信息", middleware.OperateTypeUpdate),
			sysUser.Status,
		)
		sysUserGroup.PUT("/unlock/:user_name",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:unlock"}}),
			middleware.OperateLog("用户信息", middleware.OperateTypeClean),
			sysUser.Unlock,
		)
		sysUserGroup.GET("/export",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:export"}}),
			middleware.OperateLog("用户信息", middleware.OperateTypeExport),
			sysUser.Export,
		)
		sysUserGroup.GET("/import/template",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:import"}}),
			sysUser.Template,
		)
		sysUserGroup.POST("/import",
			middleware.AuthorizeUser(middleware.AuthOption{HasPerms: []string{"system:user:import"}}),
			middleware.OperateLog("用户信息", middleware.OperateTypeImport),
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
