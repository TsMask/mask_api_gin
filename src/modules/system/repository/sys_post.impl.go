package repository

import "mask_api_gin/src/modules/system/model"

// SysPostImpl 岗位表 数据层处理
var SysPostImpl = &sysPostImpl{
	selectSql: "",
}

type sysPostImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// SelectPostPage 查询岗位分页数据集合
func (r *sysPostImpl) SelectPostPage(query map[string]string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectPostList 查询岗位数据集合
func (r *sysPostImpl) SelectPostList(sysPost model.SysPost) []model.SysPost {
	return []model.SysPost{}
}

// SelectPostById 通过岗位ID查询岗位信息
func (r *sysPostImpl) SelectPostById(postId string) model.SysPost {
	return model.SysPost{}
}

// SelectPostListByUserId 根据用户ID获取岗位选择框列表
func (r *sysPostImpl) SelectPostListByUserId(userId string) []string {
	return []string{}
}

// SelectPostsByUserName 查询用户所属岗位组
func (r *sysPostImpl) SelectPostsByUserName(userName string) []model.SysPost {
	return []model.SysPost{}
}

// DeletePostByIds 批量删除岗位信息
func (r *sysPostImpl) DeletePostByIds(postIds []string) int {
	return 0
}

// UpdatePost 修改岗位信息
func (r *sysPostImpl) UpdatePost(sysPost model.SysPost) int {
	return 0
}

// InsertPost 新增岗位信息
func (r *sysPostImpl) InsertPost(sysPost model.SysPost) string {
	return ""
}

// CheckUniquePostName 校验岗位名称
func (r *sysPostImpl) CheckUniquePostName(postName string) string {
	return ""
}

// CheckUniquePostCode 校验岗位编码
func (r *sysPostImpl) CheckUniquePostCode(postCode string) string {
	return ""
}
