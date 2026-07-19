package repository

import (
	"fmt"

	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/modules/system/model"
)

// NewSysConfig 实例化数据层
var NewSysConfig = &SysConfig{}

// SysConfig 参数配置表 数据层处理
type SysConfig struct{}

// sysConfigSortMapping 参数配置排序字段映射
var sysConfigSortMapping = map[string]string{
	"id":           "id",
	"createTime":   "create_time",
	"updateTime":   "update_time",
	"configName":   "config_name",
	"configKey":    "config_key",
	"configGroup":  "config_group",
	"configSort":   "config_sort",
	"configStatus": "config_status",
}

// SelectByPage 分页查询集合
func (r SysConfig) SelectByPage(query map[string]string) ([]model.SysConfig, int64, error) {
	tx := db.DB("").Model(&model.SysConfig{})
	if v, ok := query["configName"]; ok && v != "" {
		tx = tx.Where("config_name like ?", v+"%")
	}
	if v, ok := query["configGroup"]; ok && v != "" {
		tx = tx.Where("config_group = ?", v)
	}
	if v, ok := query["configKey"]; ok && v != "" {
		tx = tx.Where("config_key like ?", v+"%")
	}
	if v, ok := query["beginTime"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s000", v)
			tx = tx.Where("create_time >= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("create_time >= ?", v)
		}
	}
	if v, ok := query["endTime"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s999", v)
			tx = tx.Where("create_time <= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("create_time <= ?", v)
		}
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysConfig{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return rows, total, nil
	}

	// 查询数据分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序
	orders := db.SortParse(query["sort"], sysConfigSortMapping)
	tx = db.SortOrDefault(tx, orders, "config_sort asc")

	err := tx.Find(&rows).Error
	if err != nil {
		return rows, total, err
	}
	return rows, total, nil
}

// Select 查询集合
func (r SysConfig) Select(sysConfig model.SysConfig) ([]model.SysConfig, error) {
	tx := db.DB("").Model(&model.SysConfig{})
	// 查询条件拼接
	if sysConfig.ConfigName != "" {
		tx = tx.Where("config_name like ?", sysConfig.ConfigName+"%")
	}
	if sysConfig.ConfigGroup != "" {
		tx = tx.Where("config_group = ?", sysConfig.ConfigGroup)
	}
	if sysConfig.ConfigKey != "" {
		tx = tx.Where("config_key like ?", sysConfig.ConfigKey+"%")
	}
	if sysConfig.CreateTime > 0 {
		tx = tx.Where("create_time >= ?", sysConfig.CreateTime)
	}

	// 查询数据
	rows := []model.SysConfig{}
	if err := tx.Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// SelectByIds 通过ID查询信息
func (r SysConfig) SelectByIds(ids []int64) ([]model.SysConfig, error) {
	rows := []model.SysConfig{}
	if len(ids) <= 0 {
		return rows, nil
	}
	tx := db.DB("").Model(&model.SysConfig{})
	// 构建查询条件
	tx = tx.Where("id in ?", ids)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// Insert 新增信息 返回新增数据ID
func (r SysConfig) Insert(sysConfig model.SysConfig) (int64, error) {
	// 执行插入
	if err := db.DB("").Create(&sysConfig).Error; err != nil {
		return 0, err
	}
	return sysConfig.ID, nil
}

// Update 修改信息 返回受影响行数
func (r SysConfig) Update(sysConfig model.SysConfig) (int64, error) {
	if sysConfig.ID <= 0 {
		return 0, nil
	}
	tx := db.DB("").Model(&model.SysConfig{})
	// 构建查询条件
	tx = tx.Where("id = ?", sysConfig.ID)
	tx = tx.Omit("id", "create_by", "create_time")
	// 执行更新
	if err := tx.Updates(sysConfig).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// DeleteByIds 批量删除信息 返回受影响行数（物理删除）
func (r SysConfig) DeleteByIds(ids []int64) (int64, error) {
	if len(ids) <= 0 {
		return 0, nil
	}
	tx := db.DB("").Model(&model.SysConfig{})
	// 构建查询条件
	tx = tx.Where("id in ?", ids)
	// 执行物理删除
	if err := tx.Delete(&model.SysConfig{}).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// CheckUnique 检查信息是否唯一 返回数据ID
func (r SysConfig) CheckUnique(sysConfig model.SysConfig) (int64, error) {
	tx := db.DB("").Model(&model.SysConfig{})
	// 查询条件拼接
	if sysConfig.ConfigGroup != "" {
		tx = tx.Where("config_group = ?", sysConfig.ConfigGroup)
	}
	if sysConfig.ConfigKey != "" {
		tx = tx.Where("config_key = ?", sysConfig.ConfigKey)
	}
	// 查询数据
	var id int64 = 0
	if err := tx.Select("id").Limit(1).Find(&id).Error; err != nil {
		return id, err
	}
	return id, nil
}

// SelectValueByKey 通过Key查询Value
func (r SysConfig) SelectValueByKey(configKey string) (string, error) {
	if configKey == "" {
		return "", nil
	}
	tx := db.DB("").Model(&model.SysConfig{})
	tx.Where("config_key = ?", configKey)
	// 查询数据
	var configValue string = ""
	if err := tx.Select("config_value").Limit(1).Find(&configValue).Error; err != nil {
		return configValue, err
	}
	return configValue, nil
}
