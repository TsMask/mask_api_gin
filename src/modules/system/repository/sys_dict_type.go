package repository

import "mask_api_gin/src/modules/system/model"

// ISysDictType 字典类型表 数据层接口
type ISysDictType interface {
	// SelectDictTypePage 根据条件分页查询字典类型
	SelectDictTypePage(query map[string]any) map[string]any

	// SelectDictTypeList 根据条件查询字典类型
	SelectDictTypeList(sysDictType model.SysDictType) []model.SysDictType

	// SelectDictTypeByIDs 根据字典类型ID查询信息
	SelectDictTypeByIDs(dictIDs []string) []model.SysDictType

	// SelectDictTypeByType 根据字典类型查询信息
	SelectDictTypeByType(dictType string) model.SysDictType

	// CheckUniqueDictType 校验字典类型是否唯一
	CheckUniqueDictType(sysDictType model.SysDictType) string

	// InsertDictType 新增字典类型信息
	InsertDictType(sysDictType model.SysDictType) string

	// UpdateDictType 修改字典类型信息
	UpdateDictType(sysDictType model.SysDictType) int64

	// DeleteDictTypeByIDs 批量删除字典类型信息
	DeleteDictTypeByIDs(dictIDs []string) int64
}
