package repository

import "mask_api_gin/src/modules/system/model"

// ISysUserRoleRepository 用户与角色关联表 数据层接口
type ISysUserRoleRepository interface {
	// ExistUserByRoleId 存在用户使用数量
	ExistUserByRoleId(roleId string) int64

	// DeleteByUserIds 批量删除关联By用户
	DeleteByUserIds(userIds []string) int64

	// DeleteByRoleId 批量删除关联By角色
	DeleteByRoleId(roleId string, userIds []string) int64

	// BatchInsert 批量新增信息
	BatchInsert(arr []model.SysUserRole) int64
}
