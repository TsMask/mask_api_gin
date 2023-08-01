package service

import (
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
)

// 实例化服务层 SysJobLogImpl 结构体
var NewSysJobLogImpl = &SysJobLogImpl{
	sysJobLogRepository: repository.NewSysJobLogImpl,
}

// SysJobLogImpl 调度任务日志 服务层处理
type SysJobLogImpl struct {
	// 调度任务日志数据信息
	sysJobLogRepository repository.ISysJobLog
}

// SelectJobLogPage 分页查询调度任务日志集合
func (s *SysJobLogImpl) SelectJobLogPage(query map[string]any) map[string]any {
	return s.sysJobLogRepository.SelectJobLogPage(query)
}

// SelectJobLogList 查询调度任务日志集合
func (s *SysJobLogImpl) SelectJobLogList(sysJobLog model.SysJobLog) []model.SysJobLog {
	return s.sysJobLogRepository.SelectJobLogList(sysJobLog)
}

// SelectJobLogById 通过调度ID查询调度任务日志信息
func (s *SysJobLogImpl) SelectJobLogById(jobLogId string) model.SysJobLog {
	return s.sysJobLogRepository.SelectJobLogById(jobLogId)
}

// DeleteJobLogByIds 批量删除调度任务日志信息
func (s *SysJobLogImpl) DeleteJobLogByIds(jobLogIds []string) int64 {
	return s.sysJobLogRepository.DeleteJobLogByIds(jobLogIds)
}

// CleanJobLog 清空调度任务日志
func (s *SysJobLogImpl) CleanJobLog() error {
	return s.sysJobLogRepository.CleanJobLog()
}
