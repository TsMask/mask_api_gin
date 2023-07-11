package service

import "mask_api_gin/src/modules/monitor/model"

// ISysOperLog 操作日志表 服务层接口
type ISysOperLog interface {
	// SelectOperLogPage 分页查询系统操作日志集合
	SelectOperLogPage(query map[string]string) map[string]interface{}

	// SelectOperLogList 查询系统操作日志集合
	SelectOperLogList(sysOperLog model.SysOperLog) []model.SysOperLog

	// SelectOperLogById 查询操作日志详细
	SelectOperLogById(operId string) model.SysOperLog

	// InsertOperLog 新增操作日志
	InsertOperLog(sysOperLog model.SysOperLog) string

	// DeleteOperLogByIds 批量删除系统操作日志
	DeleteOperLogByIds(operIds []string) int64

	// CleanOperLog 清空操作日志
	CleanOperLog() error
}
