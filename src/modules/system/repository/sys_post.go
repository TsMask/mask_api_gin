package repository

import "mask_api_gin/src/modules/system/model"

// ISysPost 岗位表 数据层接口
type ISysPost interface {
	// SelectPostPage 查询岗位分页数据集合
	SelectPostPage(query map[string]string) map[string]interface{}

	// SelectPostList 查询岗位数据集合
	SelectPostList(sysPost model.SysPost) []model.SysPost

	// SelectPostByIds 通过岗位ID查询岗位信息
	SelectPostByIds(postIds []string) []model.SysPost

	// SelectPostListByUserId 根据用户ID获取岗位选择框列表
	SelectPostListByUserId(userId string) []model.SysPost

	// DeletePostByIds 批量删除岗位信息
	DeletePostByIds(postIds []string) int64

	// UpdatePost 修改岗位信息
	UpdatePost(sysPost model.SysPost) int64

	// InsertPost 新增岗位信息
	InsertPost(sysPost model.SysPost) string

	// CheckUniquePost 校验岗位唯一
	CheckUniquePost(sysPost model.SysPost) string
}
