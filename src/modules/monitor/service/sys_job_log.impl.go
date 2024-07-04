package service

import (
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
)

// NewSysJobLogService 服务层实例化
var NewSysJobLogService = &SysJobLogServiceImpl{
	sysJobLogRepository: repository.NewSysJobLogRepository,
}

// SysJobLogServiceImpl 调度任务日志 服务层处理
type SysJobLogServiceImpl struct {
	// 调度任务日志数据信息
	sysJobLogRepository repository.ISysJobLogRepository
}

// FindByPage 分页查询
func (s *SysJobLogServiceImpl) FindByPage(query map[string]any) map[string]any {
	return s.sysJobLogRepository.SelectByPage(query)
}

// Find 查询
func (s *SysJobLogServiceImpl) Find(sysJobLog model.SysJobLog) []model.SysJobLog {
	return s.sysJobLogRepository.Select(sysJobLog)
}

// FindById 通过ID查询
func (s *SysJobLogServiceImpl) FindById(jobLogId string) model.SysJobLog {
	return s.sysJobLogRepository.SelectById(jobLogId)
}

// RemoveByIds 批量删除
func (s *SysJobLogServiceImpl) RemoveByIds(jobLogIds []string) int64 {
	return s.sysJobLogRepository.DeleteByIds(jobLogIds)
}

// Clean 清空
func (s *SysJobLogServiceImpl) Clean() error {
	return s.sysJobLogRepository.Clean()
}
