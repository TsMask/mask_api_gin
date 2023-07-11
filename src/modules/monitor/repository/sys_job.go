package repository

import (
	"mask_api_gin/src/modules/monitor/model"
)

// ISysJob 调度任务信息 数据层接口
type ISysJob interface {
	// SelectJobPage 分页查询调度任务集合
	SelectJobPage(query map[string]string) map[string]interface{}

	// SelectJobList 查询调度任务集合
	SelectJobList(sysJob model.SysJob) []model.SysJob

	// SelectJobById 通过调度ID查询调度任务信息
	SelectJobById(jobId string) model.SysJob

	// CheckUniqueJob 校验调度任务名称和组是否唯一
	CheckUniqueJob(jobName string, jobGroup string) string

	// SelectJobByInvokeTarget 通过调用目标字符串查询调度任务信息
	SelectJobByInvokeTarget(invokeTarget string) model.SysJob

	// InsertJob 新增调度任务信息
	InsertJob(sysJob model.SysJob) string

	// UpdateJob 修改调度任务信息
	UpdateJob(sysJob model.SysJob) int

	// DeleteJobByIds 批量删除调度任务信息
	DeleteJobByIds(jobIds []string) error
}
