package service

import "mask_api_gin/src/modules/system/model"

// SysUserRoleImpl 用户与角色关联 数据层处理
var SysUserRoleImpl = &sysUserRoleImpl{
	selectSql: "",
}

type sysUserRoleImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// CountUserRoleByRoleId 通过角色ID查询角色使用数量
func (r *sysUserPostImpl) CountUserRoleByRoleId(roleId string) int {
	return 0
}

// BatchUserRole 批量新增用户角色信息
func (r *sysUserPostImpl) BatchUserRole(sysUserRoles []model.SysUserRole) int {
	return 0
}

// DeleteUserRole 批量删除用户和角色关联
func (r *sysUserPostImpl) DeleteUserRole(userIds []string) int {
	return 0
}

// DeleteUserRoleInfos 批量取消授权用户角色
func (r *sysUserPostImpl) DeleteUserRoleInfos(roleId string, userIds []string) int {
	return 0
}
