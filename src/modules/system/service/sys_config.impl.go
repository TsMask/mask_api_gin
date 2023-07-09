package service

import (
	"mask_api_gin/src/framework/constants/cachekey"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// 参数配置 服务层实现
var SysConfigImpl = &sysConfigImpl{
	sysConfigRepository: repository.SysConfigImpl,
}

type sysConfigImpl struct {
	sysConfigRepository repository.ISysConfig
}

// SelectDictDataPage 分页查询参数配置列表数据
func (r *sysConfigImpl) SelectConfigPage(query map[string]string) map[string]interface{} {
	return r.sysConfigRepository.SelectConfigPage(query)
}

// SelectConfigList 查询参数配置列表
func (r *sysConfigImpl) SelectConfigList(sysConfig model.SysConfig) []model.SysConfig {
	return r.sysConfigRepository.SelectConfigList(sysConfig)
}

// SelectConfigValueByKey 通过参数键名查询参数键值
func (r *sysConfigImpl) SelectConfigValueByKey(configKey string) string {
	cacheKey := r.getCacheKey(configKey)

	return r.sysConfigRepository.SelectConfigValueByKey(cacheKey)
}

// SelectConfigById 通过配置ID查询参数配置信息
func (r *sysConfigImpl) SelectConfigById(configId string) model.SysConfig {
	return r.sysConfigRepository.SelectConfigById(configId)
}

// CheckUniqueConfigKey 校验参数键名是否唯一
func (r *sysConfigImpl) CheckUniqueConfigKey(sysConfig model.SysConfig) bool {
	configId := r.sysConfigRepository.CheckUniqueConfigKey(sysConfig.ConfigKey)
	if configId == "" {
		return true
	}
	// 与查询得到的一致
	if configId == sysConfig.ConfigID {
		return true
	}
	return false
}

// InsertConfig 新增参数配置
func (r *sysConfigImpl) InsertConfig(sysConfig model.SysConfig) string {
	configId := r.sysConfigRepository.InsertConfig(sysConfig)
	if configId != "" {
		r.loadingConfigCache(sysConfig.ConfigKey)
	}
	return configId
}

// UpdateConfig 修改参数配置
func (r *sysConfigImpl) UpdateConfig(sysConfig model.SysConfig) int {
	rows := r.sysConfigRepository.UpdateConfig(sysConfig)
	if rows > 0 {
		r.loadingConfigCache(sysConfig.ConfigKey)
	}
	return rows
}

// DeleteConfigByIds 批量删除参数配置信息
func (r *sysConfigImpl) DeleteConfigByIds(configIds []string) int {
	for _, configId := range configIds {
		// 检查是否存在
		config := r.sysConfigRepository.SelectConfigById(configId)
		if config.ConfigID != configId {
			return 0
		}
		// 检查是否为内置参数
		if config.ConfigType == "Y" {
			return 0
		}
		// 清除缓存
		r.clearConfigCache(config.ConfigKey)
	}
	return r.sysConfigRepository.DeleteConfigByIds(configIds)
}

// ResetConfigCache 重置参数缓存数据
func (r *sysConfigImpl) ResetConfigCache() {
	r.clearConfigCache("*")
	r.loadingConfigCache("*")
}

// getCacheKey 组装缓存key
func (r *sysConfigImpl) getCacheKey(configKey string) string {
	return cachekey.SYS_CONFIG_KEY + configKey
}

// loadingConfigCache 加载参数缓存数据
func (r *sysConfigImpl) loadingConfigCache(configKey string) {
}

// clearConfigCache 清空参数缓存数据
func (r *sysConfigImpl) clearConfigCache(configKey string) int {
	return 0
}
