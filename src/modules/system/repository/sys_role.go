package repository

import "mask_api_gin/src/modules/system/model"

// ISysRole 角色表 数据层接口
type ISysRole interface {
	// SelectRolePage 根据条件分页查询角色数据
	SelectRolePage(query map[string]any, dataScopeSQL string) map[string]any

	// SelectRoleList 根据条件查询角色数据
	SelectRoleList(sysRole model.SysRole, dataScopeSQL string) []model.SysRole

	// SelectRoleListByUserId 根据用户ID获取角色选择框列表
	SelectRoleListByUserId(userId string) []model.SysRole

	// SelectRoleByIds 通过角色ID查询角色
	SelectRoleByIds(roleIds []string) []model.SysRole

	// UpdateRole 修改角色信息
	UpdateRole(sysRole model.SysRole) int64

	// InsertRole 新增角色信息
	InsertRole(sysRole model.SysRole) string

	// DeleteRoleByIds 批量删除角色信息
	DeleteRoleByIds(roleIds []string) int64

	// CheckUniqueRole 校验角色是否唯一
	CheckUniqueRole(sysRole model.SysRole) string
}
