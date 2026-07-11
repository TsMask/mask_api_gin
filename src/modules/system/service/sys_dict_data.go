package service

import (
	"fmt"
	"path/filepath"

	"mask_api_gin/src/framework/constants"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/framework/utils/file"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
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
func (s SysDictData) FindById(id int64) model.SysDictData {
	if id <= 0 {
		return model.SysDictData{}
	}
	dicts := s.sysDictDataRepository.SelectByIds([]int64{id})
	if len(dicts) > 0 {
		return dicts[0]
	}
	return model.SysDictData{}
}

// FindByType 根据字典类型查询信息
func (s SysDictData) FindByType(dictType string) []model.SysDictData {
	return s.sysDictTypeService.FindDataByType(dictType)
}

// Insert 新增信息
func (s SysDictData) Insert(sysDictData model.SysDictData) int64 {
	insertId := s.sysDictDataRepository.Insert(sysDictData)
	if insertId > 0 {
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
func (s SysDictData) DeleteByIds(ids []int64) (int64, error) {
	// 检查是否存在
	dicts := s.sysDictDataRepository.SelectByIds(ids)
	if len(dicts) <= 0 {
		return 0, fmt.Errorf("没有权限访问字典编码数据！")
	}
	if len(dicts) == len(ids) {
		for _, v := range dicts {
			// 刷新缓存
			s.sysDictTypeService.CacheClean(v.DictType)
			s.sysDictTypeService.CacheLoad(v.DictType)
		}
		rows := s.sysDictDataRepository.DeleteByIds(ids)
		return rows, nil
	}
	return 0, fmt.Errorf("删除字典数据信息失败！")
}

// CheckUniqueTypeByLabel 检查同字典类型下字典标签是否唯一
func (s SysDictData) CheckUniqueTypeByLabel(dictType, dataLabel string, id int64) bool {
	uniqueId := s.sysDictDataRepository.CheckUnique(model.SysDictData{
		DictType:  dictType,
		DataLabel: dataLabel,
	})
	if uniqueId == id {
		return true
	}
	return uniqueId == 0
}

// CheckUniqueTypeByValue 检查同字典类型下字典键值是否唯一
func (s SysDictData) CheckUniqueTypeByValue(dictType, dataValue string, id int64) bool {
	uniqueId := s.sysDictDataRepository.CheckUnique(model.SysDictData{
		DictType:  dictType,
		DataValue: dataValue,
	})
	if uniqueId == id {
		return true
	}
	return uniqueId == 0
}

// ExportData 导出数据表格
func (s SysDictData) ExportData(rows []model.SysDictData, fileName string) (string, error) {
	// 第一行表头标题
	data := [][]string{
		{"字典类型", "数据排序", "数据编号", "数据标签", "数据键值", "数据状态"},
	}
	// 从第二行开始的数据
	for _, row := range rows {
		statusValue := "停用"
		if row.DataStatus == constants.STATUS_YES {
			statusValue = "正常"
		}
		data = append(data, []string{
			row.DictType,
			fmt.Sprintf("%d", row.DataSort),
			fmt.Sprintf("%d", row.ID),
			row.DataLabel,
			row.DataValue,
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
