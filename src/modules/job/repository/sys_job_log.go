package repository

import (
	"fmt"

	"mask_api_gin/src/internal/infra/database"
	"mask_api_gin/src/modules/job/model"
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
func (r SysJobLog) SelectByPage(query map[string]string) ([]model.SysJobLog, int64, error) {
	tx := database.DB("").Model(&model.SysJobLog{})
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
		return rows, total, nil
	}

	// 查询数据分页
	pageNum, pageSize := database.PageNumSize(query["pageNum"], query["pageSize"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序
	orders := database.SortParse(query["sort"], sysJobLogSortMapping)
	tx = database.SortOrDefault(tx, orders, "id desc")

	err := tx.Find(&rows).Error
	if err != nil {
		return rows, total, err
	}
	return rows, total, nil
}

// Select 查询集合
func (r SysJobLog) Select(sysJobLog model.SysJobLog) ([]model.SysJobLog, error) {
	tx := database.DB("").Model(&model.SysJobLog{})
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
		return rows, err
	}
	return rows, nil
}

// SelectById 通过ID查询信息
func (r SysJobLog) SelectById(id int64) (model.SysJobLog, error) {
	item := model.SysJobLog{}
	if id <= 0 {
		return item, nil
	}
	tx := database.DB("").Model(&model.SysJobLog{})
	// 构建查询条件
	tx = tx.Where("id = ?", id)
	// 查询数据
	if err := tx.Limit(1).Find(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

// Insert 新增信息
func (r SysJobLog) Insert(sysJobLog model.SysJobLog) (int64, error) {
	// 执行插入
	if err := database.DB("").Create(&sysJobLog).Error; err != nil {
		return 0, err
	}
	return sysJobLog.ID, nil
}

// DeleteByIds 批量删除信息
func (r SysJobLog) DeleteByIds(ids []int64) (int64, error) {
	if len(ids) <= 0 {
		return 0, nil
	}
	tx := database.DB("").Where("id in ?", ids)
	if err := tx.Delete(&model.SysJobLog{}).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// Clean 清空集合数据
func (r SysJobLog) Clean() (int64, error) {
	tx := database.DB("").Where("1 = 1").Delete(&model.SysJobLog{})
	if err := tx.Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}
