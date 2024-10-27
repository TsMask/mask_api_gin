package service

import (
	"fmt"
	constCacheKey "mask_api_gin/src/framework/constants/cache_key"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysConfig 实例化服务层
var NewSysConfig = &SysConfig{
	sysConfigRepository: repository.NewSysConfig,
}

// SysConfig 参数配置 服务层处理
type SysConfig struct {
	sysConfigRepository *repository.SysConfig // 参数配置表
}

// FindByPage 分页查询列表数据
func (s SysConfig) FindByPage(query map[string]any) ([]model.SysConfig, int64) {
	return s.sysConfigRepository.SelectByPage(query)
}

// FindById 通过ID查询信息
func (s SysConfig) FindById(configId string) model.SysConfig {
	if configId == "" {
		return model.SysConfig{}
	}
	configs := s.sysConfigRepository.SelectByIds([]string{configId})
	if len(configs) > 0 {
		return configs[0]
	}
	return model.SysConfig{}
}

// Insert 新增信息
func (s SysConfig) Insert(sysConfig model.SysConfig) string {
	if configId := s.sysConfigRepository.Insert(sysConfig); configId != "" {
		s.CacheLoad(sysConfig.ConfigKey)
	}
	return ""
}

// Update 修改信息
func (s SysConfig) Update(sysConfig model.SysConfig) int64 {
	if rows := s.sysConfigRepository.Update(sysConfig); rows > 0 {
		s.CacheLoad(sysConfig.ConfigKey)
	}
	return 0
}

// DeleteByIds 批量删除信息
func (s SysConfig) DeleteByIds(configIds []string) (int64, error) {
	// 检查是否存在
	configs := s.sysConfigRepository.SelectByIds(configIds)
	if len(configs) <= 0 {
		return 0, fmt.Errorf("没有权限访问参数配置数据！")
	}
	for _, config := range configs {
		// 检查是否为内置参数
		if config.ConfigType == "Y" {
			return 0, fmt.Errorf("%s 配置参数属于内置参数，禁止删除！", config.ConfigId)
		}
		// 清除缓存
		s.CacheClean(config.ConfigKey)
	}
	if len(configs) == len(configIds) {
		return s.sysConfigRepository.DeleteByIds(configIds), nil
	}
	return 0, fmt.Errorf("删除参数配置信息失败！")
}

// FindValueByKey 通过参数键名查询参数值
func (s SysConfig) FindValueByKey(configKey string) string {
	cacheKey := s.getCacheKey(configKey)
	// 从缓存中读取
	if cacheValue, err := redis.Get("", cacheKey); cacheValue != "" || err != nil {
		return cacheValue
	}
	// 无缓存时读取数据放入缓存中
	if configValue := s.sysConfigRepository.SelectValueByKey(configKey); configValue != "" {
		_ = redis.Set("", cacheKey, configValue)
		return configValue
	}
	return ""
}

// CheckUniqueByKey 检查参数键名是否唯一
func (s SysConfig) CheckUniqueByKey(configKey, configId string) bool {
	uniqueId := s.sysConfigRepository.CheckUnique(model.SysConfig{
		ConfigKey: configKey,
	})
	if uniqueId == configId {
		return true
	}
	return uniqueId == ""
}

// getCacheKey 组装缓存key
func (s SysConfig) getCacheKey(configKey string) string {
	return constCacheKey.SYS_CONFIG_KEY + configKey
}

// CacheLoad 加载参数缓存数据 传入*查询全部
func (s SysConfig) CacheLoad(configKey string) {
	// 查询全部参数
	if configKey == "*" || configKey == "" {
		sysConfigs := s.sysConfigRepository.Select(model.SysConfig{})
		for _, v := range sysConfigs {
			key := s.getCacheKey(v.ConfigKey)
			_ = redis.Del("", key)
			_ = redis.Set("", key, v.ConfigValue)
		}
		return
	}
	// 指定参数
	cacheValue := s.sysConfigRepository.SelectValueByKey(configKey)
	if cacheValue != "" {
		key := s.getCacheKey(configKey)
		_ = redis.Del("", key)
		_ = redis.Set("", key, cacheValue)
	}
}

// CacheClean 清空参数缓存数据 传入*清除全部
func (s SysConfig) CacheClean(configKey string) bool {
	key := s.getCacheKey(configKey)
	keys, err := redis.GetKeys("", key)
	if err != nil {
		return false
	}
	return redis.DelKeys("", keys) == nil
}
