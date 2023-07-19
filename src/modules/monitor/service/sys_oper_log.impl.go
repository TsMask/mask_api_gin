package service

import (
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
)

// SysOperLogImpl 操作日志表 数据层处理
var SysOperLogImpl = &sysOperLogImpl{
	sysOperLogService: repository.SysOperLogImpl,
}

type sysOperLogImpl struct {
	// 操作日志信息
	sysOperLogService repository.ISysOperLog
}

// SelectOperLogPage 分页查询系统操作日志集合
func (r *sysOperLogImpl) SelectOperLogPage(query map[string]string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectOperLogList 查询系统操作日志集合
func (r *sysOperLogImpl) SelectOperLogList(sysOperLog model.SysOperLog) []model.SysOperLog {
	return []model.SysOperLog{}
}

// InsertOperLog 新增操作日志
func (r *sysOperLogImpl) InsertOperLog(sysOperLog model.SysOperLog) string {
	return r.sysOperLogService.InsertOperLog(sysOperLog)
}

// DeleteOperLogByIds 批量删除系统操作日志
func (r *sysOperLogImpl) DeleteOperLogByIds(operIds []string) int64 {
	return 0
}

// SelectOperLogById 查询操作日志详细
func (r *sysOperLogImpl) SelectOperLogById(operId string) model.SysOperLog {
	return model.SysOperLog{}
}

// CleanOperLog 清空操作日志
func (r *sysOperLogImpl) CleanOperLog() error {
	return nil
}