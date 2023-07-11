package repository

import "mask_api_gin/src/modules/system/model"

// ISysRoleMenu 角色与菜单关联表 数据层接口
type ISysRoleMenu interface {
	// CheckMenuExistRole 查询菜单使用数量
	CheckMenuExistRole(menuId string) int

	// DeleteRoleMenuByRoleId 通过角色ID删除角色和菜单关联
	DeleteRoleMenuByRoleId(roleId string) int

	// DeleteRoleMenu 批量删除角色菜单关联信息
	DeleteRoleMenu(roleIds []string) int

	// BatchRoleMenu 批量新增角色菜单信息
	BatchRoleMenu(sysRoleMenus []model.SysRoleMenu) int
}
