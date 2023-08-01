package service

import (
	"errors"
	"fmt"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// 实例化服务层 SysPostImpl 结构体
var NewSysPostImpl = &SysPostImpl{
	sysPostRepository:     repository.NewSysPostImpl,
	sysUserPostRepository: repository.NewSysUserPostImpl,
}

// SysPostImpl 岗位表 服务层处理
type SysPostImpl struct {
	// 岗位服务
	sysPostRepository repository.ISysPost
	// 用户与岗位关联服务
	sysUserPostRepository repository.ISysUserPost
}

// SelectPostPage 查询岗位分页数据集合
func (r *SysPostImpl) SelectPostPage(query map[string]any) map[string]any {
	return r.sysPostRepository.SelectPostPage(query)
}

// SelectPostList 查询岗位数据集合
func (r *SysPostImpl) SelectPostList(sysPost model.SysPost) []model.SysPost {
	return r.sysPostRepository.SelectPostList(sysPost)
}

// SelectPostById 通过岗位ID查询岗位信息
func (r *SysPostImpl) SelectPostById(postId string) model.SysPost {
	if postId == "" {
		return model.SysPost{}
	}
	posts := r.sysPostRepository.SelectPostByIds([]string{postId})
	if len(posts) > 0 {
		return posts[0]
	}
	return model.SysPost{}
}

// SelectPostListByUserId 根据用户ID获取岗位选择框列表
func (r *SysPostImpl) SelectPostListByUserId(userId string) []model.SysPost {
	return r.sysPostRepository.SelectPostListByUserId(userId)
}

// DeletePostByIds 批量删除岗位信息
func (r *SysPostImpl) DeletePostByIds(postIds []string) (int64, error) {
	// 检查是否存在
	posts := r.sysPostRepository.SelectPostByIds(postIds)
	if len(posts) <= 0 {
		return 0, errors.New("没有权限访问岗位数据！")
	}
	for _, post := range posts {
		useCount := r.sysUserPostRepository.CountUserPostByPostId(post.PostID)
		if useCount > 0 {
			msg := fmt.Sprintf("【%s】已分配给用户,不能删除", post.PostName)
			return 0, errors.New(msg)
		}
	}
	if len(posts) == len(postIds) {
		rows := r.sysPostRepository.DeletePostByIds(postIds)
		return rows, nil
	}
	return 0, errors.New("删除岗位信息失败！")
}

// UpdatePost 修改岗位信息
func (r *SysPostImpl) UpdatePost(sysPost model.SysPost) int64 {
	return r.sysPostRepository.UpdatePost(sysPost)
}

// InsertPost 新增岗位信息
func (r *SysPostImpl) InsertPost(sysPost model.SysPost) string {
	return r.sysPostRepository.InsertPost(sysPost)
}

// CheckUniquePostName 校验岗位名称
func (r *SysPostImpl) CheckUniquePostName(postName, postId string) bool {
	uniqueId := r.sysPostRepository.CheckUniquePost(model.SysPost{
		PostName: postName,
	})
	if uniqueId == postId {
		return true
	}
	return uniqueId == ""
}

// CheckUniquePostCode 校验岗位编码
func (r *SysPostImpl) CheckUniquePostCode(postCode, postId string) bool {
	uniqueId := r.sysPostRepository.CheckUniquePost(model.SysPost{
		PostCode: postCode,
	})
	if uniqueId == postId {
		return true
	}
	return uniqueId == ""
}
