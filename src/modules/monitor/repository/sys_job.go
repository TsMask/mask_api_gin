package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/monitor/model"

	"time"
)

// NewSysJob 实例化数据层
var NewSysJob = &SysJob{}

// SysJob 调度任务 数据层处理
type SysJob struct{}

// SelectByPage 分页查询集合
func (r SysJob) SelectByPage(query map[string]any) ([]model.SysJob, int64) {
	tx := db.DB("").Model(&model.SysJob{})
	// 查询条件拼接
	if v, ok := query["jobName"]; ok && v != "" {
		tx = tx.Where("job_name like concat(?, '%')", v)
	}
	if v, ok := query["jobGroup"]; ok && v != "" {
		tx = tx.Where("job_group = ?", v)
	}
	if v, ok := query["invokeTarget"]; ok && v != "" {
		tx = tx.Where("invoke_target like concat(?, '%')", v)
	}
	if v, ok := query["statusFlag"]; ok && v != "" {
		tx = tx.Where("status_flag = ?", v)
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysJob{}

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
func (r SysJob) Select(sysJob model.SysJob) []model.SysJob {
	tx := db.DB("").Model(&model.SysJob{})
	// 查询条件拼接
	if sysJob.JobName != "" {
		tx = tx.Where("job_name like concat(?, '%')", sysJob.JobName)
	}
	if sysJob.JobGroup != "" {
		tx = tx.Where("job_group = ?", sysJob.JobGroup)
	}
	if sysJob.InvokeTarget != "" {
		tx = tx.Where("invoke_target like concat(?, '%')", sysJob.InvokeTarget)
	}
	if sysJob.StatusFlag != "" {
		tx = tx.Where("status_flag = ?", sysJob.StatusFlag)
	}

	// 查询数据
	rows := []model.SysJob{}
	if err := tx.Order("job_id asc").Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// SelectByIds 通过ID查询信息
func (r SysJob) SelectByIds(jobIds []string) []model.SysJob {
	rows := []model.SysJob{}
	if len(jobIds) <= 0 {
		return rows
	}
	tx := db.DB("").Model(&model.SysJob{})
	// 构建查询条件
	tx = tx.Where("job_id in ?", jobIds)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// Insert 新增信息 返回新增数据ID
func (r SysJob) Insert(sysJob model.SysJob) string {
	if sysJob.CreateBy != "" {
		ms := time.Now().UnixMilli()
		sysJob.UpdateBy = sysJob.CreateBy
		sysJob.UpdateTime = ms
		sysJob.CreateTime = ms

	}
	// 执行插入
	if err := db.DB("").Create(&sysJob).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return ""
	}
	return sysJob.JobId
}

// Update 修改信息
func (r SysJob) Update(sysJob model.SysJob) int64 {
	if sysJob.JobId == "" {
		return 0
	}
	if sysJob.UpdateBy != "" {
		sysJob.UpdateTime = time.Now().UnixMilli()
	}
	tx := db.DB("").Model(&model.SysJob{})
	// 构建查询条件
	tx = tx.Where("job_id = ?", sysJob.JobId)
	// 执行更新
	if err := tx.Updates(sysJob).Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByIds 批量删除信息
func (r SysJob) DeleteByIds(jobIds []string) int64 {
	if len(jobIds) <= 0 {
		return 0
	}
	tx := db.DB("").Where("job_id in ?", jobIds)
	if err := tx.Delete(&model.SysJob{}).Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// CheckUniqueJob 校验信息是否唯一
func (r SysJob) CheckUniqueJob(sysJob model.SysJob) string {
	tx := db.DB("").Model(&model.SysJob{})
	// 查询条件拼接
	if sysJob.JobName != "" {
		tx = tx.Where("job_name = ?", sysJob.JobName)
	}
	if sysJob.JobGroup != "" {
		tx = tx.Where("job_group = ?", sysJob.JobGroup)
	}

	// 查询数据
	var id string = ""
	if err := tx.Select("job_id").Limit(1).Find(&id).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return id
	}
	return id
}
