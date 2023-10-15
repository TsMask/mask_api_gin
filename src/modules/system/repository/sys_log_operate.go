package repository

import "mask_api_gin/src/modules/system/model"

// ISysLogOperate 操作日志表 数据层接口
type ISysLogOperate interface {
	// SelectSysLogOperatePage 分页查询系统操作日志集合
	SelectSysLogOperatePage(query map[string]any) map[string]any

	// SelectSysLogOperateList 查询系统操作日志集合
	SelectSysLogOperateList(sysLogOperate model.SysLogOperate) []model.SysLogOperate

	// SelectSysLogOperateById 查询操作日志详细
	SelectSysLogOperateById(operId string) model.SysLogOperate

	// InsertSysLogOperate 新增操作日志
	InsertSysLogOperate(sysLogOperate model.SysLogOperate) string

	// DeleteSysLogOperateByIds 批量删除系统操作日志
	DeleteSysLogOperateByIds(operIds []string) int64

	// CleanSysLogOperate 清空操作日志
	CleanSysLogOperate() error
}
