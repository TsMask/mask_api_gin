package repository

import "mask_api_gin/src/modules/system/model"

// ISysPostRepository 岗位表 数据层接口
type ISysPostRepository interface {
	// SelectByPage 分页查询集合
	SelectByPage(query map[string]any) map[string]any

	// Select 查询集合
	Select(sysPost model.SysPost) []model.SysPost

	// SelectByIds 通过ID查询信息
	SelectByIds(postIds []string) []model.SysPost

	// Insert 新增信息
	Insert(sysPost model.SysPost) string

	// Update 修改信息
	Update(sysPost model.SysPost) int64

	// DeleteByIds 批量删除信息
	DeleteByIds(postIds []string) int64

	// SelectByUserId 根据用户ID获取岗位选择框列表
	SelectByUserId(userId string) []model.SysPost

	// CheckUnique 检查信息是否唯一
	CheckUnique(sysPost model.SysPost) string
}
