package repository

import "mask_api_gin/src/modules/system/model"

// ISysUserPost 用户与岗位关联表 数据层接口
type ISysUserPost interface {
	// CountUserPostByPostId 通过岗位ID查询岗位使用数量
	CountUserPostByPostId(postId string) int64

	// BatchUserPost 批量新增用户岗位信息
	BatchUserPost(sysUserPosts []model.SysUserPost) int64

	// DeleteUserPost 批量删除用户和岗位关联
	DeleteUserPost(userIds []string) int64
}
