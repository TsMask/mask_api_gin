package service

import (
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysLogOperate 实例化服务层
var NewSysLogOperate = &SysLogOperateService{
	SysLogOperateService: repository.NewSysLogOperate,
}

// SysLogOperateService 操作日志表 服务层处理
type SysLogOperateService struct {
	SysLogOperateService repository.ISysLogOperateRepository // 操作日志信息
}

// FindByPage 分页查询列表数据
func (r *SysLogOperateService) FindByPage(query map[string]any) map[string]any {
	return r.SysLogOperateService.SelectByPage(query)
}

// Find 查询数据
func (r *SysLogOperateService) Find(SysLogOperate model.SysLogOperate) []model.SysLogOperate {
	return r.SysLogOperateService.Select(SysLogOperate)
}

// FindById 根据ID查询信息
func (r *SysLogOperateService) FindById(operaId string) model.SysLogOperate {
	return r.SysLogOperateService.SelectById(operaId)
}

// Insert 新增信息
func (r *SysLogOperateService) Insert(SysLogOperate model.SysLogOperate) string {
	return r.SysLogOperateService.Insert(SysLogOperate)
}

// DeleteById 删除信息
func (r *SysLogOperateService) DeleteById(operaIds []string) int64 {
	return r.SysLogOperateService.DeleteByIds(operaIds)
}

// Clean 清空操作日志
func (r *SysLogOperateService) Clean() error {
	return r.SysLogOperateService.Clean()
}
