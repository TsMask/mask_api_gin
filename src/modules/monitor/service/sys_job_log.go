package service

import (
	"mask_api_gin/src/modules/monitor/model"
)

// ISysJobLogService 调度任务日志 服务层接口
type ISysJobLogService interface {
	// FindByPage 分页查询
	FindByPage(query map[string]any) map[string]any

	// Find 查询
	Find(sysJobLog model.SysJobLog) []model.SysJobLog

	// FindById 通过ID查询
	FindById(jobLogId string) model.SysJobLog

	// RemoveByIds 批量删除
	RemoveByIds(jobLogIds []string) int64

	// Clean 清空
	Clean() error
}
