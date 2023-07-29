package service

import (
	"mask_api_gin/src/modules/monitor/model"
)

// ISysJobLog 调度任务日志 服务层接口
type ISysJobLog interface {
	// SelectJobLogPage 分页查询调度任务日志集合
	SelectJobLogPage(query map[string]string) map[string]interface{}

	// SelectJobLogList 查询调度任务日志集合
	SelectJobLogList(sysJobLog model.SysJobLog) []model.SysJobLog

	// SelectJobLogById 通过调度ID查询调度任务日志信息
	SelectJobLogById(jobLogId string) model.SysJobLog

	// DeleteJobLogByIds 批量删除调度任务日志信息
	DeleteJobLogByIds(jobLogIds []string) int64

	// CleanJobLog 清空调度任务日志
	CleanJobLog() error
}
