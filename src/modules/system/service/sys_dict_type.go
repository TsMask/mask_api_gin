package service

import "mask_api_gin/src/modules/system/model"

// ISysDictType 字典类型 服务层接口
type ISysDictType interface {
	// SelectDictTypePage 根据条件分页查询字典类型
	SelectDictTypePage(query map[string]any) map[string]any

	// SelectDictTypeList 根据条件查询字典类型
	SelectDictTypeList(sysDictType model.SysDictType) []model.SysDictType

	// SelectDictTypeByID 根据字典类型ID查询信息
	SelectDictTypeByID(dictID string) model.SysDictType

	// SelectDictTypeByType 根据字典类型查询信息
	SelectDictTypeByType(dictType string) model.SysDictType

	// CheckUniqueDictName 校验字典名称是否唯一
	CheckUniqueDictName(dictName, dictID string) bool

	// CheckUniqueDictType 校验字典类型是否唯一
	CheckUniqueDictType(dictType, dictID string) bool

	// InsertDictType 新增字典类型信息
	InsertDictType(sysDictType model.SysDictType) string

	// UpdateDictType 修改字典类型信息
	UpdateDictType(sysDictType model.SysDictType) int64

	// DeleteDictTypeByIDs 批量删除字典类型信息
	DeleteDictTypeByIDs(dictIDs []string) (int64, error)

	// ResetDictCache 重置字典缓存数据
	ResetDictCache()

	// 加载字典缓存数据
	LoadingDictCache(dictType string)

	// 清空字典缓存数据
	ClearDictCache(dictType string) bool

	// DictDataCache 获取字典数据缓存数据
	DictDataCache(dictType string) []model.SysDictData
}
