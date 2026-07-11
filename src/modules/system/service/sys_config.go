package service

import (
	"fmt"
	"path/filepath"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/database/redis"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
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
func (s SysConfig) FindByPage(query map[string]string) ([]model.SysConfig, int64) {
	return s.sysConfigRepository.SelectByPage(query)
}

// FindById 通过ID查询信息
func (s SysConfig) FindById(id int64) model.SysConfig {
	if id <= 0 {
		return model.SysConfig{}
	}
	configs := s.sysConfigRepository.SelectByIds([]int64{id})
	if len(configs) > 0 {
		return configs[0]
	}
	return model.SysConfig{}
}

// Insert 新增信息
func (s SysConfig) Insert(sysConfig model.SysConfig) int64 {
	id := s.sysConfigRepository.Insert(sysConfig)
	if id > 0 {
		s.CacheLoad(sysConfig.ConfigKey)
	}
	return id
}

// Update 修改信息
func (s SysConfig) Update(sysConfig model.SysConfig) int64 {
	rows := s.sysConfigRepository.Update(sysConfig)
	if rows > 0 {
		s.CacheLoad(sysConfig.ConfigKey)
	}
	return rows
}

// DeleteByIds 批量删除信息
func (s SysConfig) DeleteByIds(ids []int64) (int64, error) {
	// 检查是否存在
	configs := s.sysConfigRepository.SelectByIds(ids)
	if len(configs) <= 0 {
		return 0, fmt.Errorf("没有权限访问参数配置数据！")
	}
	for _, config := range configs {
		// 清除缓存
		s.CacheClean(config.ConfigKey)
	}
	if len(configs) == len(ids) {
		return s.sysConfigRepository.DeleteByIds(ids), nil
	}
	return 0, fmt.Errorf("删除参数配置信息失败！")
}

// CheckUniqueByKey 检查参数键名是否唯一
func (s SysConfig) CheckUniqueByKey(configKey string, id int64) bool {
	uniqueId := s.sysConfigRepository.CheckUnique(model.SysConfig{
		ConfigKey: configKey,
	})
	if uniqueId == id {
		return true
	}
	return uniqueId == 0
}

// FindValueByKey 通过参数键名查询参数值
func (s SysConfig) FindValueByKey(configKey string) string {
	cacheKey := constants.CACHE_SYS_CONFIG + ":" + configKey
	// 从缓存中读取
	if cacheValue, err := redis.Get("", cacheKey); cacheValue != "" || err != nil {
		return cacheValue
	}
	// 无缓存时读取数据放入缓存中
	configValue := s.sysConfigRepository.SelectValueByKey(configKey)
	if configValue != "" {
		_ = redis.Set("", cacheKey, configValue, 0)
		return configValue
	}
	return ""
}

// CacheLoad 加载参数缓存数据 传入*查询全部
func (s SysConfig) CacheLoad(configKey string) {
	// 查询全部参数
	if configKey == "*" || configKey == "" {
		sysConfigs := s.sysConfigRepository.Select(model.SysConfig{})
		for _, v := range sysConfigs {
			key := constants.CACHE_SYS_CONFIG + ":" + v.ConfigKey
			_ = redis.Del("", key)
			_ = redis.Set("", key, v.ConfigValue, 0)
		}
		return
	}
	// 指定参数
	cacheValue := s.sysConfigRepository.SelectValueByKey(configKey)
	if cacheValue != "" {
		key := constants.CACHE_SYS_CONFIG + ":" + configKey
		_ = redis.Del("", key)
		_ = redis.Set("", key, cacheValue, 0)
	}
}

// CacheClean 清空参数缓存数据 传入*清除全部
func (s SysConfig) CacheClean(configKey string) bool {
	key := constants.CACHE_SYS_CONFIG + ":" + configKey
	keys, err := redis.GetKeys("", key)
	if err != nil {
		return false
	}
	return redis.DelKeys("", keys) == nil
}

// ExportData 导出数据表格
func (s SysConfig) ExportData(rows []model.SysConfig, fileName string) (string, error) {
	// 第一行表头标题
	data := [][]string{
		{"参数编号", "参数名称", "参数键名", "参数键值", "参数类型"},
	}
	// 从第二行开始的数据
	for _, row := range rows {
		data = append(data, []string{
			fmt.Sprintf("%d", row.ID),
			row.ConfigName,
			row.ConfigKey,
			row.ConfigValue,
			row.ConfigValueType,
		})
	}

	// 导出文件路径
	_, dir := file.UploadResourceDir()
	datePath := date.ParseDatePath()
	filePath := filepath.Join(dir, constants.UPLOAD_EXPORT, datePath, fileName)

	// 导出数据表格
	return filePath, file.CSVWrite(filePath, data, false)
}
