package repository

import (
	"mask_api_gin/src/internal/infra/database"
	"mask_api_gin/src/modules/job/model"
)

// NewSysJob 实例化数据层
var NewSysJob = &SysJob{}

// SysJob 调度任务 数据层处理
type SysJob struct{}

// sysJobSortMapping 调度任务排序字段映射
var sysJobSortMapping = map[string]string{
	"id":                 "id",
	"job_name":           "job_name",
	"job_group":          "job_group",
	"job_invoke_target":  "job_invoke_target",
	"job_cron":           "job_cron",
	"job_misfire_policy": "job_misfire_policy",
	"job_concurrent":     "job_concurrent",
	"job_status":         "job_status",
	"create_time":        "create_time",
	"update_time":        "update_time",
}

// SelectByPage 分页查询集合
func (r SysJob) SelectByPage(query map[string]string) ([]model.SysJob, int64, error) {
	tx := database.DB("").Model(&model.SysJob{})
	// 查询条件拼接
	if v, ok := query["job_name"]; ok && v != "" {
		tx = tx.Where("job_name like ?", v+"%")
	}
	if v, ok := query["job_group"]; ok && v != "" {
		tx = tx.Where("job_group = ?", v)
	}
	if v, ok := query["job_invoke_target"]; ok && v != "" {
		tx = tx.Where("job_invoke_target like ?", v+"%")
	}
	if v, ok := query["job_status"]; ok && v != "" {
		tx = tx.Where("job_status = ?", v)
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysJob{}

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
	orders := database.SortParse(query["sort"], sysJobSortMapping)
	tx = database.SortOrDefault(tx, orders, "id desc")

	err := tx.Find(&rows).Error
	if err != nil {
		return rows, total, err
	}
	return rows, total, nil
}

// Select 查询集合
func (r SysJob) Select(sysJob model.SysJob) ([]model.SysJob, error) {
	tx := database.DB("").Model(&model.SysJob{})
	// 查询条件拼接
	if sysJob.JobName != "" {
		tx = tx.Where("job_name like ?", sysJob.JobName+"%")
	}
	if sysJob.JobGroup != "" {
		tx = tx.Where("job_group = ?", sysJob.JobGroup)
	}
	if sysJob.JobInvokeTarget != "" {
		tx = tx.Where("job_invoke_target like ?", sysJob.JobInvokeTarget+"%")
	}
	if sysJob.JobStatus != "" {
		tx = tx.Where("job_status = ?", sysJob.JobStatus)
	}

	// 查询数据
	rows := []model.SysJob{}
	if err := tx.Order("id asc").Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// SelectByIds 通过ID查询信息
func (r SysJob) SelectByIds(ids []int64) ([]model.SysJob, error) {
	rows := []model.SysJob{}
	if len(ids) <= 0 {
		return rows, nil
	}
	tx := database.DB("").Model(&model.SysJob{})
	// 构建查询条件
	tx = tx.Where("id in ?", ids)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// Insert 新增信息 返回新增数据ID
func (r SysJob) Insert(sysJob model.SysJob) (int64, error) {
	// 执行插入
	if err := database.DB("").Create(&sysJob).Error; err != nil {
		return 0, err
	}
	return sysJob.ID, nil
}

// Update 修改信息
func (r SysJob) Update(sysJob model.SysJob) (int64, error) {
	if sysJob.ID <= 0 {
		return 0, nil
	}
	tx := database.DB("").Model(&model.SysJob{})
	// 构建查询条件
	tx = tx.Where("id = ?", sysJob.ID)
	tx = tx.Omit("id", "create_by", "create_time")
	// 执行更新
	if err := tx.Updates(sysJob).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// DeleteByIds 批量删除信息（物理删除）
func (r SysJob) DeleteByIds(ids []int64) (int64, error) {
	if len(ids) <= 0 {
		return 0, nil
	}
	tx := database.DB("").Where("id in ?", ids)
	if err := tx.Delete(&model.SysJob{}).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// CheckUnique 校验信息是否唯一
func (r SysJob) CheckUnique(sysJob model.SysJob) (int64, error) {
	tx := database.DB("").Model(&model.SysJob{})
	// 查询条件拼接
	if sysJob.JobName != "" {
		tx = tx.Where("job_name = ?", sysJob.JobName)
	}
	if sysJob.JobGroup != "" {
		tx = tx.Where("job_group = ?", sysJob.JobGroup)
	}

	// 查询数据
	var id int64 = 0
	if err := tx.Select("id").Limit(1).Find(&id).Error; err != nil {
		return id, err
	}
	return id, nil
}
