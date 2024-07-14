package service

import "mask_api_gin/src/modules/system/model"

// ISysUserService 用户 服务层接口
type ISysUserService interface {
	// FindByPage 分页查询列表数据
	FindByPage(queryMap map[string]any, dataScopeSQL string) map[string]any

	// Find 查询列表数据
	Find(sysUser model.SysUser, dataScopeSQL string) []model.SysUser

	// FindById 通过ID查询信息
	FindById(userId string) model.SysUser

	// Insert 新增信息
	Insert(sysUser model.SysUser) string

	// Update 修改信息
	Update(sysUser model.SysUser) int64

	// UpdateUserAndRolePost 修改用户信息同时更新角色和岗位
	UpdateUserAndRolePost(sysUser model.SysUser) int64

	// DeleteByIds 批量删除信息
	DeleteByIds(userIds []string) (int64, error)

	// CheckUniqueByUserName 检查用户名称是否唯一
	CheckUniqueByUserName(userName, userId string) bool

	// CheckUniqueByPhone 检查手机号码是否唯一
	CheckUniqueByPhone(phone, userId string) bool

	// CheckUniqueByEmail 检查Email是否唯一
	CheckUniqueByEmail(email, userId string) bool

	// FindByUserName 通过用户名查询用户信息
	FindByUserName(userName string) model.SysUser

	// FindAllocatedPage 根据条件分页查询分配用户角色列表
	FindAllocatedPage(queryMap map[string]any, dataScopeSQL string) map[string]any
}
