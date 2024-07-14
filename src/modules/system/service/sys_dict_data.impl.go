package service

import (
	"errors"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysDictData 实例化服务层
var NewSysDictData = &SysDictDataService{
	sysDictDataRepository: repository.NewSysDictData,
	sysDictTypeService:    NewSysDictType,
}

// SysDictDataService 字典类型数据 服务层处理
type SysDictDataService struct {
	sysDictDataRepository repository.ISysDictDataRepository // 字典数据服务
	sysDictTypeService    ISysDictTypeService               // 字典类型服务
}

// FindByPage 分页查询列表数据
func (r *SysDictDataService) FindByPage(query map[string]any) map[string]any {
	return r.sysDictDataRepository.SelectByPage(query)
}

// Find 查询数据
func (r *SysDictDataService) Find(sysDictData model.SysDictData) []model.SysDictData {
	return r.sysDictDataRepository.Select(sysDictData)
}

// FindByCode 通过Code查询信息
func (r *SysDictDataService) FindByCode(dictCode string) model.SysDictData {
	if dictCode == "" {
		return model.SysDictData{}
	}
	dictCodes := r.sysDictDataRepository.SelectByCodes([]string{dictCode})
	if len(dictCodes) > 0 {
		return dictCodes[0]
	}
	return model.SysDictData{}
}

// FindByType 根据字典类型查询信息
func (r *SysDictDataService) FindByType(dictType string) []model.SysDictData {
	return r.sysDictTypeService.DictDataCache(dictType)
}

// Insert 新增信息
func (r *SysDictDataService) Insert(sysDictData model.SysDictData) string {
	insertId := r.sysDictDataRepository.Insert(sysDictData)
	if insertId != "" {
		r.sysDictTypeService.LoadingDictCache(sysDictData.DictType)
	}
	return insertId
}

// Update 修改信息
func (r *SysDictDataService) Update(sysDictData model.SysDictData) int64 {
	rows := r.sysDictDataRepository.Update(sysDictData)
	if rows > 0 {
		r.sysDictTypeService.LoadingDictCache(sysDictData.DictType)
	}
	return rows
}

// DeleteByCodes 批量删除信息
func (r *SysDictDataService) DeleteByCodes(dictCodes []string) (int64, error) {
	// 检查是否存在
	arr := r.sysDictDataRepository.SelectByCodes(dictCodes)
	if len(arr) <= 0 {
		return 0, errors.New("没有权限访问字典编码数据！")
	}
	if len(arr) == len(dictCodes) {
		for _, v := range arr {
			// 刷新缓存
			r.sysDictTypeService.CleanDictCache(v.DictType)
			r.sysDictTypeService.LoadingDictCache(v.DictType)
		}
		rows := r.sysDictDataRepository.DeleteByCodes(dictCodes)
		return rows, nil
	}
	return 0, errors.New("删除字典数据信息失败！")
}

// CheckUniqueTypeByLabel 检查同字典类型下字典标签是否唯一
func (r *SysDictDataService) CheckUniqueTypeByLabel(dictType, dictLabel, dictCode string) bool {
	uniqueId := r.sysDictDataRepository.CheckUnique(model.SysDictData{
		DictType:  dictType,
		DictLabel: dictLabel,
	})
	if uniqueId == dictCode {
		return true
	}
	return uniqueId == ""
}

// CheckUniqueTypeByValue 检查同字典类型下字典键值是否唯一
func (r *SysDictDataService) CheckUniqueTypeByValue(dictType, dictValue, dictCode string) bool {
	uniqueId := r.sysDictDataRepository.CheckUnique(model.SysDictData{
		DictType:  dictType,
		DictValue: dictValue,
	})
	if uniqueId == dictCode {
		return true
	}
	return uniqueId == ""
}
