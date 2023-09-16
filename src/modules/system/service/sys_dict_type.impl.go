package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"mask_api_gin/src/framework/constants/cachekey"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/redis"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// 实例化服务层 SysDictTypeImpl 结构体
var NewSysDictTypeImpl = &SysDictTypeImpl{
	sysDictTypeRepository: repository.NewSysDictTypeImpl,
	sysDictDataRepository: repository.NewSysDictDataImpl,
}

// SysDictTypeImpl 字典类型 服务层处理
type SysDictTypeImpl struct {
	// 字典类型服务
	sysDictTypeRepository repository.ISysDictType
	// 字典数据服务
	sysDictDataRepository repository.ISysDictData
}

// SelectDictTypePage 根据条件分页查询字典类型
func (r *SysDictTypeImpl) SelectDictTypePage(query map[string]any) map[string]any {
	return r.sysDictTypeRepository.SelectDictTypePage(query)
}

// SelectDictTypeList 根据条件查询字典类型
func (r *SysDictTypeImpl) SelectDictTypeList(sysDictType model.SysDictType) []model.SysDictType {
	return r.sysDictTypeRepository.SelectDictTypeList(sysDictType)
}

// SelectDictTypeByID 根据字典类型ID查询信息
func (r *SysDictTypeImpl) SelectDictTypeByID(dictID string) model.SysDictType {
	if dictID == "" {
		return model.SysDictType{}
	}
	dictTypes := r.sysDictTypeRepository.SelectDictTypeByIDs([]string{dictID})
	if len(dictTypes) > 0 {
		return dictTypes[0]
	}
	return model.SysDictType{}
}

// SelectDictTypeByType 根据字典类型查询信息
func (r *SysDictTypeImpl) SelectDictTypeByType(dictType string) model.SysDictType {
	return r.sysDictTypeRepository.SelectDictTypeByType(dictType)
}

// CheckUniqueDictName 校验字典名称是否唯一
func (r *SysDictTypeImpl) CheckUniqueDictName(dictName, dictID string) bool {
	uniqueId := r.sysDictTypeRepository.CheckUniqueDictType(model.SysDictType{
		DictName: dictName,
	})
	if uniqueId == dictID {
		return true
	}
	return uniqueId == ""
}

// CheckUniqueDictType 校验字典类型是否唯一
func (r *SysDictTypeImpl) CheckUniqueDictType(dictType, dictID string) bool {
	uniqueId := r.sysDictTypeRepository.CheckUniqueDictType(model.SysDictType{
		DictType: dictType,
	})
	if uniqueId == dictID {
		return true
	}
	return uniqueId == ""
}

// InsertDictType 新增字典类型信息
func (r *SysDictTypeImpl) InsertDictType(sysDictType model.SysDictType) string {
	insertId := r.sysDictTypeRepository.InsertDictType(sysDictType)
	if insertId != "" {
		r.LoadingDictCache(sysDictType.DictType)
	}
	return insertId
}

// UpdateDictType 修改字典类型信息
func (r *SysDictTypeImpl) UpdateDictType(sysDictType model.SysDictType) int64 {
	data := r.sysDictTypeRepository.SelectDictTypeByIDs([]string{sysDictType.DictID})
	if len(data) == 0 {
		return 0
	}
	// 修改字典类型key时同步更新其字典数据的类型key
	oldDictType := data[0].DictType
	rows := r.sysDictTypeRepository.UpdateDictType(sysDictType)
	if rows > 0 && oldDictType != "" && oldDictType != sysDictType.DictType {
		r.sysDictDataRepository.UpdateDictDataType(oldDictType, sysDictType.DictType)
	}
	// 刷新缓存
	r.ClearDictCache(oldDictType)
	r.LoadingDictCache(sysDictType.DictType)
	return rows
}

// DeleteDictTypeByIDs 批量删除字典类型信息
func (r *SysDictTypeImpl) DeleteDictTypeByIDs(dictIDs []string) (int64, error) {
	// 检查是否存在
	dictTypes := r.sysDictTypeRepository.SelectDictTypeByIDs(dictIDs)
	if len(dictTypes) <= 0 {
		return 0, errors.New("没有权限访问字典类型数据！")
	}
	for _, v := range dictTypes {
		// 字典类型下级含有数据
		useCount := r.sysDictDataRepository.CountDictDataByType(v.DictType)
		if useCount > 0 {
			msg := fmt.Sprintf("【%s】存在字典数据,不能删除", v.DictName)
			return 0, errors.New(msg)
		}
		// 清除缓存
		r.ClearDictCache(v.DictType)
	}
	if len(dictTypes) == len(dictIDs) {
		rows := r.sysDictTypeRepository.DeleteDictTypeByIDs(dictIDs)
		return rows, nil
	}
	return 0, errors.New("删除字典数据信息失败！")
}

// ResetDictCache 重置字典缓存数据
func (r *SysDictTypeImpl) ResetDictCache() {
	r.ClearDictCache("*")
	r.LoadingDictCache("")
}

// getCacheKey 组装缓存key
func (r *SysDictTypeImpl) getDictCache(dictType string) string {
	return cachekey.SYS_DICT_KEY + dictType
}

// LoadingDictCache 加载字典缓存数据
func (r *SysDictTypeImpl) LoadingDictCache(dictType string) {
	sysDictData := model.SysDictData{
		Status: common.STATUS_YES,
	}

	// 指定字典类型
	if dictType != "" {
		sysDictData.DictType = dictType
		// 删除缓存
		key := r.getDictCache(dictType)
		redis.Del("", key)
	}

	sysDictDataList := r.sysDictDataRepository.SelectDictDataList(sysDictData)
	if len(sysDictDataList) == 0 {
		return
	}

	// 将字典数据按类型分组
	m := make(map[string][]model.SysDictData, 0)
	for _, v := range sysDictDataList {
		key := v.DictType
		if item, ok := m[key]; ok {
			m[key] = append(item, v)
		} else {
			m[key] = []model.SysDictData{v}
		}
	}

	// 放入缓存
	for k, v := range m {
		key := r.getDictCache(k)
		values, _ := json.Marshal(v)
		redis.Set("", key, string(values))
	}
}

// ClearDictCache 清空字典缓存数据
func (r *SysDictTypeImpl) ClearDictCache(dictType string) bool {
	key := r.getDictCache(dictType)
	keys, err := redis.GetKeys("", key)
	if err != nil {
		return false
	}
	delOk, _ := redis.DelKeys("", keys)
	return delOk
}

// DictDataCache 获取字典数据缓存数据
func (r *SysDictTypeImpl) DictDataCache(dictType string) []model.SysDictData {
	data := []model.SysDictData{}
	key := r.getDictCache(dictType)
	jsonStr, _ := redis.Get("", key)
	if len(jsonStr) > 7 {
		err := json.Unmarshal([]byte(jsonStr), &data)
		if err != nil {
			data = []model.SysDictData{}
		}
	} else {
		data = r.sysDictDataRepository.SelectDictDataList(model.SysDictData{
			Status:   common.STATUS_YES,
			DictType: dictType,
		})
		if len(data) > 0 {
			redis.Del("", key)
			values, _ := json.Marshal(data)
			redis.Set("", key, string(values))
		}
	}
	return data
}
