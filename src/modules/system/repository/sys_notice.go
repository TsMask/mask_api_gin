package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"

	"fmt"
	"time"
)

// NewSysNotice 实例化数据层
var NewSysNotice = &SysNotice{}

// SysNotice 通知公告表 数据层处理
type SysNotice struct{}

// SelectByPage 分页查询集合
func (r SysNotice) SelectByPage(query map[string]string) ([]model.SysNotice, int64) {
	tx := db.DB("").Model(&model.SysNotice{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if v, ok := query["noticeTitle"]; ok && v != "" {
		tx = tx.Where("notice_title like concat(?, '%')", v)
	}
	if v, ok := query["noticeType"]; ok && v != "" {
		tx = tx.Where("notice_type = ?", v)
	}
	if v, ok := query["createBy"]; ok && v != "" {
		tx = tx.Where("create_by like concat(?, '%')", v)
	}
	if v, ok := query["statusFlag"]; ok && v != "" {
		tx = tx.Where("status_flag = ?", v)
	}
	if v, ok := query["beginTime"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s000", v)
			tx = tx.Where("create_time >= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("create_time >= ?", v)
		}
	}
	if v, ok := query["endTime"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s000", v)
			tx = tx.Where("create_time <= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("create_time <= ?", v)
		}
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysNotice{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return rows, total
	}

	// 查询数据分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)
	err := tx.Find(&rows).Error
	if err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows, total
	}
	return rows, total
}

// Select 查询集合
func (r SysNotice) Select(sysNotice model.SysNotice) []model.SysNotice {
	tx := db.DB("").Model(&model.SysNotice{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if sysNotice.NoticeTitle != "" {
		tx = tx.Where("notice_title like concat(?, '%')", sysNotice.NoticeTitle)
	}
	if sysNotice.NoticeType != "" {
		tx = tx.Where("notice_type = ?", sysNotice.NoticeType)
	}
	if sysNotice.CreateBy != "" {
		tx = tx.Where("create_by like concat(?, '%')", sysNotice.CreateBy)
	}
	if sysNotice.StatusFlag != "" {
		tx = tx.Where("status_flag = ?", sysNotice.StatusFlag)
	}

	// 查询数据
	rows := []model.SysNotice{}
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// SelectByIds 通过ID查询信息
func (r SysNotice) SelectByIds(noticeIds []string) []model.SysNotice {
	rows := []model.SysNotice{}
	if len(noticeIds) <= 0 {
		return rows
	}
	tx := db.DB("").Model(&model.SysNotice{})
	// 构建查询条件
	tx = tx.Where("notice_id in ? and del_flag = '0'", noticeIds)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// Insert 新增信息 返回新增数据ID
func (r SysNotice) Insert(sysNotice model.SysNotice) string {
	sysNotice.DelFlag = "0"
	if sysNotice.CreateBy != "" {
		ms := time.Now().UnixMilli()
		sysNotice.UpdateBy = sysNotice.CreateBy
		sysNotice.UpdateTime = ms
		sysNotice.CreateTime = ms
	}
	// 执行插入
	if err := db.DB("").Create(&sysNotice).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return ""
	}
	return sysNotice.NoticeId
}

// Update 修改信息 返回受影响行数
func (r SysNotice) Update(sysNotice model.SysNotice) int64 {
	if sysNotice.NoticeId == "" {
		return 0
	}
	if sysNotice.UpdateBy != "" {
		sysNotice.UpdateTime = time.Now().UnixMilli()
	}
	tx := db.DB("").Model(&model.SysNotice{})
	// 构建查询条件
	tx = tx.Where("notice_id = ?", sysNotice.NoticeId)
	tx = tx.Omit("notice_id", "del_flag", "create_by", "create_time")
	// 执行更新
	if err := tx.Updates(sysNotice).Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByIds 批量删除信息 返回受影响行数
func (r SysNotice) DeleteByIds(noticeIds []string) int64 {
	if len(noticeIds) <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.SysNotice{})
	// 构建查询条件
	tx = tx.Where("notice_id in ?", noticeIds)
	// 执行更新删除标记
	if err := tx.Update("del_flag", "1").Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}
