package service

import (
	"errors"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// 实例化服务层 SysDictDataImpl 结构体
var NewSysDictDataImpl = &SysDictDataImpl{
	sysDictDataRepository: repository.NewSysDictDataImpl,
	sysDictTypeService:    NewSysDictTypeImpl,
}

// SysDictDataImpl 字典类型数据 服务层处理
type SysDictDataImpl struct {
	// 字典数据服务
	sysDictDataRepository repository.ISysDictData
	// 字典类型服务
	sysDictTypeService ISysDictType
}

// SelectDictDataPage 根据条件分页查询字典数据
func (r *SysDictDataImpl) SelectDictDataPage(query map[string]string) map[string]interface{} {
	return r.sysDictDataRepository.SelectDictDataPage(query)
}

// SelectDictDataList 根据条件查询字典数据
func (r *SysDictDataImpl) SelectDictDataList(sysDictData model.SysDictData) []model.SysDictData {
	return r.sysDictDataRepository.SelectDictDataList(sysDictData)
}

// SelectDictDataByCode 根据字典数据编码查询信息
func (r *SysDictDataImpl) SelectDictDataByCode(dictCode string) model.SysDictData {
	if dictCode == "" {
		return model.SysDictData{}
	}
	dictCodes := r.sysDictDataRepository.SelectDictDataByCodes([]string{dictCode})
	if len(dictCodes) > 0 {
		return dictCodes[0]
	}
	return model.SysDictData{}
}

// SelectDictDataByType 根据字典类型查询信息
func (r *SysDictDataImpl) SelectDictDataByType(dictType string) []model.SysDictData {
	return r.sysDictTypeService.DictDataCache(dictType)
}

// CheckUniqueDictLabel 校验字典标签是否唯一
func (r *SysDictDataImpl) CheckUniqueDictLabel(dictType, dictLabel, dictCode string) bool {
	uniqueId := r.sysDictDataRepository.CheckUniqueDictData(model.SysDictData{
		DictType:  dictType,
		DictLabel: dictLabel,
	})
	if uniqueId == dictCode {
		return true
	}
	return uniqueId == ""
}

// CheckUniqueDictValue 校验字典键值是否唯一
func (r *SysDictDataImpl) CheckUniqueDictValue(dictType, dictValue, dictCode string) bool {
	uniqueId := r.sysDictDataRepository.CheckUniqueDictData(model.SysDictData{
		DictType:  dictType,
		DictValue: dictValue,
	})
	if uniqueId == dictCode {
		return true
	}
	return uniqueId == ""
}

// DeleteDictDataByCodes 批量删除字典数据信息
func (r *SysDictDataImpl) DeleteDictDataByCodes(dictCodes []string) (int64, error) {
	// 检查是否存在
	dictDatas := r.sysDictDataRepository.SelectDictDataByCodes(dictCodes)
	if len(dictDatas) <= 0 {
		return 0, errors.New("没有权限访问字典编码数据！")
	}
	if len(dictDatas) == len(dictCodes) {
		for _, v := range dictDatas {
			// 刷新缓存
			r.sysDictTypeService.ClearDictCache(v.DictType)
			r.sysDictTypeService.LoadingDictCache(v.DictType)
		}
		rows := r.sysDictDataRepository.DeleteDictDataByCodes(dictCodes)
		return rows, nil
	}
	return 0, errors.New("删除字典数据信息失败！")
}

// InsertDictData 新增字典数据信息
func (r *SysDictDataImpl) InsertDictData(sysDictData model.SysDictData) string {
	insertId := r.sysDictDataRepository.InsertDictData(sysDictData)
	if insertId != "" {
		r.sysDictTypeService.LoadingDictCache(sysDictData.DictType)
	}
	return insertId
}

// UpdateDictData 修改字典数据信息
func (r *SysDictDataImpl) UpdateDictData(sysDictData model.SysDictData) int64 {
	rows := r.sysDictDataRepository.UpdateDictData(sysDictData)
	if rows > 0 {
		r.sysDictTypeService.LoadingDictCache(sysDictData.DictType)
	}
	return rows
}
