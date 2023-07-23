package service

import "mask_api_gin/src/modules/system/model"

// ISysPost 岗位信息 服务层接口
type ISysPost interface {
	// SelectPostPage 查询岗位分页数据集合
	SelectPostPage(query map[string]string) map[string]interface{}

	// SelectPostList 查询岗位数据集合
	SelectPostList(sysPost model.SysPost) []model.SysPost

	// SelectPostById 通过岗位ID查询岗位信息
	SelectPostById(postId string) model.SysPost

	// SelectPostListByUserId 根据用户ID获取岗位选择框列表
	SelectPostListByUserId(userId string) []model.SysPost

	// DeletePostByIds 批量删除岗位信息
	DeletePostByIds(postIds []string) (int64, error)

	// UpdatePost 修改岗位信息
	UpdatePost(sysPost model.SysPost) int64

	// InsertPost 新增岗位信息
	InsertPost(sysPost model.SysPost) string

	// CheckUniquePostName 校验岗位名称
	CheckUniquePostName(postName, postId string) bool

	// CheckUniquePostCode 校验岗位编码
	CheckUniquePostCode(postCode, postId string) bool
}
