package repository

import "mask_api_gin/src/modules/system/model"

// ISysDictData 字典类型数据表 数据层接口
type ISysDictData interface {
	// SelectDictDataPage 根据条件分页查询字典数据
	SelectDictDataPage(query map[string]string) map[string]interface{}

	// SelectDictDataList 根据条件查询字典数据
	SelectDictDataList(sysDictData model.SysDictData) []model.SysDictData

	// SelectDictDataByCodes 根据字典数据编码查询信息
	SelectDictDataByCodes(dictCodes []string) []model.SysDictData

	// CountDictDataByType 查询字典数据
	CountDictDataByType(dictType string) int64

	// CheckUniqueDictData 校验字典数据是否唯一
	CheckUniqueDictData(sysDictData model.SysDictData) string

	// DeleteDictDataByCodes 批量删除字典数据信息
	DeleteDictDataByCodes(dictCodes []string) int64

	// InsertDictData 新增字典数据信息
	InsertDictData(sysDictData model.SysDictData) string

	// UpdateDictData 修改字典数据信息
	UpdateDictData(sysDictData model.SysDictData) int64

	// UpdateDictDataType 同步修改字典类型
	UpdateDictDataType(oldDictType string, newDictType string) int64
}
