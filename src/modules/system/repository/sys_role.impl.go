package repository

import "mask_api_gin/src/modules/system/model"

// SysRoleImpl 角色表 数据层处理
var SysRoleImpl = &sysRoleImpl{
	selectSql: "",
}

type sysRoleImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// SelectRolePage 根据条件分页查询角色数据
func (r *sysRoleImpl) SelectRolePage(query map[string]string, dataScopeSQL string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectRoleList 根据条件查询角色数据
func (r *sysRoleImpl) SelectRoleList(sysRole model.SysRole, dataScopeSQL string) []model.SysRole {
	return []model.SysRole{}
}

// SelectRolePermissionByUserId 根据用户ID查询角色
func (r *sysRoleImpl) SelectRolePermissionByUserId(userId string) []model.SysRole {
	return []model.SysRole{}
}

// SelectRoleIdsByUserId 根据用户ID获取拥有角色ID
func (r *sysRoleImpl) SelectRoleIdsByUserId(userId string) []string {
	return []string{}
}

// SelectRoleById 通过角色ID查询角色
func (r *sysRoleImpl) SelectRoleById(roleId string) model.SysRole {
	return model.SysRole{}
}

// SelectRolesByUserName 根据用户名查询角色
func (r *sysRoleImpl) SelectRolesByUserName(userName string) []model.SysRole {
	return []model.SysRole{}
}

// CheckUniqueRoleName 校验角色名称是否唯一
func (r *sysRoleImpl) CheckUniqueRoleName(roleName string) string {
	return ""
}

// CheckUniqueRoleKey 校验角色权限是否唯一
func (r *sysRoleImpl) CheckUniqueRoleKey(roleKey string) string {
	return ""
}

// UpdateRole 修改角色信息
func (r *sysRoleImpl) UpdateRole(sysRole model.SysRole) int {
	return 0
}

// InsertRole 新增角色信息
func (r *sysRoleImpl) InsertRole(sysRole model.SysRole) string {
	return ""
}

// DeleteRoleByIds 批量删除角色信息
func (r *sysRoleImpl) DeleteRoleByIds(roleIds []string) int {
	return 0
}
