package repository

import (
	"mask_api_gin/src/modules/monitor/model"
)

// ISysJobLogRepository 调度任务日志表 数据层接口
type ISysJobLogRepository interface {
	// SelectByPage 分页查询集合
	SelectByPage(query map[string]any) map[string]any

	// Select 查询集合
	Select(sysJobLog model.SysJobLog) []model.SysJobLog

	// SelectById 通过ID查询信息
	SelectById(jobLogId string) model.SysJobLog

	// Insert 新增信息
	Insert(sysJobLog model.SysJobLog) string

	// DeleteByIds 批量删除信息
	DeleteByIds(jobLogId []string) int64

	// Clean 清空集合数据
	Clean() error
}
