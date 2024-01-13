package service

import (
	"mask_api_gin/src/modules/monitor/model"
)

// ISysJob 调度任务信息 服务层接口
type ISysJob interface {
	// SelectJobPage 分页查询调度任务集合
	SelectJobPage(query map[string]any) map[string]any

	// SelectJobList 查询调度任务集合
	SelectJobList(sysJob model.SysJob) []model.SysJob

	// SelectJobById 通过调度ID查询调度任务信息
	SelectJobById(jobId string) model.SysJob

	// CheckUniqueJobName 校验调度任务名称和组是否唯一
	CheckUniqueJobName(jobName, jobGroup, jobId string) bool

	// InsertJob 新增调度任务信息
	InsertJob(sysJob model.SysJob) string

	// UpdateJob 修改调度任务信息
	UpdateJob(sysJob model.SysJob) int64

	// DeleteJobByIds 批量删除调度任务信息
	DeleteJobByIds(jobIds []string) (int64, error)

	// RunQueueJob 立即运行一次调度任务
	RunQueueJob(sysJob model.SysJob) bool

	// ResetQueueJob 重置初始调度任务
	ResetQueueJob()
}
