package service

import (
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/system/model"
)

// ISysMenu 菜单 服务层接口
type ISysMenu interface {
	// SelectMenuList 查询系统菜单列表
	SelectMenuList(sysMenu model.SysMenu, userId string) []model.SysMenu

	// SelectMenuPermsByUserId 根据用户ID查询权限
	SelectMenuPermsByUserId(userId string) []string

	// SelectMenuPermsByUserId 根据用户ID查询权限
	SelectMenuTreeByUserId(userId string) []model.SysMenu

	// SelectMenuTreeSelectByUserId 查询菜单树结构信息
	SelectMenuTreeSelectByUserId(sysMenu model.SysMenu, userId string) []vo.TreeSelect

	// SelectMenuListByRoleId 根据角色ID查询菜单树信息
	SelectMenuListByRoleId(roleId string) []string

	// SelectMenuById 根据菜单ID查询信息
	SelectMenuById(menuId string) model.SysMenu

	// HasChildByMenuId 存在菜单子节点数量
	HasChildByMenuId(menuId string) int64

	// CheckMenuExistRole 查询菜单分配角色数量
	CheckMenuExistRole(menuId string) int64

	// InsertMenu 新增菜单信息
	InsertMenu(sysMenu model.SysMenu) string

	// UpdateMenu 修改菜单信息
	UpdateMenu(sysMenu model.SysMenu) int64

	// DeleteMenuById 删除菜单管理信息
	DeleteMenuById(menuId string) int64

	// CheckUniqueMenuName 校验菜单名称是否唯一
	CheckUniqueMenuName(menuName, parentId, menuId string) bool

	// CheckUniqueMenuPath 校验路由地址是否唯一（针对目录和菜单）
	CheckUniqueMenuPath(path, menuId string) bool

	// BuildRouteMenus 构建前端路由所需要的菜单
	BuildRouteMenus(sysMenus []model.SysMenu, prefix string) []vo.Router
}
