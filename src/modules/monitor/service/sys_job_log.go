package service

import (
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
)

// NewSysJobLog 服务层实例化
var NewSysJobLog = &SysJobLog{
	sysJobLogRepository: repository.NewSysJobLog,
}

// SysJobLog 调度任务日志 服务层处理
type SysJobLog struct {
	sysJobLogRepository *repository.SysJobLog // 调度任务日志数据信息
}

// FindByPage 分页查询
func (s *SysJobLog) FindByPage(query map[string]any) map[string]any {
	return s.sysJobLogRepository.SelectByPage(query)
}

// Find 查询
func (s *SysJobLog) Find(sysJobLog model.SysJobLog) []model.SysJobLog {
	return s.sysJobLogRepository.Select(sysJobLog)
}

// FindById 通过ID查询
func (s *SysJobLog) FindById(jobLogId string) model.SysJobLog {
	return s.sysJobLogRepository.SelectById(jobLogId)
}

// RemoveByIds 批量删除
func (s *SysJobLog) RemoveByIds(jobLogIds []string) int64 {
	return s.sysJobLogRepository.DeleteByIds(jobLogIds)
}

// Clean 清空
func (s *SysJobLog) Clean() error {
	return s.sysJobLogRepository.Clean()
}
