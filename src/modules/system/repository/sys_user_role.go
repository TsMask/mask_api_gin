package repository

import "mask_api_gin/src/modules/system/model"

// ISysUserRole 用户与角色关联表 数据层接口
type ISysUserRole interface {
	// CountUserRoleByRoleId 通过角色ID查询角色使用数量
	CountUserRoleByRoleId(roleId string) int

	// BatchUserRole 批量新增用户角色信息
	BatchUserRole(sysUserRoles []model.SysUserRole) int64

	// DeleteUserRole 批量删除用户和角色关联
	DeleteUserRole(userIds []string) int64

	// DeleteUserRoleInfos 批量取消授权用户角色
	DeleteUserRoleInfos(roleId string, userIds []string) int
}
