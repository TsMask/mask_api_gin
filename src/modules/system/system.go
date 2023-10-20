package system

import (
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/middleware"
	"mask_api_gin/src/framework/middleware/collectlogs"
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
			collectlogs.OperateLog(collectlogs.OptionNew("参数配置信息", collectlogs.BUSINESS_TYPE_INSERT)),
			controller.NewSysConfig.Add,
		)
		sysConfigGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:edit"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("参数配置信息", collectlogs.BUSINESS_TYPE_UPDATE)),
			controller.NewSysConfig.Edit,
		)
		sysConfigGroup.DELETE("/:configIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:remove"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("参数配置信息", collectlogs.BUSINESS_TYPE_DELETE)),
			controller.NewSysConfig.Remove,
		)
		sysConfigGroup.PUT("/refreshCache",
			repeat.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:remove"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("参数配置信息", collectlogs.BUSINESS_TYPE_CLEAN)),
			controller.NewSysConfig.RefreshCache,
		)
		sysConfigGroup.GET("/configKey/:configKey", controller.NewSysConfig.ConfigKey)
		sysConfigGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:export"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("参数配置信息", collectlogs.BUSINESS_TYPE_EXPORT)),
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
			collectlogs.OperateLog(collectlogs.OptionNew("部门信息", collectlogs.BUSINESS_TYPE_INSERT)),
			controller.NewSysDept.Add,
		)
		sysDeptGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:edit"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("部门信息", collectlogs.BUSINESS_TYPE_UPDATE)),
			controller.NewSysDept.Edit,
		)
		sysDeptGroup.DELETE("/:deptId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:remove"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("部门信息", collectlogs.BUSINESS_TYPE_DELETE)),
			controller.NewSysDept.Remove,
		)
		sysDeptGroup.GET("/list/exclude/:deptId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:list"}}),
			controller.NewSysDept.ExcludeChild,
		)
		sysDeptGroup.GET("/treeSelect",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:list", "system:user:list"}}),
			controller.NewSysDept.TreeSelect,
		)
		sysDeptGroup.GET("/roleDeptTreeSelect/:roleId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:query"}}),
			controller.NewSysDept.RoleDeptTreeSelect,
		)
	}

	// 字典数据信息
	sysDictDataGroup := router.Group("/system/dict/data")
	{
		sysDictDataGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:list"}}),
			controller.NewSysDictData.List,
		)
		sysDictDataGroup.GET("/:dictCode",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:query"}}),
			controller.NewSysDictData.Info,
		)
		sysDictDataGroup.POST("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:add"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("字典数据信息", collectlogs.BUSINESS_TYPE_INSERT)),
			controller.NewSysDictData.Add,
		)
		sysDictDataGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:edit"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("字典数据信息", collectlogs.BUSINESS_TYPE_UPDATE)),
			controller.NewSysDictData.Edit,
		)
		sysDictDataGroup.DELETE("/:dictCodes",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:remove"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("字典数据信息", collectlogs.BUSINESS_TYPE_DELETE)),
			controller.NewSysDictData.Remove,
		)
		sysDictDataGroup.GET("/type/:dictType",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:query"}}),
			controller.NewSysDictData.DictType,
		)
		sysDictDataGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:export"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("字典类型信息", collectlogs.BUSINESS_TYPE_EXPORT)),
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
			collectlogs.OperateLog(collectlogs.OptionNew("字典类型信息", collectlogs.BUSINESS_TYPE_INSERT)),
			controller.NewSysDictType.Add,
		)
		sysDictTypeGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:edit"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("字典类型信息", collectlogs.BUSINESS_TYPE_UPDATE)),
			controller.NewSysDictType.Edit,
		)
		sysDictTypeGroup.DELETE("/:dictIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:remove"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("字典类型信息", collectlogs.BUSINESS_TYPE_DELETE)),
			controller.NewSysDictType.Remove,
		)
		sysDictTypeGroup.PUT("/refreshCache",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:remove"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("字典类型信息", collectlogs.BUSINESS_TYPE_CLEAN)),
			controller.NewSysDictType.RefreshCache,
		)
		sysDictTypeGroup.GET("/getDictOptionselect",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:query"}}),
			controller.NewSysDictType.DictOptionselect,
		)
		sysDictTypeGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:export"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("字典类型信息", collectlogs.BUSINESS_TYPE_EXPORT)),
			controller.NewSysDictType.Export,
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
			collectlogs.OperateLog(collectlogs.OptionNew("菜单信息", collectlogs.BUSINESS_TYPE_INSERT)),
			controller.NewSysMenu.Add,
		)
		sysMenuGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:edit"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("菜单信息", collectlogs.BUSINESS_TYPE_UPDATE)),
			controller.NewSysMenu.Edit,
		)
		sysMenuGroup.DELETE("/:menuId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:remove"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("菜单信息", collectlogs.BUSINESS_TYPE_DELETE)),
			controller.NewSysMenu.Remove,
		)
		sysMenuGroup.GET("/treeSelect",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:list"}}),
			controller.NewSysMenu.TreeSelect,
		)
		sysMenuGroup.GET("/roleMenuTreeSelect/:roleId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:list"}}),
			controller.NewSysMenu.RoleMenuTreeSelect,
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
			collectlogs.OperateLog(collectlogs.OptionNew("参数配置信息", collectlogs.BUSINESS_TYPE_INSERT)),
			controller.NewSysNotice.Add,
		)
		sysNoticeGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:edit"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("参数配置信息", collectlogs.BUSINESS_TYPE_UPDATE)),
			controller.NewSysNotice.Edit,
		)
		sysNoticeGroup.DELETE("/:noticeIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:remove"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("参数配置信息", collectlogs.BUSINESS_TYPE_DELETE)),
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
			collectlogs.OperateLog(collectlogs.OptionNew("岗位信息", collectlogs.BUSINESS_TYPE_INSERT)),
			controller.NewSysPost.Add,
		)
		sysPostGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:edit"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("岗位信息", collectlogs.BUSINESS_TYPE_UPDATE)),
			controller.NewSysPost.Edit,
		)
		sysPostGroup.DELETE("/:postIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:remove"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("岗位信息", collectlogs.BUSINESS_TYPE_DELETE)),
			controller.NewSysPost.Remove,
		)
		sysPostGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:export"}}),
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
		sysProfileGroup.PUT("/updatePwd",
			middleware.PreAuthorize(nil),
			collectlogs.OperateLog(collectlogs.OptionNew("个人信息", collectlogs.BUSINESS_TYPE_UPDATE)),
			controller.NewSysProfile.UpdatePwd,
		)
		sysProfileGroup.POST("/avatar",
			middleware.PreAuthorize(nil),
			collectlogs.OperateLog(collectlogs.OptionNew("用户头像", collectlogs.BUSINESS_TYPE_UPDATE)),
			controller.NewSysProfile.Avatar,
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
			collectlogs.OperateLog(collectlogs.OptionNew("角色信息", collectlogs.BUSINESS_TYPE_INSERT)),
			controller.NewSysRole.Add,
		)
		sysRoleGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:edit"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("角色信息", collectlogs.BUSINESS_TYPE_UPDATE)),
			controller.NewSysRole.Edit,
		)
		sysRoleGroup.DELETE("/:roleIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:remove"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("角色信息", collectlogs.BUSINESS_TYPE_DELETE)),
			controller.NewSysRole.Remove,
		)
		sysRoleGroup.PUT("/changeStatus",
			repeat.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:edit"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("角色信息", collectlogs.BUSINESS_TYPE_UPDATE)),
			controller.NewSysRole.Status,
		)
		sysRoleGroup.PUT("/dataScope",
			repeat.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("角色信息", collectlogs.BUSINESS_TYPE_UPDATE)),
			controller.NewSysRole.DataScope,
		)
		sysRoleGroup.GET("/authUser/allocatedList",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:list"}}),
			controller.NewSysRole.AuthUserAllocatedList,
		)
		sysRoleGroup.PUT("/authUser/checked",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("角色信息", collectlogs.BUSINESS_TYPE_GRANT)),
			controller.NewSysRole.AuthUserChecked,
		)
		sysRoleGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:export"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("角色信息", collectlogs.BUSINESS_TYPE_EXPORT)),
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
			collectlogs.OperateLog(collectlogs.OptionNew("用户信息", collectlogs.BUSINESS_TYPE_INSERT)),
			controller.NewSysUser.Add,
		)
		sysUserGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("用户信息", collectlogs.BUSINESS_TYPE_UPDATE)),
			controller.NewSysUser.Edit,
		)
		sysUserGroup.DELETE("/:userIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:remove"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("用户信息", collectlogs.BUSINESS_TYPE_DELETE)),
			controller.NewSysUser.Remove,
		)
		sysUserGroup.PUT("/resetPwd",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:resetPwd"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("用户信息", collectlogs.BUSINESS_TYPE_UPDATE)),
			controller.NewSysUser.ResetPwd,
		)
		sysUserGroup.PUT("/changeStatus",
			repeat.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("用户信息", collectlogs.BUSINESS_TYPE_UPDATE)),
			controller.NewSysUser.Status,
		)
		sysUserGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:export"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("用户信息", collectlogs.BUSINESS_TYPE_EXPORT)),
			controller.NewSysUser.Export,
		)
		sysUserGroup.GET("/importTemplate",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:import"}}),
			controller.NewSysUser.Template,
		)
		sysUserGroup.POST("/importData",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:import"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("用户信息", collectlogs.BUSINESS_TYPE_INSERT)),
			controller.NewSysUser.ImportData,
		)
	}

	// 操作日志记录信息
	sysLogOperateGroup := router.Group("/system/log/operate")
	{
		sysLogOperateGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:operate:list"}}),
			controller.NewSysLogOperate.List,
		)
		sysLogOperateGroup.DELETE("/:operIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:operate:remove"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("操作日志", collectlogs.BUSINESS_TYPE_DELETE)),
			controller.NewSysLogOperate.Remove,
		)
		sysLogOperateGroup.DELETE("/clean",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:operate:remove"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("操作日志", collectlogs.BUSINESS_TYPE_CLEAN)),
			controller.NewSysLogOperate.Clean,
		)
		sysLogOperateGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:operate:export"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("操作日志", collectlogs.BUSINESS_TYPE_EXPORT)),
			controller.NewSysLogOperate.Export,
		)
	}

	// 系统登录日志信息
	sysLogLoginGroup := router.Group("/system/log/login")
	{
		sysLogLoginGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:login:list"}}),
			controller.NewSysLogLogin.List,
		)
		sysLogLoginGroup.DELETE("/:loginIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:login:remove"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("系统登录信息", collectlogs.BUSINESS_TYPE_DELETE)),
			controller.NewSysLogLogin.Remove,
		)
		sysLogLoginGroup.DELETE("/clean",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:login:remove"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("系统登录信息", collectlogs.BUSINESS_TYPE_CLEAN)),
			controller.NewSysLogLogin.Clean,
		)
		sysLogLoginGroup.PUT("/unlock/:userName",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:login:unlock"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("系统登录信息", collectlogs.BUSINESS_TYPE_CLEAN)),
			controller.NewSysLogLogin.Unlock,
		)
		sysLogLoginGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:login:export"}}),
			collectlogs.OperateLog(collectlogs.OptionNew("系统登录信息", collectlogs.BUSINESS_TYPE_EXPORT)),
			controller.NewSysLogLogin.Export,
		)
	}
}

// InitLoad 初始参数
func InitLoad() {
	// 启动时，刷新缓存-参数配置
	service.NewSysConfigImpl.ResetConfigCache()
	// 启动时，刷新缓存-字典类型数据
	service.NewSysDictTypeImpl.ResetDictCache()
}
