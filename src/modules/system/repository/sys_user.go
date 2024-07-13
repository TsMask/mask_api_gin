package repository

import "mask_api_gin/src/modules/system/model"

// ISysUserRepository 用户表 数据层接口
type ISysUserRepository interface {
	// SelectByPage 分页查询集合
	SelectByPage(queryMap map[string]any, dataScopeSQL string) map[string]any

	// Select 查询集合
	Select(sysUser model.SysUser, dataScopeSQL string) []model.SysUser

	// SelectByIds 通过ID查询信息
	SelectByIds(userIds []string) []model.SysUser

	// Insert 新增信息
	Insert(sysUser model.SysUser) string

	// Update 修改信息
	Update(sysUser model.SysUser) int64

	// DeleteByIds 批量删除信息
	DeleteByIds(userIds []string) int64

	// CheckUnique 检查信息是否唯一
	CheckUnique(sysUser model.SysUser) string

	// SelectByUserName 通过登录账号查询信息
	SelectByUserName(userName string) model.SysUser
	
	// SelectAllocatedByPage 分页查询集合By分配用户角色
	SelectAllocatedByPage(query map[string]any, dataScopeSQL string) map[string]any
}
