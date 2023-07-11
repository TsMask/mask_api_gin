package service

import (
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
)

// SysJobLogImpl 调度任务日志 业务层处理
var SysJobLogImpl = &sysJobLogImpl{
	sysJobLogRepository: repository.SysJobLogImpl,
}

type sysJobLogImpl struct {
	// 调度任务日志信息
	sysJobLogRepository repository.ISysJobLog
}

// SelectJobLogPage 分页查询调度任务日志集合
func (s *sysJobLogImpl) SelectJobLogPage(query map[string]string) map[string]interface{} {
	return s.sysJobLogRepository.SelectJobLogPage(query)
}

// SelectJobLogList 查询调度任务日志集合
func (s *sysJobLogImpl) SelectJobLogList(sysJobLog model.SysJobLog) []model.SysJobLog {
	return s.sysJobLogRepository.SelectJobLogList(sysJobLog)
}

// SelectJobLogById 通过调度ID查询调度任务日志信息
func (s *sysJobLogImpl) SelectJobLogById(jobLogId string) model.SysJobLog {
	return s.sysJobLogRepository.SelectJobLogById(jobLogId)
}

// InsertJobLog 新增调度任务日志信息
func (s *sysJobLogImpl) InsertJobLog(sysJobLog model.SysJobLog) string {
	return s.sysJobLogRepository.InsertJobLog(sysJobLog)
}

// DeleteJobLogByIds 批量删除调度任务日志信息
func (s *sysJobLogImpl) DeleteJobLogByIds(jobLogIds []string) int64 {
	return s.sysJobLogRepository.DeleteJobLogByIds(jobLogIds)
}

// CleanJobLog 清空调度任务日志
func (s *sysJobLogImpl) CleanJobLog() error {
	return s.sysJobLogRepository.CleanJobLog()
}
