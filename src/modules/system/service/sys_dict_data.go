package service

import "mask_api_gin/src/modules/system/model"

// ISysDictData 字典类型数据 服务层接口
type ISysDictData interface {
	// SelectDictDataPage 根据条件分页查询字典数据
	SelectDictDataPage(query map[string]any) map[string]any

	// SelectDictDataList 根据条件查询字典数据
	SelectDictDataList(sysDictData model.SysDictData) []model.SysDictData

	// SelectDictDataByCode 根据字典数据编码查询信息
	SelectDictDataByCode(dictCode string) model.SysDictData

	// SelectDictDataByType 根据字典类型查询信息
	SelectDictDataByType(dictType string) []model.SysDictData

	// CheckUniqueDictLabel 校验字典标签是否唯一
	CheckUniqueDictLabel(dictType, dictLabel, dictCode string) bool

	// CheckUniqueDictValue 校验字典键值是否唯一
	CheckUniqueDictValue(dictType, dictValue, dictCode string) bool

	// DeleteDictDataByCodes 批量删除字典数据信息
	DeleteDictDataByCodes(dictCodes []string) (int64, error)

	// InsertDictData 新增字典数据信息
	InsertDictData(sysDictData model.SysDictData) string

	// UpdateDictData 修改字典数据信息
	UpdateDictData(sysDictData model.SysDictData) int64
}
