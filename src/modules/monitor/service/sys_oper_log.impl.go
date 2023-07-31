package service

import (
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
)

// 实例化服务层 SysOperLogImpl 结构体
var NewSysOperLogImpl = &SysOperLogImpl{
	sysOperLogService: repository.NewSysOperLogImpl,
}

// SysOperLogImpl 操作日志表 数据层处理
type SysOperLogImpl struct {
	// 操作日志信息
	sysOperLogService repository.ISysOperLog
}

// SelectOperLogPage 分页查询系统操作日志集合
func (r *SysOperLogImpl) SelectOperLogPage(query map[string]string) map[string]interface{} {
	return r.sysOperLogService.SelectOperLogPage(query)
}

// SelectOperLogList 查询系统操作日志集合
func (r *SysOperLogImpl) SelectOperLogList(sysOperLog model.SysOperLog) []model.SysOperLog {
	return r.sysOperLogService.SelectOperLogList(sysOperLog)
}

// SelectOperLogById 查询操作日志详细
func (r *SysOperLogImpl) SelectOperLogById(operId string) model.SysOperLog {
	return r.sysOperLogService.SelectOperLogById(operId)
}

// InsertOperLog 新增操作日志
func (r *SysOperLogImpl) InsertOperLog(sysOperLog model.SysOperLog) string {
	return r.sysOperLogService.InsertOperLog(sysOperLog)
}

// DeleteOperLogByIds 批量删除系统操作日志
func (r *SysOperLogImpl) DeleteOperLogByIds(operIds []string) int64 {
	return r.sysOperLogService.DeleteOperLogByIds(operIds)
}

// CleanOperLog 清空操作日志
func (r *SysOperLogImpl) CleanOperLog() error {
	return r.sysOperLogService.CleanOperLog()
}
