package service

import (
	"mask_api_gin/src/modules/monitor/model"
)

// ISysJobService 调度任务信息 服务层接口
type ISysJobService interface {
	// SelectJobPage 分页查询调度任务集合
	SelectJobPage(query map[string]string) map[string]interface{}

	// SelectJobList 查询调度任务集合
	SelectJobList(sysJob model.SysJob) []model.SysJob

	// SelectJobById 通过调度ID查询调度任务信息
	SelectJobById(jobId string) model.SysJob

	// CheckUniqueJob 校验调度任务名称和组是否唯一
	CheckUniqueJob(sysJob model.SysJob) bool

	// InsertJob 新增调度任务信息
	InsertJob(sysJob model.SysJob) string

	// UpdateJob 修改调度任务信息
	UpdateJob(sysJob model.SysJob) int

	// DeleteJobByIds 批量删除调度任务信息
	DeleteJobByIds(jobIds []string) int

	// ChangeStatus 任务调度状态修改
	ChangeStatus(sysJob model.SysJob) bool

	// InsertQueueJob 添加调度任务
	InsertQueueJob(sysJob model.SysJob, repeat bool) bool

	// DeleteQueueJob 删除调度任务
	DeleteQueueJob(sysJob model.SysJob) error

	// RunQueueJob 立即运行一次调度任务
	RunQueueJob(sysJob model.SysJob) bool

	// ResetQueueJob 重置初始调度任务
	ResetQueueJob() error
}
