package repository

import "mask_api_gin/src/modules/system/model"

// ISysDictData 字典类型数据表 数据层接口
type ISysDictData interface {
	// SelectDictDataPage 根据条件分页查询字典数据
	SelectDictDataPage(query map[string]string) map[string]interface{}

	// SelectDictDataList 根据条件查询字典数据
	SelectDictDataList(sysDictData model.SysDictData) []model.SysDictData

	// SelectDictLabel 根据字典类型和字典键值查询字典数据信息
	SelectDictLabel(dictType string, dictValue string) (string, error)

	// SelectDictDataByCode 根据字典数据编码查询信息
	SelectDictDataByCode(dictCode string) model.SysDictData

	// CountDictDataByType 查询字典数据
	CountDictDataByType(dictType string) string

	// CheckUniqueDictLabel 校验字典标签是否唯一
	CheckUniqueDictLabel(dictType string, dictLabel string) string

	// CheckUniqueDictValue 校验字典键值是否唯一
	CheckUniqueDictValue(dictType string, dictValue string) string

	// DeleteDictDataByCodes 批量删除字典数据信息
	DeleteDictDataByCodes(dictCodes []string) int

	// InsertDictData 新增字典数据信息
	InsertDictData(sysDictData model.SysDictData) string

	// UpdateDictData 修改字典数据信息
	UpdateDictData(sysDictData model.SysDictData) int

	// UpdateDictDataType 同步修改字典类型
	UpdateDictDataType(oldDictType string, newDictType string) int
}
