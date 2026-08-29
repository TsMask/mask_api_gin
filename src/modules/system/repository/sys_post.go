package repository

import (
	"mask_api_gin/src/internal/infra/database"
	"mask_api_gin/src/modules/system/model"
)

// NewSysPost 实例化数据层
var NewSysPost = &SysPost{}

// SysPost 岗位表 数据层处理
type SysPost struct{}

// sysPostSortMapping 岗位排序字段映射
var sysPostSortMapping = map[string]string{
	"id":          "id",
	"create_time": "create_time",
	"update_time": "update_time",
	"post_code":   "post_code",
	"post_name":   "post_name",
	"post_sort":   "post_sort",
	"post_status": "post_status",
}

// SelectByPage 分页查询集合
func (r SysPost) SelectByPage(query map[string]string) ([]model.SysPost, int64, error) {
	tx := database.DB("").Model(&model.SysPost{})
	// 查询条件拼接
	if v, ok := query["post_code"]; ok && v != "" {
		tx = tx.Where("post_code like ?", v+"%")
	}
	if v, ok := query["post_name"]; ok && v != "" {
		tx = tx.Where("post_name like ?", v+"%")
	}
	if v, ok := query["post_status"]; ok && v != "" {
		tx = tx.Where("post_status = ?", v)
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysPost{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil {
		return rows, total, err
	}
	if total <= 0 {
		return rows, total, nil
	}

	// 查询数据分页
	pageNum, pageSize := database.PageNumSize(query["page_num"], query["page_size"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序，默认按 post_sort ASC
	orders := database.SortParse(query["sort"], sysPostSortMapping)
	tx = database.SortOrDefault(tx, orders, "post_sort asc")

	err := tx.Find(&rows).Error
	if err != nil {
		return rows, total, err
	}
	return rows, total, nil
}

// Select 查询集合
func (r SysPost) Select(sysPost model.SysPost) ([]model.SysPost, error) {
	tx := database.DB("").Model(&model.SysPost{})
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
		return rows, err
	}
	return rows, nil
}

// SelectByIds 通过ID查询信息
func (r SysPost) SelectByIds(postIds []int64) ([]model.SysPost, error) {
	rows := []model.SysPost{}
	if len(postIds) <= 0 {
		return rows, nil
	}
	tx := database.DB("").Model(&model.SysPost{})
	// 构建查询条件
	tx = tx.Where("id in ?", postIds)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// Insert 新增信息 返回新增数据ID
func (r SysPost) Insert(sysPost model.SysPost) (int64, error) {
	// 执行插入
	if err := database.DB("").Create(&sysPost).Error; err != nil {
		return 0, err
	}
	return sysPost.ID, nil
}

// Update 修改信息 返回受影响行数
func (r SysPost) Update(sysPost model.SysPost) (int64, error) {
	if sysPost.ID <= 0 {
		return 0, nil
	}
	tx := database.DB("").Model(&model.SysPost{})
	// 构建查询条件
	tx = tx.Where("id = ?", sysPost.ID)
	tx = tx.Omit("id", "create_by", "create_time")
	// 执行更新
	if err := tx.Updates(sysPost).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// DeleteByIds 批量删除信息 返回受影响行数（物理删除）
func (r SysPost) DeleteByIds(postIds []int64) (int64, error) {
	if len(postIds) <= 0 {
		return 0, nil
	}
	tx := database.DB("").Model(&model.SysPost{})
	// 构建查询条件
	tx = tx.Where("id in ?", postIds)
	// 执行物理删除
	if err := tx.Delete(&model.SysPost{}).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// CheckUnique 检查信息是否唯一 返回ID
func (r SysPost) CheckUnique(sysPost model.SysPost) (int64, error) {
	tx := database.DB("").Model(&model.SysPost{})
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
		return id, err
	}
	return id, nil
}

// SelectByUserId 根据用户ID获取岗位选择框列表
func (r SysPost) SelectByUserId(userId int64) ([]model.SysPost, error) {
	rows := []model.SysPost{}
	if userId <= 0 {
		return rows, nil
	}
	tx := database.DB("").Model(&model.SysPost{})
	// 构建查询条件（sys_post 主键 id；sys_user_post 外键 post_id/user_id）
	tx = tx.Where("id in (select post_id from sys_user_post where user_id = ?)", userId)
	// 查询数据
	if err := tx.Order("post_sort").Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}
