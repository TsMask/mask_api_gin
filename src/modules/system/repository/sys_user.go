package repository

import "mask_api_gin/src/modules/system/model"

// ISysUser 用户表 数据层接口
type ISysUser interface {
	// SelectUserPage 根据条件分页查询用户列表
	SelectUserPage(query map[string]string, dataScopeSQL string) map[string]interface{}

	// SelectAllocatedPage 根据条件分页查询分配用户角色列表
	SelectAllocatedPage(query map[string]string, dataScopeSQL string) map[string]interface{}

	// SelectUserList 根据条件查询用户列表
	SelectUserList(sysUser model.SysUser, dataScopeSQL string) []model.SysUser

	// SelectUserById 通过用户ID查询用户
	SelectUserById(userID string) model.SysUser

	// SelectUserByUserName 通过用户登录账号查询用户
	SelectUserByUserName(userName string) model.SysUser

	// InsertUser 新增用户信息
	InsertUser(sysUser model.SysUser) string

	// UpdateUser 修改用户信息
	UpdateUser(sysUser model.SysUser) int

	// DeleteUserByIds 批量删除用户信息
	DeleteUserByIds(userIds []string) int

	// CheckUniqueUserName 校验用户名称是否唯一
	CheckUniqueUserName(userName string) string

	// CheckUniquePhone 校验手机号码是否唯一
	CheckUniquePhone(phonenumber string) string

	// CheckUniqueEmail 校验email是否唯一
	CheckUniqueEmail(email string) string
}
