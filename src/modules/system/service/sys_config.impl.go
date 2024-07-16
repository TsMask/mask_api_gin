package service

import (
	"fmt"
	constCacheKey "mask_api_gin/src/framework/constants/cache_key"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysConfig 实例化服务层
var NewSysConfig = &SysConfigService{
	sysConfigRepository: repository.NewSysConfig,
}

// SysConfigService 参数配置 服务层处理
type SysConfigService struct {
	sysConfigRepository repository.ISysConfigRepository // 参数配置表
}

// FindByPage 分页查询列表数据
func (r *SysConfigService) FindByPage(query map[string]any) map[string]any {
	return r.sysConfigRepository.SelectByPage(query)
}

// FindById 通过ID查询信息
func (r *SysConfigService) FindById(configId string) model.SysConfig {
	if configId == "" {
		return model.SysConfig{}
	}
	configs := r.sysConfigRepository.SelectByIds([]string{configId})
	if len(configs) > 0 {
		return configs[0]
	}
	return model.SysConfig{}
}

// Insert 新增信息
func (r *SysConfigService) Insert(sysConfig model.SysConfig) string {
	if configId := r.sysConfigRepository.Insert(sysConfig); configId != "" {
		r.CacheLoad(sysConfig.ConfigKey)
	}
	return ""
}

// Update 修改信息
func (r *SysConfigService) Update(sysConfig model.SysConfig) int64 {
	if rows := r.sysConfigRepository.Update(sysConfig); rows > 0 {
		r.CacheLoad(sysConfig.ConfigKey)
	}
	return 0
}

// DeleteByIds 批量删除信息
func (r *SysConfigService) DeleteByIds(configIds []string) (int64, error) {
	// 检查是否存在
	configs := r.sysConfigRepository.SelectByIds(configIds)
	if len(configs) <= 0 {
		return 0, fmt.Errorf("没有权限访问参数配置数据！")
	}
	for _, config := range configs {
		// 检查是否为内置参数
		if config.ConfigType == "Y" {
			return 0, fmt.Errorf("%s 配置参数属于内置参数，禁止删除！", config.ConfigID)
		}
		// 清除缓存
		r.CacheClean(config.ConfigKey)
	}
	if len(configs) == len(configIds) {
		return r.sysConfigRepository.DeleteByIds(configIds), nil
	}
	return 0, fmt.Errorf("删除参数配置信息失败！")
}

// FindValueByKey 通过参数键名查询参数值
func (r *SysConfigService) FindValueByKey(configKey string) string {
	cacheKey := r.getCacheKey(configKey)
	// 从缓存中读取
	if cacheValue, err := redis.Get("", cacheKey); cacheValue != "" || err != nil {
		return cacheValue
	}
	// 无缓存时读取数据放入缓存中
	if configValue := r.sysConfigRepository.SelectValueByKey(configKey); configValue != "" {
		_ = redis.Set("", cacheKey, configValue)
		return configValue
	}
	return ""
}

// CheckUniqueByKey 检查参数键名是否唯一
func (r *SysConfigService) CheckUniqueByKey(configKey, configId string) bool {
	uniqueId := r.sysConfigRepository.CheckUnique(model.SysConfig{
		ConfigKey: configKey,
	})
	if uniqueId == configId {
		return true
	}
	return uniqueId == ""
}

// getCacheKey 组装缓存key
func (r *SysConfigService) getCacheKey(configKey string) string {
	return constCacheKey.SysConfigKey + configKey
}

// CacheLoad 加载参数缓存数据 传入*查询全部
func (r *SysConfigService) CacheLoad(configKey string) {
	// 查询全部参数
	if configKey == "*" || configKey == "" {
		sysConfigs := r.sysConfigRepository.Select(model.SysConfig{})
		for _, v := range sysConfigs {
			key := r.getCacheKey(v.ConfigKey)
			_ = redis.Del("", key)
			_ = redis.Set("", key, v.ConfigValue)
		}
		return
	}
	// 指定参数
	cacheValue := r.sysConfigRepository.SelectValueByKey(configKey)
	if cacheValue != "" {
		key := r.getCacheKey(configKey)
		_ = redis.Del("", key)
		_ = redis.Set("", key, cacheValue)
	}
	return
}

// CacheClean 清空参数缓存数据 传入*清除全部
func (r *SysConfigService) CacheClean(configKey string) bool {
	key := r.getCacheKey(configKey)
	keys, err := redis.GetKeys("", key)
	if err != nil {
		return false
	}
	return redis.DelKeys("", keys) == nil
}
