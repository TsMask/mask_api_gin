package service

import (
	"fmt"
	"path/filepath"
	"time"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
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
func (s SysPost) FindByPage(query map[string]string) ([]model.SysPost, int64) {
	rows, total, err := s.sysPostRepository.SelectByPage(query)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return rows, total
}

// Find 查询列表数据
func (s SysPost) Find(sysPost model.SysPost) []model.SysPost {
	rows, err := s.sysPostRepository.Select(sysPost)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return rows
}

// FindById 通过ID查询信息
func (s SysPost) FindById(postId int64) model.SysPost {
	if postId <= 0 {
		return model.SysPost{}
	}
	posts, err := s.sysPostRepository.SelectByIds([]int64{postId})
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	if len(posts) > 0 {
		return posts[0]
	}
	return model.SysPost{}
}

// Insert 新增信息
func (s SysPost) Insert(sysPost model.SysPost) int64 {
	if sysPost.CreateBy != "" {
		ms := time.Now().UnixMilli()
		sysPost.UpdateBy = sysPost.CreateBy
		sysPost.UpdateTime = ms
		sysPost.CreateTime = ms
	}
	rows, err := s.sysPostRepository.Insert(sysPost)
	if err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
	}
	return rows
}

// Update 修改信息
func (s SysPost) Update(sysPost model.SysPost) int64 {
	if sysPost.UpdateBy != "" {
		sysPost.UpdateTime = time.Now().UnixMilli()
	}
	rows, err := s.sysPostRepository.Update(sysPost)
	if err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return rows
}

// DeleteByIds 批量删除信息
func (s SysPost) DeleteByIds(postIds []int64) (int64, error) {
	// 检查是否存在
	posts, err := s.sysPostRepository.SelectByIds(postIds)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	if len(posts) <= 0 {
		return 0, fmt.Errorf("没有权限访问岗位数据！")
	}
	for _, post := range posts {
		useCount, err := s.sysUserPostRepository.ExistUserByPostId(post.ID)
		if err != nil {
			logger.Errorf("query err => %v", err.Error())
		}
		if useCount > 0 {
			return 0, fmt.Errorf("【%s】已分配给用户,不能删除", post.PostName)
		}
	}
	if len(posts) == len(postIds) {
		return s.sysPostRepository.DeleteByIds(postIds)
	}
	return 0, fmt.Errorf("删除岗位信息失败！")
}

// CheckUniqueByName 检查岗位名称是否唯一
func (s SysPost) CheckUniqueByName(postName string, postId int64) bool {
	uniqueId, err := s.sysPostRepository.CheckUnique(model.SysPost{
		PostName: postName,
	})
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	if uniqueId == postId {
		return true
	}
	return uniqueId == 0
}

// CheckUniqueByCode 检查岗位编码是否唯一
func (s SysPost) CheckUniqueByCode(postCode string, postId int64) bool {
	uniqueId, err := s.sysPostRepository.CheckUnique(model.SysPost{
		PostCode: postCode,
	})
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	if uniqueId == postId {
		return true
	}
	return uniqueId == 0
}

// FindByUserId 根据用户ID获取岗位选择框列表
func (s SysPost) FindByUserId(userId int64) []model.SysPost {
	rows, err := s.sysPostRepository.SelectByUserId(userId)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return rows
}

// ExportData 导出数据表格
func (s SysPost) ExportData(rows []model.SysPost, fileName string) (string, error) {
	// 第一行表头标题
	data := [][]string{
		{"岗位编号", "岗位编码", "岗位名称", "岗位排序", "岗位状态"},
	}
	// 从第二行开始的数据
	for _, row := range rows {
		statusValue := "停用"
		if row.PostStatus == constants.STATUS_YES {
			statusValue = "正常"
		}
		data = append(data, []string{
			fmt.Sprintf("%d", row.ID),
			row.PostCode,
			row.PostName,
			fmt.Sprintf("%d", row.PostSort),
			statusValue,
		})
	}

	// 导出文件路径
	_, dir := file.UploadResourceDir()
	datePath := date.ParseDatePath()
	filePath := filepath.Join(dir, constants.UPLOAD_EXPORT, datePath, fileName)

	// 导出数据表格
	return filePath, file.CSVWrite(filePath, data, false)
}
