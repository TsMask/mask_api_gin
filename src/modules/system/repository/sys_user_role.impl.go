package repository

import "mask_api_gin/src/modules/system/model"

// 用户与角色关联表 数据层处理
var SysUserRoleImpl = new(sysUserRoleImpl)

type sysUserRoleImpl struct{}

// CountUserRoleByRoleId 通过角色ID查询角色使用数量
func (r *sysUserRoleImpl) CountUserRoleByRoleId(roleId string) int {
	// 实现具体逻辑
	return 0
}

// BatchUserRole 批量新增用户角色信息
func (r *sysUserRoleImpl) BatchUserRole(sysUserRoles []model.SysUserRole) int {
	// 实现具体逻辑
	return 0
}

// DeleteUserRole 批量删除用户和角色关联
func (r *sysUserRoleImpl) DeleteUserRole(userIds []string) int {
	// 实现具体逻辑
	return 0
}

// DeleteUserRoleInfos 批量取消授权用户角色
func (r *sysUserRoleImpl) DeleteUserRoleInfos(roleId string, userIds []string) int {
	// 实现具体逻辑
	return 0
}
