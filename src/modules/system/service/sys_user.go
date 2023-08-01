package service

import "mask_api_gin/src/modules/system/model"

// ISysUser 用户 服务层接口
type ISysUser interface {
	// SelectUserPage 根据条件分页查询用户列表
	SelectUserPage(query map[string]any, dataScopeSQL string) map[string]any

	// SelectUserList 根据条件查询用户列表
	SelectUserList(sysUser model.SysUser, dataScopeSQL string) []model.SysUser

	// SelectAllocatedPage 根据条件分页查询分配用户角色列表
	SelectAllocatedPage(query map[string]any, dataScopeSQL string) map[string]any

	// SelectUserByUserName 通过用户名查询用户
	SelectUserByUserName(userName string) model.SysUser

	// SelectUserById 通过用户ID查询用户
	SelectUserById(userId string) model.SysUser

	// InsertUser 新增用户信息
	InsertUser(sysUser model.SysUser) string

	// UpdateUser 修改用户信息
	UpdateUser(sysUser model.SysUser) int64

	// UpdateUserAndRolePost 修改用户信息同时更新角色和岗位
	UpdateUserAndRolePost(sysUser model.SysUser) int64

	// DeleteUserByIds 批量删除用户信息
	DeleteUserByIds(userIds []string) (int64, error)

	// CheckUniqueUserName 校验用户名称是否唯一
	CheckUniqueUserName(userName, userId string) bool

	// CheckUniquePhone 校验手机号码是否唯一
	CheckUniquePhone(phonenumber, userId string) bool

	// CheckUniqueEmail 校验email是否唯一
	CheckUniqueEmail(email, userId string) bool

	// ImportUser 导入用户数据
	ImportUser(rows []map[string]string, isUpdateSupport bool, operName string) string
}
