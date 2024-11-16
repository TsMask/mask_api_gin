package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"
	"time"
)

// NewSysDictType 实例化数据层
var NewSysDictType = &SysDictType{}

// SysDictType 字典类型表 数据层处理
type SysDictType struct{}

// SelectByPage 分页查询集合
func (r SysDictType) SelectByPage(query map[string]any) ([]model.SysDictType, int64) {
	tx := db.DB("").Model(&model.SysDictType{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if v, ok := query["dictName"]; ok && v != "" {
		tx = tx.Where("dict_name like concat(?, '%')", v)
	}
	if v, ok := query["dictType"]; ok && v != "" {
		tx = tx.Where("dict_type like concat(?, '%')", v)
	}
	if v, ok := query["statusFlag"]; ok && v != "" {
		tx = tx.Where("status_flag = ?", v)
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysDictType{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return rows, total
	}

	// 查询数据分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	err := tx.Limit(pageSize).Offset(pageSize * pageNum).Find(&rows).Error
	if err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows, total
	}
	return rows, total
}

// Select 查询集合
func (r SysDictType) Select(sysDictType model.SysDictType) []model.SysDictType {
	tx := db.DB("").Model(&model.SysDictType{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if sysDictType.DictName != "" {
		tx = tx.Where("dict_name like concat(?, '%')", sysDictType.DictName)
	}
	if sysDictType.DictType != "" {
		tx = tx.Where("dict_type like concat(?, '%')", sysDictType.DictType)
	}
	if sysDictType.StatusFlag != "" {
		tx = tx.Where("status_flag = ?", sysDictType.StatusFlag)
	}

	// 查询数据
	rows := []model.SysDictType{}
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// SelectByIds 通过ID查询信息
func (r SysDictType) SelectByIds(dictIds []int64) []model.SysDictType {
	rows := []model.SysDictType{}
	if len(dictIds) <= 0 {
		return rows
	}
	tx := db.DB("").Model(&model.SysDictType{})
	// 构建查询条件
	tx = tx.Where("dict_id in ? and del_flag = '0'", dictIds)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// Insert 新增信息 返回新增数据ID
func (r SysDictType) Insert(sysDictType model.SysDictType) int64 {
	sysDictType.DelFlag = "0"
	if sysDictType.CreateBy != "" {
		sysDictType.CreateTime = time.Now().UnixMilli()
	}
	// 执行插入
	if err := db.DB("").Create(&sysDictType).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
	}
	return sysDictType.DictId
}

// Update 修改信息 返回受影响的行数
func (r SysDictType) Update(sysDictType model.SysDictType) int64 {
	if sysDictType.DictId <= 0 {
		return 0
	}
	if sysDictType.UpdateBy != "" {
		sysDictType.UpdateTime = time.Now().UnixMilli()
	}
	tx := db.DB("").Model(&model.SysDictType{})
	// 构建查询条件
	tx = tx.Where("dict_id = ?", sysDictType.DictId)
	// 执行更新
	if err := tx.Updates(sysDictType).Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByIds 批量删除信息 返回受影响的行数
func (r SysDictType) DeleteByIds(dictIds []int64) int64 {
	if len(dictIds) <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.SysDictType{})
	// 构建查询条件
	tx = tx.Where("dict_id in ?", dictIds)
	// 执行更新删除标记
	if err := tx.Update("del_flag", "1").Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// CheckUnique 检查信息是否唯一 返回数据ID
func (r SysDictType) CheckUnique(sysDictType model.SysDictType) int64 {
	tx := db.DB("").Model(&model.SysDictType{})
	tx = tx.Where("del_flag = 0")
	// 查询条件拼接
	if sysDictType.DictName != "" {
		tx = tx.Where("dict_name = ?", sysDictType.DictName)
	}
	if sysDictType.DictType != "" {
		tx = tx.Where("dict_type = ?", sysDictType.DictType)
	}
	// 查询数据
	var id int64 = 0
	if err := tx.Select("dict_id").Limit(1).Find(&id).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return id
	}
	return id
}

// SelectByType 通过字典类型查询信息
func (r SysDictType) SelectByType(dictType string) model.SysDictType {
	item := model.SysDictType{}
	if dictType == "" {
		return item
	}
	tx := db.DB("").Model(&model.SysDictType{})
	tx.Where("dict_type = ? and del_flag = '0'", dictType)
	// 查询数据
	if err := tx.Limit(1).Find(&item).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return item
	}
	return item
}
