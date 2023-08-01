package repository

import (
	"mask_api_gin/src/modules/monitor/model"
)

// ISysJob 调度任务表 数据层接口
type ISysJob interface {
	// SelectJobPage 分页查询调度任务集合
	SelectJobPage(query map[string]any) map[string]any

	// SelectJobList 查询调度任务集合
	SelectJobList(sysJob model.SysJob) []model.SysJob

	// SelectJobByIds 通过调度ID查询调度任务信息
	SelectJobByIds(jobIds []string) []model.SysJob

	// CheckUniqueJob 校验调度任务是否唯一
	CheckUniqueJob(sysJob model.SysJob) string

	// InsertJob 新增调度任务信息
	InsertJob(sysJob model.SysJob) string

	// UpdateJob 修改调度任务信息
	UpdateJob(sysJob model.SysJob) int64

	// DeleteJobByIds 批量删除调度任务信息
	DeleteJobByIds(jobIds []string) int64
}
