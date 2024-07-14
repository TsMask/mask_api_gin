package service

import "mask_api_gin/src/modules/system/model"

// ISysDictTypeService 字典类型 服务层接口
type ISysDictTypeService interface {
	// FindByPage 分页查询列表数据
	FindByPage(query map[string]any) map[string]any

	// Find 查询数据
	Find(sysDictType model.SysDictType) []model.SysDictType

	// FindById 通过ID查询信息
	FindById(dictId string) model.SysDictType

	// FindByType 根据字典类型查询信息
	FindByType(dictType string) model.SysDictType

	// Insert 新增信息
	Insert(sysDictType model.SysDictType) string

	// Update 修改信息
	Update(sysDictType model.SysDictType) int64

	// DeleteByIds 批量删除信息
	DeleteByIds(dictIds []string) (int64, error)

	// CheckUniqueByName 检查字典名称是否唯一
	CheckUniqueByName(dictName, dictId string) bool

	// CheckUniqueByType 检查字典类型是否唯一
	CheckUniqueByType(dictType, dictId string) bool

	// ResetCache 重置字典缓存数据
	ResetCache()

	// FindDataByType 获取字典数据缓存数据
	FindDataByType(dictType string) []model.SysDictData
}
