package repository

import "mask_api_gin/src/modules/system/model"

// SysUserPostImpl 用户与岗位关联表 数据层处理
var SysUserPostImpl = &sysUserPostImpl{
	selectSql: "",
}

type sysUserPostImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// CountUserPostByPostId 通过岗位ID查询岗位使用数量
func (r *sysUserPostImpl) CountUserPostByPostId(postId string) int {
	return 0
}

// DeleteUserPost 批量删除用户和岗位关联
func (r *sysUserPostImpl) DeleteUserPost(userIds []string) int {
	return 0
}

// BatchUserPost 批量新增用户岗位信息
func (r *sysUserPostImpl) BatchUserPost(sysUserPosts []model.SysUserPost) int {
	return 0
}
