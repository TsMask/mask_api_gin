package repository

import "mask_api_gin/src/modules/system/model"

// ISysRoleMenuRepository 角色与菜单关联表 数据层接口
type ISysRoleMenuRepository interface {
	// ExistRoleByMenuId 存在角色使用数量
	ExistRoleByMenuId(menuId string) int64

	// DeleteByRoleIds 批量删除关联By角色
	DeleteByRoleIds(roleIds []string) int64

	// DeleteByMenuIds 批量删除关联By菜单
	DeleteByMenuIds(menuIds []string) int64

	// BatchInsert 批量新增信息
	BatchInsert(arr []model.SysRoleMenu) int64
}
