package repository

import (
	"time"

	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"

	"gorm.io/gorm"
)

// NewSysDictData 实例化数据层
var NewSysDictData = &SysDictData{}

// SysDictData 字典类型数据表 数据层处理
type SysDictData struct{}

// sysDictDataSortMapping 字典数据排序字段映射
var sysDictDataSortMapping = map[string]string{
	"id":         "id",
	"createTime": "create_time",
	"updateTime": "update_time",
	"dictType":   "dict_type",
	"dataSort":   "data_sort",
	"tagClass":   "tag_class",
	"tagType":    "tag_type",
	"dataStatus": "data_status",
}

// SelectByPage 分页查询集合
func (r SysDictData) SelectByPage(query map[string]string) ([]model.SysDictData, int64) {
	tx := db.DB("").Model(&model.SysDictData{})
	// 查询条件拼接
	if v, ok := query["dictType"]; ok && v != "" {
		tx = tx.Where("dict_type = ?", v)
	}
	if v, ok := query["dataLabel"]; ok && v != "" {
		tx = tx.Where("data_label like ?", v+"%")
	}
	if v, ok := query["dataStatus"]; ok && v != "" {
		tx = tx.Where("data_status = ?", v)
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysDictData{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return rows, total
	}

	// 查询数据分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序，默认按 data_sort ASC
	orders := db.SortParse(query["sort"], sysDictDataSortMapping)
	tx = db.SortOrDefault(tx, orders, "data_sort asc")

	err := tx.Find(&rows).Error
	if err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows, total
	}
	return rows, total
}

// Select 查询集合
func (r SysDictData) Select(sysDictData model.SysDictData) []model.SysDictData {
	tx := db.DB("").Model(&model.SysDictData{})
	// 查询条件拼接
	if sysDictData.DataLabel != "" {
		tx = tx.Where("data_label like ?", sysDictData.DataLabel+"%")
	}
	if sysDictData.DictType != "" {
		tx = tx.Where("dict_type = ?", sysDictData.DictType)
	}
	if sysDictData.DataStatus != "" {
		tx = tx.Where("data_status = ?", sysDictData.DataStatus)
	}

	// 查询数据
	rows := []model.SysDictData{}
	if err := tx.Order("data_sort asc").Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// SelectByIds 通过ID查询信息
func (r SysDictData) SelectByIds(ids []int64) []model.SysDictData {
	rows := []model.SysDictData{}
	if len(ids) <= 0 {
		return rows
	}
	tx := db.DB("").Model(&model.SysDictData{})
	// 构建查询条件
	tx = tx.Where("id in ?", ids)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// Insert 新增信息 返回新增数据ID
func (r SysDictData) Insert(sysDictData model.SysDictData) int64 {
	if sysDictData.CreateBy != "" {
		ms := time.Now().UnixMilli()
		sysDictData.UpdateBy = sysDictData.CreateBy
		sysDictData.UpdateTime = ms
		sysDictData.CreateTime = ms
	}
	// 执行插入
	if err := db.DB("").Create(&sysDictData).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
	}
	return sysDictData.ID
}

// Update 修改信息 返回受影响行数
func (r SysDictData) Update(sysDictData model.SysDictData) int64 {
	if sysDictData.ID <= 0 {
		return 0
	}
	if sysDictData.UpdateBy != "" {
		sysDictData.UpdateTime = time.Now().UnixMilli()
	}
	tx := db.DB("").Model(&model.SysDictData{})
	// 构建查询条件
	tx = tx.Where("id = ?", sysDictData.ID)
	tx = tx.Omit("id", "create_by", "create_time")
	// 执行更新
	if err := tx.Updates(sysDictData).Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByIds 批量删除信息 返回受影响行数（物理删除）
func (r SysDictData) DeleteByIds(ids []int64) int64 {
	if len(ids) <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.SysDictData{})
	// 构建查询条件
	tx = tx.Where("id in ?", ids)
	// 执行物理删除
	if err := tx.Delete(&model.SysDictData{}).Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// CheckUnique 检查信息是否唯一 返回数据ID
func (r SysDictData) CheckUnique(sysDictData model.SysDictData) int64 {
	tx := db.DB("").Model(&model.SysDictData{})
	// 查询条件拼接
	if sysDictData.DictType != "" {
		tx = tx.Where("dict_type = ?", sysDictData.DictType)
	}
	if sysDictData.DataLabel != "" {
		tx = tx.Where("data_label = ?", sysDictData.DataLabel)
	}
	if sysDictData.DataValue != "" {
		tx = tx.Where("data_value = ?", sysDictData.DataValue)
	}
	// 查询数据
	var id int64 = 0
	if err := tx.Select("id").Limit(1).Find(&id).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return id
	}
	return id
}

// ExistDataByDictType 存在数据数量
func (r SysDictData) ExistDataByDictType(dictType string) int64 {
	if dictType == "" {
		return 0
	}
	tx := db.DB("").Model(&model.SysDictData{})
	tx = tx.Where("dict_type = ?", dictType)
	// 查询数据
	var count int64 = 0
	if err := tx.Count(&count).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return count
	}
	return count
}

// UpdateDataByDictType 更新一组字典类型 返回受影响行数
func (r SysDictData) UpdateDataByDictType(d *gorm.DB, oldDictType string, newDictType string) int64 {
	if oldDictType == "" || newDictType == "" {
		return 0
	}
	if d == nil {
		d = db.DB("")
	}
	tx := d.Model(&model.SysDictData{})
	// 构建查询条件
	tx = tx.Where("dict_type = ?", oldDictType)
	// 执行更新字典类型
	if err := tx.Update("dict_type", newDictType).Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}
