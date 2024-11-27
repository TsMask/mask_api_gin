package service

import (
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"

	"fmt"
)

// NewSysDictData 实例化服务层
var NewSysDictData = &SysDictData{
	sysDictDataRepository: repository.NewSysDictData,
	sysDictTypeService:    NewSysDictType,
}

// SysDictData 字典类型数据 服务层处理
type SysDictData struct {
	sysDictDataRepository *repository.SysDictData // 字典数据服务
	sysDictTypeService    *SysDictType            // 字典类型服务
}

// FindByPage 分页查询列表数据
func (s SysDictData) FindByPage(query map[string]string) ([]model.SysDictData, int64) {
	return s.sysDictDataRepository.SelectByPage(query)
}

// Find 查询数据
func (s SysDictData) Find(sysDictData model.SysDictData) []model.SysDictData {
	return s.sysDictDataRepository.Select(sysDictData)
}

// FindById 通过ID查询信息
func (s SysDictData) FindById(dictId string) model.SysDictData {
	if dictId == "" {
		return model.SysDictData{}
	}
	arr := s.sysDictDataRepository.SelectByIds([]string{dictId})
	if len(arr) > 0 {
		return arr[0]
	}
	return model.SysDictData{}
}

// FindByType 根据字典类型查询信息
func (s SysDictData) FindByType(dictType string) []model.SysDictData {
	return s.sysDictTypeService.FindDataByType(dictType)
}

// Insert 新增信息
func (s SysDictData) Insert(sysDictData model.SysDictData) string {
	insertId := s.sysDictDataRepository.Insert(sysDictData)
	if insertId != "" {
		s.sysDictTypeService.CacheLoad(sysDictData.DictType)
	}
	return insertId
}

// Update 修改信息
func (s SysDictData) Update(sysDictData model.SysDictData) int64 {
	rows := s.sysDictDataRepository.Update(sysDictData)
	if rows > 0 {
		s.sysDictTypeService.CacheLoad(sysDictData.DictType)
	}
	return rows
}

// DeleteByIds 批量删除信息
func (s SysDictData) DeleteByIds(dictIds []string) (int64, error) {
	// 检查是否存在
	arr := s.sysDictDataRepository.SelectByIds(dictIds)
	if len(arr) <= 0 {
		return 0, fmt.Errorf("没有权限访问字典编码数据！")
	}
	if len(arr) == len(dictIds) {
		for _, v := range arr {
			// 刷新缓存
			s.sysDictTypeService.CacheClean(v.DictType)
			s.sysDictTypeService.CacheLoad(v.DictType)
		}
		rows := s.sysDictDataRepository.DeleteByIds(dictIds)
		return rows, nil
	}
	return 0, fmt.Errorf("删除字典数据信息失败！")
}

// CheckUniqueTypeByLabel 检查同字典类型下字典标签是否唯一
func (s SysDictData) CheckUniqueTypeByLabel(dictType, dictLabel string, dataId string) bool {
	uniqueId := s.sysDictDataRepository.CheckUnique(model.SysDictData{
		DictType:  dictType,
		DataLabel: dictLabel,
	})
	if uniqueId == dataId {
		return true
	}
	return uniqueId == ""
}

// CheckUniqueTypeByValue 检查同字典类型下字典键值是否唯一
func (s SysDictData) CheckUniqueTypeByValue(dictType, dictValue string, dataId string) bool {
	uniqueId := s.sysDictDataRepository.CheckUnique(model.SysDictData{
		DictType:  dictType,
		DataValue: dictValue,
	})
	if uniqueId == dataId {
		return true
	}
	return uniqueId == ""
}
