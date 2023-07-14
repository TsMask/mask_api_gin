package service

import (
	"mask_api_gin/src/modules/system/model"
	pkgModel "mask_api_gin/src/pkg/model"
)

// ISysMenu 菜单 服务层接口
type ISysMenu interface {
	// SelectMenuList 查询系统菜单列表
	SelectMenuList(sysMenu model.SysMenu, userId string) []model.SysMenu

	// SelectMenuPermsByRoleId 根据角色ID查询权限
	SelectMenuPermsByRoleId(roleId string) []string

	// SelectMenuPermsByUserId 根据用户ID查询权限
	SelectMenuPermsByUserId(userId string) []string

	// SelectMenuTreeByUserId 根据用户ID查询菜单
	SelectMenuTreeByUserId(userId string) []model.SysMenu

	// SelectMenuListByRoleId 根据角色ID查询菜单树信息
	SelectMenuListByRoleId(roleId string, menuCheckStrictly bool) []string

	// SelectMenuById 根据菜单ID查询信息
	SelectMenuById(menuId string) model.SysMenu

	// HasChildByMenuId 是否存在菜单子节点
	HasChildByMenuId(menuId string) int

	// CheckMenuExistRole 查询菜单是否存在角色
	CheckMenuExistRole(menuId string) int

	// InsertMenu 新增菜单信息
	InsertMenu(sysMenu model.SysMenu) string

	// UpdateMenu 修改菜单信息
	UpdateMenu(sysMenu model.SysMenu) int

	// DeleteMenuById 删除菜单管理信息
	DeleteMenuById(menuId string) int

	// CheckUniqueMenuName 校验菜单名称是否唯一
	CheckUniqueMenuName(menuName, parentId string) string

	// CheckUniqueMenuPath 校验路由地址是否唯一（针对目录和菜单）
	CheckUniqueMenuPath(path string) string

	// BuildRouteMenus 构建前端路由所需要的菜单
	BuildRouteMenus([]model.SysMenu) []pkgModel.Router
}
