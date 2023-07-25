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

	// 参数配置信息
	sysConfigGroup := router.Group("/system/config")
	{
		sysConfigGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:list"}}),
			controller.SysConfig.List,
		)
		sysConfigGroup.GET("/:configId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:query"}}),
			controller.SysConfig.Info,
		)
		sysConfigGroup.POST("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:add"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_INSERT)),
			controller.SysConfig.Add,
		)
		sysConfigGroup.PUT("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:edit"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysConfig.Edit,
		)
		sysConfigGroup.DELETE("/:configIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:remove"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.SysConfig.Remove,
		)
		sysConfigGroup.PUT("/refreshCache",
			repeat.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:remove"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_CLEAN)),
			controller.SysConfig.RefreshCache,
		)
		sysConfigGroup.GET("/configKey/:configKey", controller.SysConfig.ConfigKey)
		sysConfigGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:config:export"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_EXPORT)),
			controller.SysConfig.Export,
		)
	}

	// 部门信息
	sysDeptGroup := router.Group("/system/dept")
	{
		sysDeptGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:list"}}),
			controller.SysDept.List,
		)
		sysDeptGroup.GET("/:deptId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:query"}}),
			controller.SysDept.Info,
		)
		sysDeptGroup.POST("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:add"}}),
			operlog.OperLog(operlog.OptionNew("部门信息", operlog.BUSINESS_TYPE_INSERT)),
			controller.SysDept.Add,
		)
		sysDeptGroup.PUT("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:edit"}}),
			operlog.OperLog(operlog.OptionNew("部门信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysDept.Edit,
		)
		sysDeptGroup.DELETE("/:deptId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:remove"}}),
			operlog.OperLog(operlog.OptionNew("部门信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.SysDept.Remove,
		)
		sysDeptGroup.GET("/list/exclude/:deptId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:list"}}),
			controller.SysDept.ExcludeChild,
		)
		sysDeptGroup.GET("/treeSelect",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:list", "system:user:list"}}),
			controller.SysDept.TreeSelect,
		)
		sysDeptGroup.GET("/roleDeptTreeSelect/:roleId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dept:query"}}),
			controller.SysDept.RoleDeptTreeSelect,
		)
	}

	// 字典数据信息
	sysDictDataGroup := router.Group("/system/dict/data")
	{
		sysDictDataGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:list"}}),
			controller.SysDictData.List,
		)
		sysDictDataGroup.GET("/:dictCode",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:query"}}),
			controller.SysDictData.Info,
		)
		sysDictDataGroup.POST("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:add"}}),
			operlog.OperLog(operlog.OptionNew("字典数据信息", operlog.BUSINESS_TYPE_INSERT)),
			controller.SysDictData.Add,
		)
		sysDictDataGroup.PUT("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:edit"}}),
			operlog.OperLog(operlog.OptionNew("字典数据信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysDictData.Edit,
		)
		sysDictDataGroup.DELETE("/:dictCodes",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:remove"}}),
			operlog.OperLog(operlog.OptionNew("字典数据信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.SysDictData.Remove,
		)
		sysDictDataGroup.GET("/type/:dictType",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:query"}}),
			controller.SysDictData.DictType,
		)
		sysDictDataGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:export"}}),
			operlog.OperLog(operlog.OptionNew("字典类型信息", operlog.BUSINESS_TYPE_EXPORT)),
			controller.SysDictData.Export,
		)
	}

	// 字典类型信息
	sysDictTypeGroup := router.Group("/system/dict/type")
	{
		sysDictTypeGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:list"}}),
			controller.SysDictType.List,
		)
		sysDictTypeGroup.GET("/:dictId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:query"}}),
			controller.SysDictType.Info,
		)
		sysDictTypeGroup.POST("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:add"}}),
			operlog.OperLog(operlog.OptionNew("字典类型信息", operlog.BUSINESS_TYPE_INSERT)),
			controller.SysDictType.Add,
		)
		sysDictTypeGroup.PUT("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:edit"}}),
			operlog.OperLog(operlog.OptionNew("字典类型信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysDictType.Edit,
		)
		sysDictTypeGroup.DELETE("/:dictIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:remove"}}),
			operlog.OperLog(operlog.OptionNew("字典类型信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.SysDictType.Remove,
		)
		sysDictTypeGroup.PUT("/refreshCache",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:remove"}}),
			operlog.OperLog(operlog.OptionNew("字典类型信息", operlog.BUSINESS_TYPE_CLEAN)),
			controller.SysDictType.RefreshCache,
		)
		sysDictTypeGroup.GET("/getDictOptionselect",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:query"}}),
			controller.SysDictType.DictOptionselect,
		)
		sysDictTypeGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:dict:export"}}),
			operlog.OperLog(operlog.OptionNew("字典类型信息", operlog.BUSINESS_TYPE_EXPORT)),
			controller.SysDictType.Export,
		)
	}

	// 菜单信息
	sysMenuGroup := router.Group("/system/menu")
	{
		sysMenuGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:list"}}),
			controller.SysMenu.List,
		)
		sysMenuGroup.GET("/:menuId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:query"}}),
			controller.SysMenu.Info,
		)
		sysMenuGroup.POST("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:add"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_INSERT)),
			controller.SysMenu.Add,
		)
		sysMenuGroup.PUT("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:edit"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysMenu.Edit,
		)
		sysMenuGroup.DELETE("/:menuId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:remove"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.SysMenu.Remove,
		)
		sysMenuGroup.GET("/treeSelect",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:list"}}),
			controller.SysMenu.TreeSelect,
		)
		sysMenuGroup.GET("/roleMenuTreeSelect/:roleId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:menu:list"}}),
			controller.SysMenu.RoleMenuTreeSelect,
		)
	}

	// 通知公告信息
	sysNoticeGroup := router.Group("/system/notice")
	{
		sysNoticeGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:list"}}),
			controller.SysNotice.List,
		)
		sysNoticeGroup.GET("/:noticeId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:query"}}),
			controller.SysNotice.Info,
		)
		sysNoticeGroup.POST("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:add"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_INSERT)),
			controller.SysNotice.Add,
		)
		sysNoticeGroup.PUT("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:edit"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysNotice.Edit,
		)
		sysNoticeGroup.DELETE("/:noticeIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:notice:remove"}}),
			operlog.OperLog(operlog.OptionNew("参数配置信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.SysNotice.Remove,
		)
	}

	// 岗位信息
	sysPostGroup := router.Group("/system/post")
	{
		sysPostGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:list"}}),
			controller.SysPost.List,
		)
		sysPostGroup.GET("/:postId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:query"}}),
			controller.SysPost.Info,
		)
		sysPostGroup.POST("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:add"}}),
			operlog.OperLog(operlog.OptionNew("岗位信息", operlog.BUSINESS_TYPE_INSERT)),
			controller.SysPost.Add,
		)
		sysPostGroup.PUT("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:edit"}}),
			operlog.OperLog(operlog.OptionNew("岗位信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysPost.Edit,
		)
		sysPostGroup.DELETE("/:postIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:remove"}}),
			operlog.OperLog(operlog.OptionNew("岗位信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.SysPost.Remove,
		)
		sysPostGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:post:export"}}),
			controller.SysPost.Export,
		)
	}

	// 个人信息
	sysProfileGroup := router.Group("/system/user/profile")
	{
		sysProfileGroup.GET("/",
			middleware.PreAuthorize(nil),
			controller.SysProfile.Info,
		)
		sysProfileGroup.PUT("/",
			middleware.PreAuthorize(nil),
			controller.SysProfile.UpdateProfile,
		)
		sysProfileGroup.PUT("/updatePwd",
			middleware.PreAuthorize(nil),
			operlog.OperLog(operlog.OptionNew("个人信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysProfile.UpdatePwd,
		)
		sysProfileGroup.POST("/avatar",
			middleware.PreAuthorize(nil),
			operlog.OperLog(operlog.OptionNew("用户头像", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysProfile.Avatar,
		)
	}

	// 角色信息
	sysRoleGroup := router.Group("/system/role")
	{
		sysRoleGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:list"}}),
			controller.SysRole.List,
		)
		sysRoleGroup.GET("/:roleId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:query"}}),
			controller.SysRole.Info,
		)
		sysRoleGroup.POST("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:add"}}),
			operlog.OperLog(operlog.OptionNew("角色信息", operlog.BUSINESS_TYPE_INSERT)),
			controller.SysRole.Add,
		)
		sysRoleGroup.PUT("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:edit"}}),
			operlog.OperLog(operlog.OptionNew("角色信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysRole.Edit,
		)
		sysRoleGroup.DELETE("/:roleIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:remove"}}),
			operlog.OperLog(operlog.OptionNew("角色信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.SysRole.Remove,
		)
		sysRoleGroup.PUT("/changeStatus",
			repeat.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:role:edit"}}),
			operlog.OperLog(operlog.OptionNew("角色信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysRole.Status,
		)
		sysRoleGroup.PUT("/dataScope",
			repeat.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			operlog.OperLog(operlog.OptionNew("角色信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysRole.DataScope,
		)
		sysRoleGroup.GET("/authUser/allocatedList",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:list"}}),
			controller.SysRole.AuthUserAllocatedList,
		)
		sysRoleGroup.PUT("/authUser/checked",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			operlog.OperLog(operlog.OptionNew("角色信息", operlog.BUSINESS_TYPE_GRANT)),
			controller.SysRole.AuthUserChecked,
		)
		sysRoleGroup.POST("/export",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:export"}}),
			operlog.OperLog(operlog.OptionNew("角色信息", operlog.BUSINESS_TYPE_EXPORT)),
			controller.SysRole.Export,
		)
	}

	// 用户信息
	sysUserGroup := router.Group("/system/user")
	{
		sysUserGroup.GET("/list",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:list"}}),
			controller.SysUser.List,
		)
		sysUserGroup.GET("/:userId",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:query"}}),
			controller.SysUser.Info,
		)
		sysUserGroup.POST("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:add"}}),
			operlog.OperLog(operlog.OptionNew("用户信息", operlog.BUSINESS_TYPE_INSERT)),
			controller.SysUser.Add,
		)
		sysUserGroup.PUT("/",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			operlog.OperLog(operlog.OptionNew("用户信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysUser.Edit,
		)
		sysUserGroup.DELETE("/:userIds",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:remove"}}),
			operlog.OperLog(operlog.OptionNew("用户信息", operlog.BUSINESS_TYPE_DELETE)),
			controller.SysUser.Remove,
		)
		sysUserGroup.PUT("/resetPwd",
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:resetPwd"}}),
			operlog.OperLog(operlog.OptionNew("用户信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysUser.ResetPwd,
		)
		sysUserGroup.PUT("/changeStatus",
			repeat.RepeatSubmit(5),
			middleware.PreAuthorize(map[string][]string{"hasPerms": {"system:user:edit"}}),
			operlog.OperLog(operlog.OptionNew("用户信息", operlog.BUSINESS_TYPE_UPDATE)),
			controller.SysUser.Status,
		)
		// 用户信息列表导入模板下载 TODO
		// 用户信息列表导入 TODO
		// 用户信息列表导出 TODO
	}
}

// InitLoad 初始参数
func InitLoad() {
	// 启动时，刷新缓存-参数配置
	service.SysConfigImpl.ResetConfigCache()
}
