package repository

import "mask_api_gin/src/modules/system/model"

// ISysMenu 菜单表 数据层接口
type ISysMenu interface {
	// SelectMenuList 查询系统菜单列表
	SelectMenuList(sysMenu model.SysMenu, userId string) []model.SysMenu

	// SelectMenuPermsByUserId 根据用户ID查询权限
	SelectMenuPermsByUserId(userId string) []string

	// SelectMenuTreeByUserId 根据用户ID查询菜单
	SelectMenuTreeByUserId(userId string) []model.SysMenu

	// SelectMenuListByRoleId 根据角色ID查询菜单树信息
	SelectMenuListByRoleId(roleId string, menuCheckStrictly bool) []string

	// SelectMenuByIds 根据菜单ID查询信息
	SelectMenuByIds(menuIds []string) []model.SysMenu

	// HasChildByMenuIdAndStatus 存在菜单子节点数量与状态
	HasChildByMenuIdAndStatus(menuId, status string) int64

	// InsertMenu 新增菜单信息
	InsertMenu(sysMenu model.SysMenu) string

	// UpdateMenu 修改菜单信息
	UpdateMenu(sysMenu model.SysMenu) int64

	// DeleteMenuById 删除菜单管理信息
	DeleteMenuById(menuId string) int64

	// CheckUniqueMenu 校验菜单是否唯一
	CheckUniqueMenu(sysMenu model.SysMenu) string
}
