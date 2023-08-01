package service

import (
	"errors"
	"mask_api_gin/src/framework/constants/cachekey"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// 实例化服务层 SysConfigImpl 结构体
var NewSysConfigImpl = &SysConfigImpl{
	sysConfigRepository: repository.NewSysConfigImpl,
}

// SysConfigImpl 参数配置 服务层处理
type SysConfigImpl struct {
	// 参数配置表
	sysConfigRepository repository.ISysConfig
}

// SelectDictDataPage 分页查询参数配置列表数据
func (r *SysConfigImpl) SelectConfigPage(query map[string]any) map[string]any {
	return r.sysConfigRepository.SelectConfigPage(query)
}

// SelectConfigList 查询参数配置列表
func (r *SysConfigImpl) SelectConfigList(sysConfig model.SysConfig) []model.SysConfig {
	return r.sysConfigRepository.SelectConfigList(sysConfig)
}

// SelectConfigValueByKey 通过参数键名查询参数键值
func (r *SysConfigImpl) SelectConfigValueByKey(configKey string) string {
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
func (r *SysConfigImpl) SelectConfigById(configId string) model.SysConfig {
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
func (r *SysConfigImpl) CheckUniqueConfigKey(configKey, configId string) bool {
	uniqueId := r.sysConfigRepository.CheckUniqueConfig(model.SysConfig{
		ConfigKey: configKey,
	})
	if uniqueId == configId {
		return true
	}
	return uniqueId == ""
}

// InsertConfig 新增参数配置
func (r *SysConfigImpl) InsertConfig(sysConfig model.SysConfig) string {
	configId := r.sysConfigRepository.InsertConfig(sysConfig)
	if configId != "" {
		r.loadingConfigCache(sysConfig.ConfigKey)
	}
	return configId
}

// UpdateConfig 修改参数配置
func (r *SysConfigImpl) UpdateConfig(sysConfig model.SysConfig) int64 {
	rows := r.sysConfigRepository.UpdateConfig(sysConfig)
	if rows > 0 {
		r.loadingConfigCache(sysConfig.ConfigKey)
	}
	return rows
}

// DeleteConfigByIds 批量删除参数配置信息
func (r *SysConfigImpl) DeleteConfigByIds(configIds []string) (int64, error) {
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
func (r *SysConfigImpl) ResetConfigCache() {
	r.clearConfigCache("*")
	r.loadingConfigCache("*")
}

// getCacheKey 组装缓存key
func (r *SysConfigImpl) getCacheKey(configKey string) string {
	return cachekey.SYS_CONFIG_KEY + configKey
}

// loadingConfigCache 加载参数缓存数据
func (r *SysConfigImpl) loadingConfigCache(configKey string) {
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
func (r *SysConfigImpl) clearConfigCache(configKey string) bool {
	key := r.getCacheKey(configKey)
	keys := redis.GetKeys(key)
	return redis.DelKeys(keys)
}
