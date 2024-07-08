package repository

import "mask_api_gin/src/modules/system/model"

// ISysConfigRepository  参数配置表 数据层接口
type ISysConfigRepository interface {
	// SelectByPage 分页查询集合
	SelectByPage(query map[string]any) map[string]any

	// Select 查询集合
	Select(sysConfig model.SysConfig) []model.SysConfig

	// SelectByIds 通过ID查询信息
	SelectByIds(configIds []string) []model.SysConfig

	// Insert 新增信息
	Insert(sysConfig model.SysConfig) string

	// Update 修改信息
	Update(sysConfig model.SysConfig) int64

	// DeleteByIds 批量删除信息
	DeleteByIds(configIds []string) int64

	// CheckUnique 检查信息是否唯一
	CheckUnique(sysConfig model.SysConfig) string

	// SelectValueByKey 通过Key查询Value
	SelectValueByKey(configKey string) string
}
