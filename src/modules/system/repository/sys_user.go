package repository

import "mask_api_gin/src/modules/system/model"

// ISysUser 用户表 数据层接口
type ISysUser interface {
	// SelectUserPage 根据条件分页查询用户列表
	SelectUserPage(queryMap map[string]any, dataScopeSQL string) map[string]any

	// SelectAllocatedPage 根据条件分页查询分配用户角色列表
	SelectAllocatedPage(query map[string]any, dataScopeSQL string) map[string]any

	// SelectUserList 根据条件查询用户列表
	SelectUserList(sysUser model.SysUser, dataScopeSQL string) []model.SysUser

	// SelectUserByIds 通过用户ID查询用户
	SelectUserByIds(userIds []string) []model.SysUser

	// SelectUserByUserName 通过用户登录账号查询用户
	SelectUserByUserName(userName string) model.SysUser

	// InsertUser 新增用户信息
	InsertUser(sysUser model.SysUser) string

	// UpdateUser 修改用户信息
	UpdateUser(sysUser model.SysUser) int64

	// DeleteUserByIds 批量删除用户信息
	DeleteUserByIds(userIds []string) int64

	// CheckUniqueUser 校验用户信息是否唯一
	CheckUniqueUser(sysUser model.SysUser) string
}
