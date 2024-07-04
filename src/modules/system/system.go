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
			middleware.OperateLog(middleware.OptionNew("参数配置信息", middleware.BusinessTypeInsert)),
			controller.NewSysConfig.Add,
		)
		sysConfigGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:edit"}}),
			middleware.OperateLog(middleware.OptionNew("参数配置信息", middleware.BusinessTypeUpdate)),
			controller.NewSysConfig.Edit,
		)
		sysConfigGroup.DELETE("/:configIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:remove"}}),
			middleware.OperateLog(middleware.OptionNew("参数配置信息", middleware.BusinessTypeDelete)),
			controller.NewSysConfig.Remove,
		)
		sysConfigGroup.PUT("/refreshCache",
			middleware.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:remove"}}),
			middleware.OperateLog(middleware.OptionNew("参数配置信息", middleware.BusinessTypeClean)),
			controller.NewSysConfig.RefreshCache,
		)
		sysConfigGroup.GET("/configKey/:configKey", controller.NewSysConfig.ConfigKey)
		sysConfigGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:export"}}),
			middleware.OperateLog(middleware.OptionNew("参数配置信息", middleware.BusinessTypeExport)),
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
			middleware.OperateLog(middleware.OptionNew("部门信息", middleware.BusinessTypeInsert)),
			controller.NewSysDept.Add,
		)
		sysDeptGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:edit"}}),
			middleware.OperateLog(middleware.OptionNew("部门信息", middleware.BusinessTypeUpdate)),
			controller.NewSysDept.Edit,
		)
		sysDeptGroup.DELETE("/:deptId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:remove"}}),
			middleware.OperateLog(middleware.OptionNew("部门信息", middleware.BusinessTypeDelete)),
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
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:query", "system:user:edit"}}),
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
			middleware.OperateLog(middleware.OptionNew("字典数据信息", middleware.BusinessTypeInsert)),
			controller.NewSysDictData.Add,
		)
		sysDictDataGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:edit"}}),
			middleware.OperateLog(middleware.OptionNew("字典数据信息", middleware.BusinessTypeUpdate)),
			controller.NewSysDictData.Edit,
		)
		sysDictDataGroup.DELETE("/:dictCodes",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:remove"}}),
			middleware.OperateLog(middleware.OptionNew("字典数据信息", middleware.BusinessTypeDelete)),
			controller.NewSysDictData.Remove,
		)
		sysDictDataGroup.GET("/type/:dictType",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:query"}}),
			controller.NewSysDictData.DictType,
		)
		sysDictDataGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:export"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.BusinessTypeExport)),
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
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.BusinessTypeInsert)),
			controller.NewSysDictType.Add,
		)
		sysDictTypeGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:edit"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.BusinessTypeUpdate)),
			controller.NewSysDictType.Edit,
		)
		sysDictTypeGroup.DELETE("/:dictIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:remove"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.BusinessTypeDelete)),
			controller.NewSysDictType.Remove,
		)
		sysDictTypeGroup.PUT("/refreshCache",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:remove"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.BusinessTypeClean)),
			controller.NewSysDictType.RefreshCache,
		)
		sysDictTypeGroup.GET("/getDictOptionSelect",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:query"}}),
			controller.NewSysDictType.DictOptionSelect,
		)
		sysDictTypeGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:export"}}),
			middleware.OperateLog(middleware.OptionNew("字典类型信息", middleware.BusinessTypeExport)),
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
			middleware.OperateLog(middleware.OptionNew("菜单信息", middleware.BusinessTypeInsert)),
			controller.NewSysMenu.Add,
		)
		sysMenuGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:edit"}}),
			middleware.OperateLog(middleware.OptionNew("菜单信息", middleware.BusinessTypeUpdate)),
			controller.NewSysMenu.Edit,
		)
		sysMenuGroup.DELETE("/:menuId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:remove"}}),
			middleware.OperateLog(middleware.OptionNew("菜单信息", middleware.BusinessTypeDelete)),
			controller.NewSysMenu.Remove,
		)
		sysMenuGroup.GET("/treeSelect",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:list"}}),
			controller.NewSysMenu.TreeSelect,
		)
		sysMenuGroup.GET("/roleMenuTreeSelect/:roleId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:list", "system:role:query"}}),
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
			middleware.OperateLog(middleware.OptionNew("通知公告信息", middleware.BusinessTypeInsert)),
			controller.NewSysNotice.Add,
		)
		sysNoticeGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:edit"}}),
			middleware.OperateLog(middleware.OptionNew("通知公告信息", middleware.BusinessTypeUpdate)),
			controller.NewSysNotice.Edit,
		)
		sysNoticeGroup.DELETE("/:noticeIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:remove"}}),
			middleware.OperateLog(middleware.OptionNew("通知公告信息", middleware.BusinessTypeDelete)),
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
			middleware.OperateLog(middleware.OptionNew("岗位信息", middleware.BusinessTypeInsert)),
			controller.NewSysPost.Add,
		)
		sysPostGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:edit"}}),
			middleware.OperateLog(middleware.OptionNew("岗位信息", middleware.BusinessTypeUpdate)),
			controller.NewSysPost.Edit,
		)
		sysPostGroup.DELETE("/:postIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:remove"}}),
			middleware.OperateLog(middleware.OptionNew("岗位信息", middleware.BusinessTypeDelete)),
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
			middleware.OperateLog(middleware.OptionNew("个人信息", middleware.BusinessTypeUpdate)),
			controller.NewSysProfile.UpdatePwd,
		)
		sysProfileGroup.POST("/avatar",
			middleware.PreAuthorize(nil),
			middleware.OperateLog(middleware.OptionNew("用户头像", middleware.BusinessTypeUpdate)),
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
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.BusinessTypeInsert)),
			controller.NewSysRole.Add,
		)
		sysRoleGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:edit"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.BusinessTypeUpdate)),
			controller.NewSysRole.Edit,
		)
		sysRoleGroup.DELETE("/:roleIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:remove"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.BusinessTypeDelete)),
			controller.NewSysRole.Remove,
		)
		sysRoleGroup.PUT("/changeStatus",
			middleware.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:edit"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.BusinessTypeUpdate)),
			controller.NewSysRole.Status,
		)
		sysRoleGroup.PUT("/dataScope",
			middleware.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.BusinessTypeUpdate)),
			controller.NewSysRole.DataScope,
		)
		sysRoleGroup.GET("/authUser/allocatedList",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:list"}}),
			controller.NewSysRole.AuthUserAllocatedList,
		)
		sysRoleGroup.PUT("/authUser/checked",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.BusinessTypeGrant)),
			controller.NewSysRole.AuthUserChecked,
		)
		sysRoleGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:export"}}),
			middleware.OperateLog(middleware.OptionNew("角色信息", middleware.BusinessTypeExport)),
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
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.BusinessTypeInsert)),
			controller.NewSysUser.Add,
		)
		sysUserGroup.PUT("",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.BusinessTypeUpdate)),
			controller.NewSysUser.Edit,
		)
		sysUserGroup.DELETE("/:userIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:remove"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.BusinessTypeDelete)),
			controller.NewSysUser.Remove,
		)
		sysUserGroup.PUT("/resetPwd",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:resetPwd"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.BusinessTypeUpdate)),
			controller.NewSysUser.ResetPwd,
		)
		sysUserGroup.PUT("/changeStatus",
			middleware.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.BusinessTypeUpdate)),
			controller.NewSysUser.Status,
		)
		sysUserGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:export"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.BusinessTypeExport)),
			controller.NewSysUser.Export,
		)
		sysUserGroup.GET("/importTemplate",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:import"}}),
			controller.NewSysUser.Template,
		)
		sysUserGroup.POST("/importData",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:import"}}),
			middleware.OperateLog(middleware.OptionNew("用户信息", middleware.BusinessTypeInsert)),
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
		sysLogOperateGroup.DELETE("/:operaIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:operate:remove"}}),
			middleware.OperateLog(middleware.OptionNew("操作日志", middleware.BusinessTypeDelete)),
			controller.NewSysLogOperate.Remove,
		)
		sysLogOperateGroup.DELETE("/clean",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:operate:remove"}}),
			middleware.OperateLog(middleware.OptionNew("操作日志", middleware.BusinessTypeClean)),
			controller.NewSysLogOperate.Clean,
		)
		sysLogOperateGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:operate:export"}}),
			middleware.OperateLog(middleware.OptionNew("操作日志", middleware.BusinessTypeExport)),
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
			middleware.OperateLog(middleware.OptionNew("系统登录信息", middleware.BusinessTypeDelete)),
			controller.NewSysLogLogin.Remove,
		)
		sysLogLoginGroup.DELETE("/clean",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:login:remove"}}),
			middleware.OperateLog(middleware.OptionNew("系统登录信息", middleware.BusinessTypeClean)),
			controller.NewSysLogLogin.Clean,
		)
		sysLogLoginGroup.PUT("/unlock/:userName",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:login:unlock"}}),
			middleware.OperateLog(middleware.OptionNew("系统登录信息", middleware.BusinessTypeClean)),
			controller.NewSysLogLogin.Unlock,
		)
		sysLogLoginGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:log:login:export"}}),
			middleware.OperateLog(middleware.OptionNew("系统登录信息", middleware.BusinessTypeExport)),
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
