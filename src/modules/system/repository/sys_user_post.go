package repository

import "mask_api_gin/src/modules/system/model"

// ISysUserPostRepository 用户与岗位关联表 数据层接口
type ISysUserPostRepository interface {
	// ExistUserByPostId 存在用户使用数量
	ExistUserByPostId(postId string) int64

	// DeleteByUserIds 批量删除关联By用户
	DeleteByUserIds(userIds []string) int64

	// BatchInsert 批量新增信息
	BatchInsert(arr []model.SysUserPost) int64
}
