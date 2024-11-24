package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"

	"time"
)

// NewSysPost 实例化数据层
var NewSysPost = &SysPost{}

// SysPost 岗位表 数据层处理
type SysPost struct{}

// SelectByPage 分页查询集合
func (r SysPost) SelectByPage(query map[string]any) ([]model.SysPost, int64) {
	tx := db.DB("").Model(&model.SysPost{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if v, ok := query["postCode"]; ok && v != "" {
		tx = tx.Where("post_code like concat(?, '%')", v)
	}
	if v, ok := query["postName"]; ok && v != "" {
		tx = tx.Where("post_name like concat(?, '%')", v)
	}
	if v, ok := query["statusFlag"]; ok && v != "" {
		tx = tx.Where("status_flag = ?", v)
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysPost{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return rows, total
	}

	// 查询数据分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	err := tx.Limit(pageSize).Offset(pageSize * pageNum).Find(&rows).Error
	if err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows, total
	}
	return rows, total
}

// Select 查询集合
func (r SysPost) Select(sysPost model.SysPost) []model.SysPost {
	tx := db.DB("").Model(&model.SysPost{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if sysPost.PostCode != "" {
		tx = tx.Where("post_code like concat(?, '%')", sysPost.PostCode)
	}
	if sysPost.PostName != "" {
		tx = tx.Where("post_name like concat(?, '%')", sysPost.PostName)
	}
	if sysPost.StatusFlag != "" {
		tx = tx.Where("status_flag = ?", sysPost.StatusFlag)
	}

	// 查询数据
	rows := []model.SysPost{}
	if err := tx.Order("post_sort asc").Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// SelectByIds 通过ID查询信息
func (r SysPost) SelectByIds(postIds []string) []model.SysPost {
	rows := []model.SysPost{}
	if len(postIds) <= 0 {
		return rows
	}
	tx := db.DB("").Model(&model.SysPost{})
	// 构建查询条件
	tx = tx.Where("post_id in ? and del_flag = '0'", postIds)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// Insert 新增信息 返回新增数据ID
func (r SysPost) Insert(sysPost model.SysPost) string {
	sysPost.DelFlag = "0"
	if sysPost.CreateBy != "" {
		ms := time.Now().UnixMilli()
		sysPost.UpdateBy = sysPost.CreateBy
		sysPost.UpdateTime = ms
		sysPost.CreateTime = ms
	}
	// 执行插入
	if err := db.DB("").Create(&sysPost).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return ""
	}
	return sysPost.PostId
}

// Update 修改信息 返回受影响行数
func (r SysPost) Update(sysPost model.SysPost) int64 {
	if sysPost.PostId == "" {
		return 0
	}
	if sysPost.UpdateBy != "" {
		sysPost.UpdateTime = time.Now().UnixMilli()
	}
	tx := db.DB("").Model(&model.SysPost{})
	// 构建查询条件
	tx = tx.Where("post_id = ?", sysPost.PostId)
	// 执行更新
	if err := tx.Updates(sysPost).Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByIds 批量删除信息 返回受影响行数
func (r SysPost) DeleteByIds(postIds []string) int64 {
	if len(postIds) <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.SysPost{})
	// 构建查询条件
	tx = tx.Where("post_id in ?", postIds)
	// 执行更新删除标记
	if err := tx.Update("del_flag", "1").Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// SelectByUserId 根据用户ID获取岗位选择框列表
func (r SysPost) SelectByUserId(userId string) []model.SysPost {
	rows := []model.SysPost{}
	if userId == "" {
		return rows
	}
	tx := db.DB("").Model(&model.SysPost{})
	// 构建查询条件
	tx = tx.Where("post_id in (select post_id from sys_user_post  where user_id = ?)", userId)

	// 查询数据
	if err := tx.Order("post_id").Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// CheckUnique 检查信息是否唯一 返回ID
func (r SysPost) CheckUnique(sysPost model.SysPost) string {
	tx := db.DB("").Model(&model.SysPost{})
	tx = tx.Where("del_flag = 0")
	// 查询条件拼接
	if sysPost.PostName != "" {
		tx = tx.Where("post_name= ?", sysPost.PostName)
	}
	if sysPost.PostCode != "" {
		tx = tx.Where("post_code = ?", sysPost.PostCode)
	}

	// 查询数据
	var id string = ""
	if err := tx.Select("post_id").Limit(1).Find(&id).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return id
	}
	return id
}
