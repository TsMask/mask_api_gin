package repository

import "mask_api_gin/src/modules/system/model"

// ISysRoleMenu 角色与菜单关联表 数据层接口
type ISysRoleMenu interface {
	// CheckMenuExistRole 查询菜单分配给角色使用数量
	CheckMenuExistRole(menuId string) int64

	// DeleteRoleMenu 批量删除角色和菜单关联
	DeleteRoleMenu(roleIds []string) int64

	// DeleteMenuRole 批量删除菜单和角色关联
	DeleteMenuRole(menuIds []string) int64

	// BatchRoleMenu 批量新增角色菜单信息
	BatchRoleMenu(sysRoleMenus []model.SysRoleMenu) int64
}
