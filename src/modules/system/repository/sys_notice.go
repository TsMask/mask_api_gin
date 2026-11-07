package repository

import (
	"fmt"

	"mask_api_gin/src/internal/infra/database"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/modules/system/model"
)

// NewSysNotice 实例化数据层
var NewSysNotice = &SysNotice{}

// SysNotice 通知公告表 数据层处理
type SysNotice struct{}

// sysNoticeSortMapping 通知公告排序字段映射
var sysNoticeSortMapping = map[string]string{
	"id":            "id",
	"create_by":     "create_by",
	"create_time":   "create_time",
	"update_time":   "update_time",
	"notice_title":  "notice_title",
	"notice_type":   "notice_type",
	"notice_status": "notice_status",
}

// SelectByPage 分页查询集合
func (r SysNotice) SelectByPage(query map[string]string) ([]model.SysNotice, int64, error) {
	tx := database.DB("").Model(&model.SysNotice{})
	// 查询条件拼接
	if v, ok := query["notice_title"]; ok && v != "" {
		tx = tx.Where("notice_title like ?", v+"%")
	}
	if v, ok := query["notice_type"]; ok && v != "" {
		tx = tx.Where("notice_type = ?", v)
	}
	if v, ok := query["create_by"]; ok && v != "" {
		tx = tx.Where("create_by like ?", v+"%")
	}
	if v, ok := query["notice_status"]; ok && v != "" {
		tx = tx.Where("notice_status = ?", v)
	}
	if v, ok := query["begin_time"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s000", v)
			tx = tx.Where("create_time >= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("create_time >= ?", v)
		}
	}
	if v, ok := query["end_time"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s999", v)
			tx = tx.Where("create_time <= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("create_time <= ?", v)
		}
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysNotice{}

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

	// 排序
	orders := database.SortParse(query["sort"], sysNoticeSortMapping)
	tx = database.SortOrDefault(tx, orders, "id desc")

	err := tx.Find(&rows).Error
	if err != nil {
		return rows, total, err
	}
	return rows, total, nil
}

// Select 查询集合
func (r SysNotice) Select(sysNotice model.SysNotice) ([]model.SysNotice, error) {
	tx := database.DB("").Model(&model.SysNotice{})
	// 查询条件拼接
	if sysNotice.NoticeTitle != "" {
		tx = tx.Where("notice_title like ?", sysNotice.NoticeTitle+"%")
	}
	if sysNotice.NoticeType != "" {
		tx = tx.Where("notice_type = ?", sysNotice.NoticeType)
	}
	if sysNotice.CreateBy != "" {
		tx = tx.Where("create_by like ?", sysNotice.CreateBy+"%")
	}
	if sysNotice.NoticeStatus != "" {
		tx = tx.Where("notice_status = ?", sysNotice.NoticeStatus)
	}

	// 查询数据
	rows := []model.SysNotice{}
	if err := tx.Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// SelectByIds 通过ID查询信息
func (r SysNotice) SelectByIds(noticeIds []int64) ([]model.SysNotice, error) {
	rows := []model.SysNotice{}
	if len(noticeIds) <= 0 {
		return rows, nil
	}
	tx := database.DB("").Model(&model.SysNotice{})
	// 构建查询条件
	tx = tx.Where("id in ?", noticeIds)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// Insert 新增信息 返回新增数据ID
func (r SysNotice) Insert(sysNotice model.SysNotice) (int64, error) {
	// 执行插入
	if err := database.DB("").Create(&sysNotice).Error; err != nil {
		return 0, err
	}
	return sysNotice.ID, nil
}

// Update 修改信息 返回受影响行数
func (r SysNotice) Update(sysNotice model.SysNotice) (int64, error) {
	if sysNotice.ID <= 0 {
		return 0, nil
	}
	tx := database.DB("").Model(&model.SysNotice{})
	// 构建查询条件
	tx = tx.Where("id = ?", sysNotice.ID)
	tx = tx.Omit("id", "create_by", "create_time")
	// 执行更新
	if err := tx.Updates(sysNotice).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// DeleteByIds 批量删除信息 返回受影响行数（软删除）
func (r SysNotice) DeleteByIds(noticeIds []int64, updateBy string) (int64, error) {
	if len(noticeIds) <= 0 {
		return 0, nil
	}
	tx := database.DB("").Model(&model.SysNotice{})
	// 构建查询条件
	tx = tx.Where("id in ?", noticeIds)
	tx = tx.Where("notice_status <> ?", constants.NOTICE_STATUS_DELETED)
	// 执行更新删除状态，已删除记录不再重复更新
	updates := map[string]interface{}{
		"notice_status": constants.NOTICE_STATUS_DELETED,
	}
	if updateBy != "" {
		updates["update_by"] = updateBy
	}
	if err := tx.Updates(updates).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}
