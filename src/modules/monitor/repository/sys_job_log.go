package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/modules/monitor/model"

	"fmt"
)

// NewSysJobLog 实例化数据层
var NewSysJobLog = &SysJobLog{}

// SysJobLog 调度任务日志表 数据层处理
type SysJobLog struct{}

// SelectByPage 分页查询集合
func (r SysJobLog) SelectByPage(query map[string]any) ([]model.SysJobLog, int64) {
	tx := db.DB("").Model(&model.SysJobLog{})
	// 查询条件拼接
	if v, ok := query["jobName"]; ok && v != "" {
		tx = tx.Where("job_name = ?", v)
	}
	if v, ok := query["jobGroup"]; ok && v != "" {
		tx = tx.Where("job_group = ?", v)
	}
	if v, ok := query["statusFlag"]; ok && v != "" {
		tx = tx.Where("status_flag = ?", v)
	}
	if v, ok := query["invokeTarget"]; ok && v != "" {
		tx = tx.Where("invoke_target like concat(?, '%')", v)
	}
	if v, ok := query["beginTime"]; ok && v != "" {
		tx = tx.Where("create_time >= ?", v)
	}
	if v, ok := query["endTime"]; ok && v != "" {
		tx = tx.Where("create_time <= ?", v)
	}
	if v, ok := query["params[beginTime]"]; ok && v != "" {
		beginDate := date.ParseStrToDate(fmt.Sprint(v), date.YYYY_MM_DD)
		tx = tx.Where("create_time >= ?", beginDate.UnixMilli())
	}
	if v, ok := query["params[endTime]"]; ok && v != "" {
		endDate := date.ParseStrToDate(fmt.Sprint(v), date.YYYY_MM_DD)
		tx = tx.Where("create_time <= ?", endDate.UnixMilli())
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
	err := tx.Limit(pageSize).Offset(pageSize * pageNum).Find(&rows).Error
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
		tx = tx.Where("job_name like concat(?, '%')", sysJobLog.JobName)
	}
	if sysJobLog.JobGroup != "" {
		tx = tx.Where("job_group = ?", sysJobLog.JobGroup)
	}
	if sysJobLog.StatusFlag != "" {
		tx = tx.Where("status_flag = ?", sysJobLog.StatusFlag)
	}
	if sysJobLog.InvokeTarget != "" {
		tx = tx.Where("invoke_target like concat(?, '%')", sysJobLog.InvokeTarget)
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
func (r SysJobLog) SelectById(jobLogId string) model.SysJobLog {
	item := model.SysJobLog{}
	if jobLogId == "" {
		return item
	}
	tx := db.DB("").Model(&model.SysJobLog{})
	// 构建查询条件
	tx = tx.Where("log_id = ?", jobLogId)
	// 查询数据
	if err := tx.Limit(1).Find(&item).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return item
	}
	return item
}

// Insert 新增信息
func (r SysJobLog) Insert(sysJobLog model.SysJobLog) string {
	// 执行插入
	if err := db.DB("").Create(&sysJobLog).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return ""
	}
	return sysJobLog.LogId
}

// DeleteByIds 批量删除信息
func (r SysJobLog) DeleteByIds(logIds []string) int64 {
	if len(logIds) <= 0 {
		return 0
	}
	tx := db.DB("").Where("log_id in ?", logIds)
	if err := tx.Delete(&model.SysJobLog{}).Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// Clean 清空集合数据
func (r SysJobLog) Clean() error {
	sql := "truncate table sys_job_log"
	_, err := db.ExecDB("", sql, nil)
	return err
}
