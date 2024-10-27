package service

import (
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysLogOperate 实例化服务层
var NewSysLogOperate = &SysLogOperate{
	SysLogOperate: repository.NewSysLogOperate,
}

// SysLogOperate 操作日志表 服务层处理
type SysLogOperate struct {
	SysLogOperate *repository.SysLogOperate // 操作日志信息
}

// FindByPage 分页查询列表数据
func (s SysLogOperate) FindByPage(query map[string]any) map[string]any {
	return s.SysLogOperate.SelectByPage(query)
}

// Find 查询数据
func (s SysLogOperate) Find(SysLogOperate model.SysLogOperate) []model.SysLogOperate {
	return s.SysLogOperate.Select(SysLogOperate)
}

// FindById 根据ID查询信息
func (s SysLogOperate) FindById(operaId string) model.SysLogOperate {
	return s.SysLogOperate.SelectById(operaId)
}

// Insert 新增信息
func (s SysLogOperate) Insert(SysLogOperate model.SysLogOperate) string {
	return s.SysLogOperate.Insert(SysLogOperate)
}

// DeleteById 删除信息
func (s SysLogOperate) DeleteById(operaIds []string) int64 {
	return s.SysLogOperate.DeleteByIds(operaIds)
}

// Clean 清空操作日志
func (s SysLogOperate) Clean() error {
	return s.SysLogOperate.Clean()
}
