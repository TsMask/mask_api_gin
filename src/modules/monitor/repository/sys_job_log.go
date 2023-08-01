package repository

import (
	"mask_api_gin/src/modules/monitor/model"
)

// 调度任务日志表 数据层接口
type ISysJobLog interface {
	// 分页查询调度任务日志集合
	SelectJobLogPage(query map[string]any) map[string]any

	// 查询调度任务日志集合
	SelectJobLogList(sysJobLog model.SysJobLog) []model.SysJobLog

	// 通过调度ID查询调度任务日志信息
	SelectJobLogById(jobLogId string) model.SysJobLog

	// 新增调度任务日志信息
	InsertJobLog(sysJobLog model.SysJobLog) string

	// 批量删除调度任务日志信息
	DeleteJobLogByIds(jobLogId []string) int64

	// 清空调度任务日志
	CleanJobLog() error
}
