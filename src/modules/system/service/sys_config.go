package service

import "mask_api_gin/src/modules/system/model"

// ISysConfig 参数配置 服务层接口
type ISysConfig interface {
	// SelectDictDataPage 分页查询参数配置列表数据
	SelectConfigPage(query map[string]string) map[string]interface{}

	// SelectConfigValueByKey 通过参数键名查询参数键值
	SelectConfigValueByKey(configKey string) string

	// SelectConfigById 通过配置ID查询参数配置信息
	SelectConfigById(configId string) model.SysConfig

	// CheckUniqueConfigKey 校验参数键名是否唯一
	CheckUniqueConfigKey(configKey, configId string) bool

	// InsertConfig 新增参数配置
	InsertConfig(sysConfig model.SysConfig) string

	// UpdateConfig 修改参数配置
	UpdateConfig(sysConfig model.SysConfig) int64

	// DeleteConfigByIds 批量删除参数配置信息
	DeleteConfigByIds(configIds []string) (int64, error)

	// ResetConfigCache 重置参数缓存数据
	ResetConfigCache()
}
