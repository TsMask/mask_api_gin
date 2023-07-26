package service

import "mask_api_gin/src/modules/system/model"

// ISysRole 角色 服务层接口
type ISysRole interface {
	// SelectRolePage 根据条件分页查询角色数据
	SelectRolePage(query map[string]string, dataScopeSQL string) map[string]interface{}

	// SelectRoleList 根据条件查询角色数据
	SelectRoleList(sysRole model.SysRole, dataScopeSQL string) []model.SysRole

	// SelectRoleListByUserId 根据用户ID获取角色选择框列表
	SelectRoleListByUserId(userId string) []model.SysRole

	// SelectRoleById 通过角色ID查询角色
	SelectRoleById(roleId string) model.SysRole

	// UpdateRole 修改角色信息
	UpdateRole(sysRole model.SysRole) int64

	// InsertRole 新增角色信息
	InsertRole(sysRole model.SysRole) string

	// DeleteRoleByIds 批量删除角色信息
	DeleteRoleByIds(roleIds []string) (int64, error)

	// CheckUniqueRoleName 校验角色名称是否唯一
	CheckUniqueRoleName(roleName, roleId string) bool

	// CheckUniqueRoleKey 校验角色权限是否唯一
	CheckUniqueRoleKey(roleKey, roleId string) bool

	// AuthDataScope 修改数据权限信息
	AuthDataScope(sysRole model.SysRole) int64

	// DeleteAuthUsers 批量取消授权用户角色
	DeleteAuthUsers(roleId string, userIds []string) int64

	// InsertAuthUsers 批量新增授权用户角色
	InsertAuthUsers(roleId string, userIds []string) int64
}
