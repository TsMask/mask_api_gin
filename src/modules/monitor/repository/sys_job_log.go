package repository

import (
	"fmt"

	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/monitor/model"
)

// NewSysJobLog 实例化数据层
var NewSysJobLog = &SysJobLog{}

// SysJobLog 调度任务日志表 数据层处理
type SysJobLog struct{}

// sysJobLogSortMapping 调度任务日志排序字段映射
var sysJobLogSortMapping = map[string]string{
	"id":              "id",
	"jobName":         "job_name",
	"jobGroup":        "job_group",
	"jobInvokeTarget": "job_invoke_target",
	"logStatus":       "log_status",
	"createTime":      "create_time",
}

// SelectByPage 分页查询集合
func (r SysJobLog) SelectByPage(query map[string]string) ([]model.SysJobLog, int64) {
	tx := db.DB("").Model(&model.SysJobLog{})
	// 查询条件拼接
	if v, ok := query["jobName"]; ok && v != "" {
		tx = tx.Where("job_name = ?", v)
	}
	if v, ok := query["jobGroup"]; ok && v != "" {
		tx = tx.Where("job_group = ?", v)
	}
	if v, ok := query["logStatus"]; ok && v != "" {
		tx = tx.Where("log_status = ?", v)
	}
	if v, ok := query["jobInvokeTarget"]; ok && v != "" {
		tx = tx.Where("job_invoke_target like ?", v+"%")
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
			v = fmt.Sprintf("%s999", v)
			tx = tx.Where("create_time <= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("create_time <= ?", v)
		}
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysJobLog{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return rows, total
	}

	// 查询数据分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序
	orders := db.SortParse(query["sort"], sysJobLogSortMapping)
	tx = db.SortOrDefault(tx, orders, "id desc")

	err := tx.Find(&rows).Error
	if err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows, total
	}
	return rows, total
}

// Select 查询集合
func (r SysJobLog) Select(sysJobLog model.SysJobLog) []model.SysJobLog {
	tx := db.DB("").Model(&model.SysJobLog{})
	// 查询条件拼接
	if sysJobLog.JobName != "" {
		tx = tx.Where("job_name like ?", sysJobLog.JobName+"%")
	}
	if sysJobLog.JobGroup != "" {
		tx = tx.Where("job_group = ?", sysJobLog.JobGroup)
	}
	if sysJobLog.LogStatus != "" {
		tx = tx.Where("log_status = ?", sysJobLog.LogStatus)
	}
	if sysJobLog.JobInvokeTarget != "" {
		tx = tx.Where("job_invoke_target like ?", sysJobLog.JobInvokeTarget+"%")
	}

	// 查询数据
	rows := []model.SysJobLog{}
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// SelectById 通过ID查询信息
func (r SysJobLog) SelectById(id int64) model.SysJobLog {
	item := model.SysJobLog{}
	if id <= 0 {
		return item
	}
	tx := db.DB("").Model(&model.SysJobLog{})
	// 构建查询条件
	tx = tx.Where("id = ?", id)
	// 查询数据
	if err := tx.Limit(1).Find(&item).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return item
	}
	return item
}

// Insert 新增信息
func (r SysJobLog) Insert(sysJobLog model.SysJobLog) int64 {
	// 执行插入
	if err := db.DB("").Create(&sysJobLog).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
	}
	return sysJobLog.ID
}

// DeleteByIds 批量删除信息
func (r SysJobLog) DeleteByIds(ids []int64) int64 {
	if len(ids) <= 0 {
		return 0
	}
	tx := db.DB("").Where("id in ?", ids)
	if err := tx.Delete(&model.SysJobLog{}).Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// Clean 清空集合数据
func (r SysJobLog) Clean() int64 {
	tx := db.DB("").Where("1 = 1").Delete(&model.SysJobLog{})
	if err := tx.Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}
