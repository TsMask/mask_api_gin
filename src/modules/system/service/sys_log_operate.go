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
func (s SysLogOperate) FindByPage(query map[string]string) ([]model.SysLogOperate, int64) {
	return s.SysLogOperate.SelectByPage(query)
}

// Insert 新增信息
func (s SysLogOperate) Insert(SysLogOperate model.SysLogOperate) string {
	return s.SysLogOperate.Insert(SysLogOperate)
}

// Clean 清空操作日志
func (s SysLogOperate) Clean() error {
	return s.SysLogOperate.Clean()
}
