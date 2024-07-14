package service

import "mask_api_gin/src/modules/system/model"

// ISysRoleService 角色 服务层接口
type ISysRoleService interface {
	// FindByPage 分页查询列表数据
	FindByPage(query map[string]any, dataScopeSQL string) map[string]any

	// Find 查询列表数据
	Find(sysRole model.SysRole, dataScopeSQL string) []model.SysRole

	// FindById 通过ID查询信息
	FindById(roleId string) model.SysRole

	// Insert 新增信息
	Insert(sysRole model.SysRole) string

	// Update 修改信息
	Update(sysRole model.SysRole) int64

	// DeleteByIds 批量删除信息
	DeleteByIds(roleIds []string) (int64, error)

	// FindByUserId 根据用户ID获取角色选择框列表
	FindByUserId(userId string) []model.SysRole

	// CheckUniqueByName 检查角色名称是否唯一
	CheckUniqueByName(roleName, roleId string) bool

	// CheckUniqueByKey 检查角色权限是否唯一
	CheckUniqueByKey(roleKey, roleId string) bool

	// UpdateAndDataScope 修改信息同时更新数据权限信息
	UpdateAndDataScope(sysRole model.SysRole) int64

	// InsertAuthUsers 批量新增授权用户角色
	InsertAuthUsers(roleId string, userIds []string) int64

	// DeleteAuthUsers 批量取消授权用户角色
	DeleteAuthUsers(roleId string, userIds []string) int64
}
