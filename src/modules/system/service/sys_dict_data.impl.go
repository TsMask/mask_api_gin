package service

import (
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// SysDictDataImpl 字典类型数据 数据层处理
var SysDictDataImpl = &sysDictDataImpl{
	sysDictDataRepository: repository.SysDictDataImpl,
}

type sysDictDataImpl struct {
	// 字典类型数据服务
	sysDictDataRepository repository.ISysDictData
}

// SelectDictDataPage 根据条件分页查询字典数据
func (r *sysDictDataImpl) SelectDictDataPage(query map[string]string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectDictDataList 根据条件查询字典数据
func (r *sysDictDataImpl) SelectDictDataList(sysDictData model.SysDictData) []model.SysDictData {
	return []model.SysDictData{}
}

// SelectDictLabel 根据字典类型和字典键值查询字典数据信息
func (r *sysDictDataImpl) SelectDictLabel(dictType string, dictValue string) (string, error) {
	return "", nil
}

// SelectDictDataByCode 根据字典数据编码查询信息
func (r *sysDictDataImpl) SelectDictDataByCode(dictCode string) model.SysDictData {
	return model.SysDictData{}
}

// CountDictDataByType 查询字典数据
func (r *sysDictDataImpl) CountDictDataByType(dictType string) string {
	return ""
}

// CheckUniqueDictLabel 校验字典标签是否唯一
func (r *sysDictDataImpl) CheckUniqueDictLabel(dictType string, dictLabel string) string {
	return ""
}

// CheckUniqueDictValue 校验字典键值是否唯一
func (r *sysDictDataImpl) CheckUniqueDictValue(dictType string, dictValue string) string {
	return ""
}

// DeleteDictDataByCodes 批量删除字典数据信息
func (r *sysDictDataImpl) DeleteDictDataByCodes(dictCodes []string) int {
	return 0
}

// InsertDictData 新增字典数据信息
func (r *sysDictDataImpl) InsertDictData(sysDictData model.SysDictData) string {
	return ""
}

// UpdateDictData 修改字典数据信息
func (r *sysDictDataImpl) UpdateDictData(sysDictData model.SysDictData) int {
	return 0
}

// UpdateDictDataType 同步修改字典类型
func (r *sysDictDataImpl) UpdateDictDataType(oldDictType string, newDictType string) int {
	return 0
}
