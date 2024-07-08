package repository

import "mask_api_gin/src/modules/system/model"

// ISysRoleRepository 角色表 数据层接口
type ISysRoleRepository interface {
	// SelectByPage 分页查询集合
	SelectByPage(query map[string]any, dataScopeSQL string) map[string]any

	// Select 查询集合
	Select(sysRole model.SysRole, dataScopeSQL string) []model.SysRole

	// SelectByIds 通过ID查询信息
	SelectByIds(roleIds []string) []model.SysRole

	// Update 修改信息
	Update(sysRole model.SysRole) int64

	// Insert 新增信息
	Insert(sysRole model.SysRole) string

	// DeleteByIds 批量删除信息
	DeleteByIds(roleIds []string) int64

	// SelectByUserId 根据用户ID获取角色选择框列表
	SelectByUserId(userId string) []model.SysRole

	// CheckUnique 检查信息是否唯一
	CheckUnique(sysRole model.SysRole) string
}
