package repository

import (
	"mask_api_gin/src/modules/monitor/model"
)

// ISysJobRepository 调度任务 数据层接口
type ISysJobRepository interface {
	// SelectByPage 分页查询集合
	SelectByPage(query map[string]any) map[string]any

	// Select 查询集合
	Select(sysJob model.SysJob) []model.SysJob

	// SelectByIds 通过ID查询信息
	SelectByIds(jobIds []string) []model.SysJob

	// Insert 新增信息
	Insert(sysJob model.SysJob) string

	// Update 修改信息
	Update(sysJob model.SysJob) int64

	// DeleteByIds 批量删除信息
	DeleteByIds(jobIds []string) int64

	// CheckUniqueJob 校验信息是否唯一
	CheckUniqueJob(sysJob model.SysJob) string
}
