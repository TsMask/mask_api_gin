package service

import "mask_api_gin/src/modules/system/model"

// ISysLogOperateService 操作日志表 服务层接口
type ISysLogOperateService interface {
	// FindByPage 分页查询列表数据
	FindByPage(query map[string]any) map[string]any

	// Find 查询数据
	Find(sysLogOperate model.SysLogOperate) []model.SysLogOperate

	// FindById 根据ID查询信息
	FindById(operaId string) model.SysLogOperate

	// Insert 新增信息
	Insert(sysLogOperate model.SysLogOperate) string

	// DeleteById 删除信息
	DeleteById(operaIds []string) int64

	// Clean 清空操作日志
	Clean() error
}
