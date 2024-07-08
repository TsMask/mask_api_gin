package repository

import "mask_api_gin/src/modules/system/model"

// ISysLogOperateRepository 操作日志表 数据层接口
type ISysLogOperateRepository interface {
	// SelectByPage 分页查询集合
	SelectByPage(query map[string]any) map[string]any

	// Select 查询集合
	Select(sysLogOperate model.SysLogOperate) []model.SysLogOperate

	// SelectById 通过ID查询信息
	SelectById(operaId string) model.SysLogOperate

	// Insert 新增信息
	Insert(sysLogOperate model.SysLogOperate) string

	// DeleteByIds 批量删除信息
	DeleteByIds(operaIds []string) int64

	// Clean 清空信息
	Clean() error
}
