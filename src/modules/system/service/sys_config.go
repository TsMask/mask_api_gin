package service

import "mask_api_gin/src/modules/system/model"

// ISysConfigService 参数配置 服务层接口
type ISysConfigService interface {
	// FindByPage 分页查询列表数据
	FindByPage(query map[string]any) map[string]any

	// FindById 通过ID查询信息
	FindById(configId string) model.SysConfig

	// Insert 新增信息
	Insert(sysConfig model.SysConfig) string

	// Update 修改信息
	Update(sysConfig model.SysConfig) int64

	// DeleteByIds 批量删除信息
	DeleteByIds(configIds []string) (int64, error)

	// FindValueByKey 通过参数键名查询参数值
	FindValueByKey(configKey string) string

	// CheckUniqueByKey 检查参数键名是否唯一
	CheckUniqueByKey(configKey, configId string) bool

	// CacheLoad 加载参数缓存数据 传入*查询全部
	CacheLoad(configKey string)

	// CacheClean 清空参数缓存数据 传入*清除全部
	CacheClean(configKey string) bool
}
