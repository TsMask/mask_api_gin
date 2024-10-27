package service

import (
	"encoding/json"
	"fmt"
	constCacheKey "mask_api_gin/src/framework/constants/cache_key"
	constSystem "mask_api_gin/src/framework/constants/system"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysDictType 实例化服务层
var NewSysDictType = &SysDictType{
	sysDictTypeRepository: repository.NewSysDictType,
	sysDictDataRepository: repository.NewSysDictData,
}

// SysDictType 字典类型 服务层处理
type SysDictType struct {
	sysDictTypeRepository *repository.SysDictType // 字典类型服务
	sysDictDataRepository *repository.SysDictData // 字典数据服务
}

// FindByPage 分页查询列表数据
func (s SysDictType) FindByPage(query map[string]any) ([]model.SysDictType, int64) {
	return s.sysDictTypeRepository.SelectByPage(query)
}

// Find 查询数据
func (s SysDictType) Find(sysDictType model.SysDictType) []model.SysDictType {
	return s.sysDictTypeRepository.Select(sysDictType)
}

// FindById 通过ID查询信息
func (s SysDictType) FindById(dictId string) model.SysDictType {
	if dictId == "" {
		return model.SysDictType{}
	}
	dictTypes := s.sysDictTypeRepository.SelectByIds([]string{dictId})
	if len(dictTypes) > 0 {
		return dictTypes[0]
	}
	return model.SysDictType{}
}

// FindByType 根据字典类型查询信息
func (s SysDictType) FindByType(dictType string) model.SysDictType {
	return s.sysDictTypeRepository.SelectByType(dictType)
}

// Insert 新增信息
func (s SysDictType) Insert(sysDictType model.SysDictType) string {
	insertId := s.sysDictTypeRepository.Insert(sysDictType)
	if insertId != "" {
		s.CacheLoad(sysDictType.DictType)
	}
	return insertId
}

// Update 修改信息
func (s SysDictType) Update(sysDictType model.SysDictType) int64 {
	arr := s.sysDictTypeRepository.SelectByIds([]string{sysDictType.DictId})
	if len(arr) == 0 {
		return 0
	}
	// 同字典类型被修改时，同步更新修改
	oldDictType := arr[0].DictType
	rows := s.sysDictTypeRepository.Update(sysDictType)
	if rows > 0 && oldDictType != "" && oldDictType != sysDictType.DictType {
		s.sysDictDataRepository.UpdateDataByDictType(oldDictType, sysDictType.DictType)
	}
	// 刷新缓存
	s.CacheLoad(sysDictType.DictType)
	return rows
}

// DeleteByIds 批量删除信息
func (s SysDictType) DeleteByIds(dictIds []string) (int64, error) {
	// 检查是否存在
	arr := s.sysDictTypeRepository.SelectByIds(dictIds)
	if len(arr) <= 0 {
		return 0, fmt.Errorf("没有权限访问字典类型数据！")
	}
	for _, v := range arr {
		// 字典类型下级含有数据
		if useCount := s.sysDictDataRepository.ExistDataByDictType(v.DictType); useCount > 0 {
			return 0, fmt.Errorf("【%s】存在字典数据,不能删除", v.DictName)
		}
		// 清除缓存
		s.CacheClean(v.DictType)
	}
	if len(arr) == len(dictIds) {
		return s.sysDictTypeRepository.DeleteByIds(dictIds), nil
	}
	return 0, fmt.Errorf("删除字典数据信息失败！")
}

// CheckUniqueByName 检查字典名称是否唯一
func (s SysDictType) CheckUniqueByName(dictName, dictId string) bool {
	uniqueId := s.sysDictTypeRepository.CheckUnique(model.SysDictType{
		DictName: dictName,
	})
	if uniqueId == dictId {
		return true
	}
	return uniqueId == ""
}

// CheckUniqueByType 检查字典类型是否唯一
func (s SysDictType) CheckUniqueByType(dictType, dictId string) bool {
	uniqueId := s.sysDictTypeRepository.CheckUnique(model.SysDictType{
		DictType: dictType,
	})
	if uniqueId == dictId {
		return true
	}
	return uniqueId == ""
}

// getCacheKey 组装缓存key
func (s SysDictType) getCacheKey(dictType string) string {
	return constCacheKey.SYS_DICT_KEY + dictType
}

// CacheLoad 加载字典缓存数据 传入*查询全部
func (s SysDictType) CacheLoad(dictType string) {
	sysDictData := model.SysDictData{
		DictType: dictType,
		Status:   constSystem.STATUS_YES,
	}

	// 指定字典类型
	if dictType == "*" || dictType == "" {
		sysDictData.DictType = dictType
	}

	arr := s.sysDictDataRepository.Select(sysDictData)
	if len(arr) == 0 {
		return
	}

	// 将字典数据按类型分组
	m := make(map[string][]model.SysDictData)
	for _, v := range arr {
		key := v.DictType
		if item, ok := m[key]; ok {
			m[key] = append(item, v)
		} else {
			m[key] = []model.SysDictData{v}
		}
	}

	// 放入缓存
	for k, v := range m {
		key := s.getCacheKey(k)
		_ = redis.Del("", key)
		values, _ := json.Marshal(v)
		_ = redis.Set("", key, string(values))
	}
}

// CacheClean 清空字典缓存数据 传入*清除全部
func (s SysDictType) CacheClean(dictType string) bool {
	key := s.getCacheKey(dictType)
	keys, err := redis.GetKeys("", key)
	if err != nil {
		return false
	}
	return redis.DelKeys("", keys) == nil
}

// FindDataByType 获取字典数据缓存数据
func (s SysDictType) FindDataByType(dictType string) []model.SysDictData {
	var data []model.SysDictData
	key := s.getCacheKey(dictType)
	jsonStr, _ := redis.Get("", key)
	if len(jsonStr) > 7 {
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			data = []model.SysDictData{}
		}
	} else {
		data = s.sysDictDataRepository.Select(model.SysDictData{
			Status:   constSystem.STATUS_YES,
			DictType: dictType,
		})
		if len(data) > 0 {
			_ = redis.Del("", key)
			values, _ := json.Marshal(data)
			_ = redis.Set("", key, string(values))
		}
	}
	return data
}
