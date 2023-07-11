package service

import "mask_api_gin/src/modules/system/model"

// ISysRole 角色 服务层接口
type ISysRole interface {
	// SelectRolePage 根据条件分页查询角色数据
	SelectRolePage(query map[string]string, dataScopeSQL string) map[string]interface{}

	// SelectRoleList 根据条件查询角色数据
	SelectRoleList(sysRole model.SysRole, dataScopeSQL string) []model.SysRole

	// SelectRolePermissionByUserId 根据用户ID查询角色
	SelectRolePermissionByUserId(userId string) []model.SysRole

	// SelectRoleIdsByUserId 根据用户ID获取拥有角色ID
	SelectRoleIdsByUserId(userId string) []string

	// SelectRoleById 通过角色ID查询角色
	SelectRoleById(roleId string) model.SysRole

	// SelectRolesByUserName 根据用户名查询角色
	SelectRolesByUserName(userName string) []model.SysRole

	// CheckUniqueRoleName 校验角色名称是否唯一
	CheckUniqueRoleName(roleName string) (string, error)

	// CheckUniqueRoleKey 校验角色权限是否唯一
	CheckUniqueRoleKey(roleKey string) (string, error)

	// UpdateRole 修改角色信息
	UpdateRole(sysRole model.SysRole) int

	// InsertRole 新增角色信息
	InsertRole(sysRole model.SysRole) string

	// DeleteRoleByIds 批量删除角色信息
	DeleteRoleByIds(roleIds []string) int
}
