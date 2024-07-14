package service

import "mask_api_gin/src/modules/system/model"

// ISysPostService 岗位信息 服务层接口
type ISysPostService interface {
	// FindByPage 分页查询列表数据
	FindByPage(query map[string]any) map[string]any

	// Find 查询列表数据
	Find(sysPost model.SysPost) []model.SysPost

	// FindById 通过ID查询信息
	FindById(postId string) model.SysPost

	// Insert 新增信息
	Insert(sysPost model.SysPost) string

	// Update 修改信息
	Update(sysPost model.SysPost) int64

	// DeleteByIds 批量删除信息
	DeleteByIds(postIds []string) (int64, error)

	// CheckUniqueByName 检查岗位名称是否唯一
	CheckUniqueByName(postName, postId string) bool

	// CheckUniqueByCode 检查岗位编码是否唯一
	CheckUniqueByCode(postCode, postId string) bool

	// FindByUserId 根据用户ID获取岗位选择框列表
	FindByUserId(userId string) []model.SysPost
}
