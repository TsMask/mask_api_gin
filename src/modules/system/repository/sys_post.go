package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"
)

// NewSysPost 实例化数据层
var NewSysPost = &SysPost{}

// SysPost 岗位表 数据层处理
type SysPost struct{}

// sysPostSortMapping 岗位排序字段映射
var sysPostSortMapping = map[string]string{
	"id":         "id",
	"createTime": "create_time",
	"updateTime": "update_time",
	"postCode":   "post_code",
	"postName":   "post_name",
	"postSort":   "post_sort",
	"postStatus": "post_status",
}

// SelectByPage 分页查询集合
func (r SysPost) SelectByPage(query map[string]string) ([]model.SysPost, int64) {
	tx := db.DB("").Model(&model.SysPost{})
	// 查询条件拼接
	if v, ok := query["postCode"]; ok && v != "" {
		tx = tx.Where("post_code like ?", v+"%")
	}
	if v, ok := query["postName"]; ok && v != "" {
		tx = tx.Where("post_name like ?", v+"%")
	}
	if v, ok := query["postStatus"]; ok && v != "" {
		tx = tx.Where("post_status = ?", v)
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
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序，默认按 post_sort ASC
	orders := db.SortParse(query["sort"], sysPostSortMapping)
	tx = db.SortOrDefault(tx, orders, "post_sort asc")

	err := tx.Find(&rows).Error
	if err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows, total
	}
	return rows, total
}

// Select 查询集合
func (r SysPost) Select(sysPost model.SysPost) []model.SysPost {
	tx := db.DB("").Model(&model.SysPost{})
	// 查询条件拼接
	if sysPost.PostCode != "" {
		tx = tx.Where("post_code like ?", sysPost.PostCode+"%")
	}
	if sysPost.PostName != "" {
		tx = tx.Where("post_name like ?", sysPost.PostName+"%")
	}
	if sysPost.PostStatus != "" {
		tx = tx.Where("post_status = ?", sysPost.PostStatus)
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
func (r SysPost) SelectByIds(postIds []int64) []model.SysPost {
	rows := []model.SysPost{}
	if len(postIds) <= 0 {
		return rows
	}
	tx := db.DB("").Model(&model.SysPost{})
	// 构建查询条件
	tx = tx.Where("id in ?", postIds)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// Insert 新增信息 返回新增数据ID
func (r SysPost) Insert(sysPost model.SysPost) int64 {
	// 执行插入
	if err := db.DB("").Create(&sysPost).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
	}
	return sysPost.ID
}

// Update 修改信息 返回受影响行数
func (r SysPost) Update(sysPost model.SysPost) int64 {
	if sysPost.ID <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.SysPost{})
	// 构建查询条件
	tx = tx.Where("id = ?", sysPost.ID)
	tx = tx.Omit("id", "create_by", "create_time")
	// 执行更新
	if err := tx.Updates(sysPost).Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByIds 批量删除信息 返回受影响行数（物理删除）
func (r SysPost) DeleteByIds(postIds []int64) int64 {
	if len(postIds) <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.SysPost{})
	// 构建查询条件
	tx = tx.Where("id in ?", postIds)
	// 执行物理删除
	if err := tx.Delete(&model.SysPost{}).Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// CheckUnique 检查信息是否唯一 返回ID
func (r SysPost) CheckUnique(sysPost model.SysPost) int64 {
	tx := db.DB("").Model(&model.SysPost{})
	// 查询条件拼接
	if sysPost.PostName != "" {
		tx = tx.Where("post_name= ?", sysPost.PostName)
	}
	if sysPost.PostCode != "" {
		tx = tx.Where("post_code = ?", sysPost.PostCode)
	}

	// 查询数据
	var id int64 = 0
	if err := tx.Select("id").Limit(1).Find(&id).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return id
	}
	return id
}

// SelectByUserId 根据用户ID获取岗位选择框列表
func (r SysPost) SelectByUserId(userId int64) []model.SysPost {
	rows := []model.SysPost{}
	if userId <= 0 {
		return rows
	}
	tx := db.DB("").Model(&model.SysPost{})
	// 构建查询条件（sys_post 主键 id；sys_user_post 外键 post_id/user_id）
	tx = tx.Where("id in (select post_id from sys_user_post where user_id = ?)", userId)
	// 查询数据
	if err := tx.Order("post_sort").Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}
