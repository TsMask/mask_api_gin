package repository

import "mask_api_gin/src/modules/system/model"

// SysUserImpl 用户表 数据层处理
var SysUserImpl = &sysUserImpl{
	selectSql: "",
}

type sysUserImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// SelectUserPage 根据条件分页查询用户列表
func (r *sysUserImpl) SelectUserPage(query map[string]string, dataScopeSQL string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectUserList 根据条件查询用户列表
func (r *sysUserImpl) SelectUserList(sysUser model.SysUser, dataScopeSQL string) []model.SysUser {
	return []model.SysUser{}
}

// SelectAllocatedPage 根据条件分页查询分配用户角色列表
func (r *sysUserImpl) SelectAllocatedPage(query map[string]string, dataScopeSQL string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectUserByUserName 通过用户名查询用户
func (r *sysUserImpl) SelectUserByUserName(userName string) model.SysUser {
	return model.SysUser{}
}

// SelectUserById 通过用户ID查询用户
func (r *sysUserImpl) SelectUserById(userId string) model.SysUser {
	return model.SysUser{}
}

// InsertUser 新增用户信息
func (r *sysUserImpl) InsertUser(sysUser model.SysUser) string {
	return ""
}

// UpdateUser 修改用户信息
func (r *sysUserImpl) UpdateUser(sysUser model.SysUser) int {
	return 0
}

// DeleteUserByIds 批量删除用户信息
func (r *sysUserImpl) DeleteUserByIds(userIds []string) int {
	return 0
}

// CheckUniqueUserName 校验用户名称是否唯一
func (r *sysUserImpl) CheckUniqueUserName(userName string) string {
	return ""
}

// CheckUniquePhone 校验手机号码是否唯一
func (r *sysUserImpl) CheckUniquePhone(phonenumber string) string {
	return ""
}

// CheckUniqueEmail 校验email是否唯一
func (r *sysUserImpl) CheckUniqueEmail(email string) string {
	return ""
}
