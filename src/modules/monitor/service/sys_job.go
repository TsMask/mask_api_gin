package service

import (
	"mask_api_gin/src/modules/monitor/model"
)

// ISysJob 调度任务 服务层接口
type ISysJob interface {
	// FindByPage 分页查询
	FindByPage(query map[string]any) map[string]any

	// Find 查询
	Find(sysJob model.SysJob) []model.SysJob

	// FindById 通过ID查询
	FindById(jobId string) model.SysJob

	// Insert 新增调度任务信息
	Insert(sysJob model.SysJob) string

	// Update 修改
	Update(sysJob model.SysJob) int64

	// DeleteByIds 批量删除
	DeleteByIds(jobIds []string) (int64, error)

	// CheckUniqueByJobName 校验调度任务名称和组是否唯一
	CheckUniqueByJobName(jobName, jobGroup, jobId string) bool

	// Run 立即运行一次调度任务
	Run(sysJob model.SysJob) bool

	// Reset 重置初始调度任务
	Reset()
}
