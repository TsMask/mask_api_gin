package repository

import "mask_api_gin/src/modules/system/model"

// ISysConfig 参数配置表 数据层接口
type ISysConfig interface {
	// SelectDictDataPage 分页查询参数配置列表数据
	SelectConfigPage(query map[string]any) map[string]any

	// SelectConfigList 查询参数配置列表
	SelectConfigList(sysConfig model.SysConfig) []model.SysConfig

	// SelectConfigValueByKey 通过参数键名查询参数键值
	SelectConfigValueByKey(configKey string) string

	// SelectConfigByIds 通过配置ID查询参数配置信息
	SelectConfigByIds(configIds []string) []model.SysConfig

	// CheckUniqueConfig 校验配置参数是否唯一
	CheckUniqueConfig(sysConfig model.SysConfig) string

	// InsertConfig 新增参数配置
	InsertConfig(sysConfig model.SysConfig) string

	// UpdateConfig 修改参数配置
	UpdateConfig(sysConfig model.SysConfig) int64

	// DeleteConfigByIds 批量删除参数配置信息
	DeleteConfigByIds(configIds []string) int64
}
