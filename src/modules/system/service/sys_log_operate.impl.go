package service

import (
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// 实例化服务层 SysLogOperateImpl 结构体
var NewSysLogOperateImpl = &SysLogOperateImpl{
	SysLogOperateService: repository.NewSysLogOperateImpl,
}

// SysLogOperateImpl 操作日志表 服务层处理
type SysLogOperateImpl struct {
	// 操作日志信息
	SysLogOperateService repository.ISysLogOperate
}

// SelectSysLogOperatePage 分页查询系统操作日志集合
func (r *SysLogOperateImpl) SelectSysLogOperatePage(query map[string]any) map[string]any {
	return r.SysLogOperateService.SelectSysLogOperatePage(query)
}

// SelectSysLogOperateList 查询系统操作日志集合
func (r *SysLogOperateImpl) SelectSysLogOperateList(SysLogOperate model.SysLogOperate) []model.SysLogOperate {
	return r.SysLogOperateService.SelectSysLogOperateList(SysLogOperate)
}

// SelectSysLogOperateById 查询操作日志详细
func (r *SysLogOperateImpl) SelectSysLogOperateById(operId string) model.SysLogOperate {
	return r.SysLogOperateService.SelectSysLogOperateById(operId)
}

// InsertSysLogOperate 新增操作日志
func (r *SysLogOperateImpl) InsertSysLogOperate(SysLogOperate model.SysLogOperate) string {
	return r.SysLogOperateService.InsertSysLogOperate(SysLogOperate)
}

// DeleteSysLogOperateByIds 批量删除系统操作日志
func (r *SysLogOperateImpl) DeleteSysLogOperateByIds(operIds []string) int64 {
	return r.SysLogOperateService.DeleteSysLogOperateByIds(operIds)
}

// CleanSysLogOperate 清空操作日志
func (r *SysLogOperateImpl) CleanSysLogOperate() error {
	return r.SysLogOperateService.CleanSysLogOperate()
}
