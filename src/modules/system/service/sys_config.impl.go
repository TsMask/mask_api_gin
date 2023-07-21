package service

import (
	"errors"
	"mask_api_gin/src/framework/constants/cachekey"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// 参数配置 服务层实现
var SysConfigImpl = &sysConfigImpl{
	sysConfigRepository: repository.SysConfigImpl,
}

type sysConfigImpl struct {
	// 参数配置表
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
	// 从缓存中读取
	cacheValue := redis.Get(cacheKey)
	if cacheValue != "" {
		return cacheValue
	}
	// 无缓存时读取数据放入缓存中
	configValue := r.sysConfigRepository.SelectConfigValueByKey(configKey)
	if configValue != "" {
		redis.Set(cacheKey, configValue)
		return configValue
	}
	return ""
}

// SelectConfigById 通过配置ID查询参数配置信息
func (r *sysConfigImpl) SelectConfigById(configId string) model.SysConfig {
	if configId == "" {
		return model.SysConfig{}
	}
	configs := r.sysConfigRepository.SelectConfigByIds([]string{configId})
	if len(configs) > 0 {
		return configs[0]
	}
	return model.SysConfig{}
}

// CheckUniqueConfigKey 校验参数键名是否唯一
func (r *sysConfigImpl) CheckUniqueConfigKey(configKey, configId string) bool {
	uniqueId := r.sysConfigRepository.CheckUniqueConfig(model.SysConfig{
		ConfigKey: configKey,
	})
	if uniqueId == configId {
		return true
	}
	return uniqueId == ""
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
func (r *sysConfigImpl) UpdateConfig(sysConfig model.SysConfig) int64 {
	rows := r.sysConfigRepository.UpdateConfig(sysConfig)
	if rows > 0 {
		r.loadingConfigCache(sysConfig.ConfigKey)
	}
	return rows
}

// DeleteConfigByIds 批量删除参数配置信息
func (r *sysConfigImpl) DeleteConfigByIds(configIds []string) (int64, error) {
	// 检查是否存在
	configs := r.sysConfigRepository.SelectConfigByIds(configIds)
	if len(configs) <= 0 {
		return 0, errors.New("没有权限访问参数配置数据！")
	}
	for _, config := range configs {
		// 检查是否为内置参数
		if config.ConfigType == "Y" {
			return 0, errors.New(config.ConfigID + " 配置参数属于内置参数，禁止删除！")
		}
		// 清除缓存
		r.clearConfigCache(config.ConfigKey)
	}
	if len(configs) == len(configIds) {
		rows := r.sysConfigRepository.DeleteConfigByIds(configIds)
		return rows, nil
	}
	return 0, errors.New("删除参数配置信息失败！")
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
	// 查询全部参数
	if configKey == "*" {
		sysConfigs := r.SelectConfigList(model.SysConfig{})
		for _, v := range sysConfigs {
			key := r.getCacheKey(v.ConfigKey)
			redis.Del(key)
			redis.Set(key, v.ConfigValue)
		}
		return
	}
	// 指定参数
	if configKey != "" {
		cacheValue := r.sysConfigRepository.SelectConfigValueByKey(configKey)
		if cacheValue != "" {
			key := r.getCacheKey(configKey)
			redis.Del(key)
			redis.Set(key, cacheValue)
		}
		return
	}
}

// clearConfigCache 清空参数缓存数据
func (r *sysConfigImpl) clearConfigCache(configKey string) bool {
	key := r.getCacheKey(configKey)
	keys := redis.GetKeys(key)
	return redis.DelKeys(keys)
}
