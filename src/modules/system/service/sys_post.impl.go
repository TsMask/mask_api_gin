package service

import (
	"fmt"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysPost 实例化服务层
var NewSysPost = &SysPostService{
	sysPostRepository:     repository.NewSysPost,
	sysUserPostRepository: repository.NewSysUserPost,
}

// SysPostService 岗位表 服务层处理
type SysPostService struct {
	sysPostRepository     repository.ISysPostRepository     // 岗位服务
	sysUserPostRepository repository.ISysUserPostRepository // 用户与岗位关联服务
}

// FindByPage 分页查询列表数据
func (r *SysPostService) FindByPage(query map[string]any) map[string]any {
	return r.sysPostRepository.SelectByPage(query)
}

// Find 查询列表数据
func (r *SysPostService) Find(sysPost model.SysPost) []model.SysPost {
	return r.sysPostRepository.Select(sysPost)
}

// FindById 通过ID查询信息
func (r *SysPostService) FindById(postId string) model.SysPost {
	if postId == "" {
		return model.SysPost{}
	}
	posts := r.sysPostRepository.SelectByIds([]string{postId})
	if len(posts) > 0 {
		return posts[0]
	}
	return model.SysPost{}
}

// Insert 新增信息
func (r *SysPostService) Insert(sysPost model.SysPost) string {
	return r.sysPostRepository.Insert(sysPost)
}

// Update 修改信息
func (r *SysPostService) Update(sysPost model.SysPost) int64 {
	return r.sysPostRepository.Update(sysPost)
}

// DeleteByIds 批量删除信息
func (r *SysPostService) DeleteByIds(postIds []string) (int64, error) {
	// 检查是否存在
	posts := r.sysPostRepository.SelectByIds(postIds)
	if len(posts) <= 0 {
		return 0, fmt.Errorf("没有权限访问岗位数据！")
	}
	for _, post := range posts {
		if useCount := r.sysUserPostRepository.ExistUserByPostId(post.PostID); useCount > 0 {
			return 0, fmt.Errorf("【%s】已分配给用户,不能删除", post.PostName)
		}
	}
	if len(posts) == len(postIds) {
		return r.sysPostRepository.DeleteByIds(postIds), nil
	}
	return 0, fmt.Errorf("删除岗位信息失败！")
}

// CheckUniqueByName 检查岗位名称是否唯一
func (r *SysPostService) CheckUniqueByName(postName, postId string) bool {
	uniqueId := r.sysPostRepository.CheckUnique(model.SysPost{
		PostName: postName,
	})
	if uniqueId == postId {
		return true
	}
	return uniqueId == ""
}

// CheckUniqueByCode 检查岗位编码是否唯一
func (r *SysPostService) CheckUniqueByCode(postCode, postId string) bool {
	uniqueId := r.sysPostRepository.CheckUnique(model.SysPost{
		PostCode: postCode,
	})
	if uniqueId == postId {
		return true
	}
	return uniqueId == ""
}

// FindByUserId 根据用户ID获取岗位选择框列表
func (r *SysPostService) FindByUserId(userId string) []model.SysPost {
	return r.sysPostRepository.SelectByUserId(userId)
}
