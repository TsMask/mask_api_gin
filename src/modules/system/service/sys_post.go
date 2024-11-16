package service

import (
	"fmt"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysPost 实例化服务层
var NewSysPost = &SysPost{
	sysPostRepository:     repository.NewSysPost,
	sysUserPostRepository: repository.NewSysUserPost,
}

// SysPostService 岗位表 服务层处理
type SysPost struct {
	sysPostRepository     *repository.SysPost     // 岗位服务
	sysUserPostRepository *repository.SysUserPost // 用户与岗位关联服务
}

// FindByPage 分页查询列表数据
func (s SysPost) FindByPage(query map[string]any) ([]model.SysPost, int64) {
	return s.sysPostRepository.SelectByPage(query)
}

// Find 查询列表数据
func (s SysPost) Find(sysPost model.SysPost) []model.SysPost {
	return s.sysPostRepository.Select(sysPost)
}

// FindById 通过ID查询信息
func (s SysPost) FindById(postId int64) model.SysPost {
	if postId == 0 {
		return model.SysPost{}
	}
	posts := s.sysPostRepository.SelectByIds([]int64{postId})
	if len(posts) > 0 {
		return posts[0]
	}
	return model.SysPost{}
}

// Insert 新增信息
func (s SysPost) Insert(sysPost model.SysPost) int64 {
	return s.sysPostRepository.Insert(sysPost)
}

// Update 修改信息
func (s SysPost) Update(sysPost model.SysPost) int64 {
	return s.sysPostRepository.Update(sysPost)
}

// DeleteByIds 批量删除信息
func (s SysPost) DeleteByIds(postIds []int64) (int64, error) {
	// 检查是否存在
	posts := s.sysPostRepository.SelectByIds(postIds)
	if len(posts) <= 0 {
		return 0, fmt.Errorf("没有权限访问岗位数据！")
	}
	for _, post := range posts {
		if useCount := s.sysUserPostRepository.ExistUserByPostId(post.PostId); useCount > 0 {
			return 0, fmt.Errorf("【%s】已分配给用户,不能删除", post.PostName)
		}
	}
	if len(posts) == len(postIds) {
		return s.sysPostRepository.DeleteByIds(postIds), nil
	}
	return 0, fmt.Errorf("删除岗位信息失败！")
}

// CheckUniqueByName 检查岗位名称是否唯一
func (s SysPost) CheckUniqueByName(postName string, postId int64) bool {
	uniqueId := s.sysPostRepository.CheckUnique(model.SysPost{
		PostName: postName,
	})
	if uniqueId == postId {
		return true
	}
	return uniqueId == 0
}

// CheckUniqueByCode 检查岗位编码是否唯一
func (s SysPost) CheckUniqueByCode(postCode string, postId int64) bool {
	uniqueId := s.sysPostRepository.CheckUnique(model.SysPost{
		PostCode: postCode,
	})
	if uniqueId == postId {
		return true
	}
	return uniqueId == 0
}

// FindByUserId 根据用户ID获取岗位选择框列表
func (s SysPost) FindByUserId(userId int64) []model.SysPost {
	return s.sysPostRepository.SelectByUserId(userId)
}
