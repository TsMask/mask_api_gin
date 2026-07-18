package service

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"gorm.io/gorm"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/database/redis"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
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
func (s SysDictType) FindByPage(query map[string]string) ([]model.SysDictType, int64) {
	return s.sysDictTypeRepository.SelectByPage(query)
}

// Find 查询数据
func (s SysDictType) Find(sysDictType model.SysDictType) []model.SysDictType {
	return s.sysDictTypeRepository.Select(sysDictType)
}

// FindById 通过ID查询信息
func (s SysDictType) FindById(id int64) model.SysDictType {
	if id <= 0 {
		return model.SysDictType{}
	}
	dictTypes := s.sysDictTypeRepository.SelectByIds([]int64{id})
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
func (s SysDictType) Insert(sysDictType model.SysDictType) int64 {
	if sysDictType.CreateBy != "" {
		ms := time.Now().UnixMilli()
		sysDictType.UpdateBy = sysDictType.CreateBy
		sysDictType.UpdateTime = ms
		sysDictType.CreateTime = ms
	}
	insertId := s.sysDictTypeRepository.Insert(sysDictType)
	if insertId > 0 {
		s.CacheLoad(sysDictType.DictType)
	}
	return insertId
}

// Update 修改信息
func (s SysDictType) Update(sysDictType model.SysDictType) int64 {
	if sysDictType.UpdateBy != "" {
		sysDictType.UpdateTime = time.Now().UnixMilli()
	}
	arr := s.sysDictTypeRepository.SelectByIds([]int64{sysDictType.ID})
	if len(arr) == 0 {
		return 0
	}
	// 同字典类型被修改时，同步更新修改
	oldDictType := arr[0].DictType
	var rows int64
	err := db.Transaction("", func(tx *gorm.DB) error {
		affected, err := s.sysDictTypeRepository.Update(tx, sysDictType)
		rows = affected
		if err != nil {
			return err
		}
		if rows > 0 && oldDictType != "" && oldDictType != sysDictType.DictType {
			if _, err := s.sysDictDataRepository.UpdateDataByDictType(tx, oldDictType, sysDictType.DictType); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0
	}
	// 刷新缓存
	s.CacheLoad(sysDictType.DictType)
	return rows
}

// DeleteByIds 批量删除信息
func (s SysDictType) DeleteByIds(ids []int64) (int64, error) {
	// 检查是否存在
	arr := s.sysDictTypeRepository.SelectByIds(ids)
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
	if len(arr) == len(ids) {
		return s.sysDictTypeRepository.DeleteByIds(ids), nil
	}
	return 0, fmt.Errorf("删除字典数据信息失败！")
}

// CheckUniqueByName 检查字典名称是否唯一
func (s SysDictType) CheckUniqueByName(dictName string, id int64) bool {
	uniqueId := s.sysDictTypeRepository.CheckUnique(model.SysDictType{
		DictName: dictName,
	})
	if uniqueId == id {
		return true
	}
	return uniqueId == 0
}

// CheckUniqueByType 检查字典类型是否唯一
func (s SysDictType) CheckUniqueByType(dictType string, id int64) bool {
	uniqueId := s.sysDictTypeRepository.CheckUnique(model.SysDictType{
		DictType: dictType,
	})
	if uniqueId == id {
		return true
	}
	return uniqueId == 0
}

// FindDataByType 获取字典数据缓存数据
func (s SysDictType) FindDataByType(dictType string) []model.SysDictData {
	var data []model.SysDictData
	key := constants.CACHE_SYS_DICT + ":" + dictType
	jsonStr, _ := redis.Get("", key)
	if len(jsonStr) > 7 {
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			data = []model.SysDictData{}
		}
	} else {
		data = s.sysDictDataRepository.Select(model.SysDictData{
			DataStatus: constants.STATUS_YES,
			DictType:   dictType,
		})
		if len(data) > 0 {
			_ = redis.Del("", key)
			values, _ := json.Marshal(data)
			_ = redis.Set("", key, string(values), 0)
		}
	}
	return data
}

// CacheLoad 加载字典缓存数据 传入*查询全部
func (s SysDictType) CacheLoad(dictType string) {
	sysDictData := model.SysDictData{
		DictType:   dictType,
		DataStatus: constants.STATUS_YES,
	}

	// 指定字典类型
	if dictType == "*" || dictType == "" {
		sysDictData.DictType = ""
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
		key := constants.CACHE_SYS_DICT + ":" + k
		_ = redis.Del("", key)
		values, _ := json.Marshal(v)
		_ = redis.Set("", key, string(values), 0)
	}
}

// CacheClean 清空字典缓存数据 传入*清除全部
func (s SysDictType) CacheClean(dictType string) bool {
	key := constants.CACHE_SYS_DICT + ":" + dictType
	keys, err := redis.GetKeys("", key)
	if err != nil {
		return false
	}
	return redis.DelKeys("", keys) == nil
}

// ExportData 导出数据表格
func (s SysDictType) ExportData(rows []model.SysDictType, fileName string) (string, error) {
	// 第一行表头标题
	data := [][]string{
		{"字典编号", "字典名称", "字典类型", "字典状态"},
	}
	// 从第二行开始的数据
	for _, row := range rows {
		statusValue := "停用"
		if row.DictStatus == constants.STATUS_YES {
			statusValue = "正常"
		}
		data = append(data, []string{
			fmt.Sprintf("%d", row.ID),
			row.DictName,
			row.DictType,
			statusValue,
		})
	}

	// 导出文件路径
	_, dir := file.UploadResourceDir()
	datePath := date.ParseDatePath()
	filePath := filepath.Join(dir, constants.UPLOAD_EXPORT, datePath, fileName)

	// 导出数据表格
	return filePath, file.CSVWrite(filePath, data, false)
}
