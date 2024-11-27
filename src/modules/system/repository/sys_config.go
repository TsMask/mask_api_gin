package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"

	"fmt"
	"time"
)

// NewSysConfig 实例化数据层
var NewSysConfig = &SysConfig{}

// SysConfig 参数配置表 数据层处理
type SysConfig struct{}

// SelectByPage 分页查询集合
func (r SysConfig) SelectByPage(query map[string]string) ([]model.SysConfig, int64) {
	tx := db.DB("").Model(&model.SysConfig{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if v, ok := query["configName"]; ok && v != "" {
		tx = tx.Where("config_name like concat(?, '%')", v)
	}
	if v, ok := query["configType"]; ok && v != "" {
		tx = tx.Where("config_type = ?", v)
	}
	if v, ok := query["configKey"]; ok && v != "" {
		tx = tx.Where("config_key like concat(?, '%')", v)
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
			v = fmt.Sprintf("%s000", v)
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
func (r SysConfig) Select(sysConfig model.SysConfig) []model.SysConfig {
	tx := db.DB("").Model(&model.SysConfig{})
	tx = tx.Where("del_flag = '0'")
	// 查询条件拼接
	if sysConfig.ConfigName != "" {
		tx = tx.Where("config_name like concat(?, '%')", sysConfig.ConfigName)
	}
	if sysConfig.ConfigType != "" {
		tx = tx.Where("config_type = ?", sysConfig.ConfigType)
	}
	if sysConfig.ConfigKey != "" {
		tx = tx.Where("config_key like concat(?, '%')", sysConfig.ConfigKey)
	}
	if sysConfig.CreateTime > 0 {
		tx = tx.Where("create_time >= ?", sysConfig.CreateTime)
	}

	// 查询数据
	rows := []model.SysConfig{}
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// SelectByIds 通过ID查询信息
func (r SysConfig) SelectByIds(configIds []string) []model.SysConfig {
	rows := []model.SysConfig{}
	if len(configIds) <= 0 {
		return rows
	}
	tx := db.DB("").Model(&model.SysConfig{})
	// 构建查询条件
	tx = tx.Where("config_id in ? and del_flag = '0'", configIds)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// Insert 新增信息 返回新增数据ID
func (r SysConfig) Insert(sysConfig model.SysConfig) string {
	sysConfig.DelFlag = "0"
	if sysConfig.CreateBy != "" {
		ms := time.Now().UnixMilli()
		sysConfig.UpdateBy = sysConfig.CreateBy
		sysConfig.UpdateTime = ms
		sysConfig.CreateTime = ms
	}
	// 执行插入
	if err := db.DB("").Create(&sysConfig).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return ""
	}
	return sysConfig.ConfigId
}

// Update 修改信息 返回受影响行数
func (r SysConfig) Update(sysConfig model.SysConfig) int64 {
	if sysConfig.ConfigId == "" {
		return 0
	}
	if sysConfig.UpdateBy != "" {
		sysConfig.UpdateTime = time.Now().UnixMilli()
	}
	tx := db.DB("").Model(&model.SysConfig{})
	// 构建查询条件
	tx = tx.Where("config_id = ?", sysConfig.ConfigId)
	tx = tx.Omit("config_id", "del_flag", "create_by", "create_time")
	// 执行更新
	if err := tx.Updates(sysConfig).Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByIds 批量删除信息 返回受影响行数
func (r SysConfig) DeleteByIds(configIds []string) int64 {
	if len(configIds) <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.SysConfig{})
	// 构建查询条件
	tx = tx.Where("config_id in ?", configIds)
	// 执行更新删除标记
	if err := tx.Update("del_flag", "1").Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// CheckUnique 检查信息是否唯一 返回数据ID
func (r SysConfig) CheckUnique(sysConfig model.SysConfig) string {
	tx := db.DB("").Model(&model.SysConfig{})
	tx = tx.Where("del_flag = 0")
	// 查询条件拼接
	if sysConfig.ConfigType != "" {
		tx = tx.Where("config_type = ?", sysConfig.ConfigType)
	}
	if sysConfig.ConfigKey != "" {
		tx = tx.Where("config_key = ?", sysConfig.ConfigKey)
	}
	// 查询数据
	var id string = ""
	if err := tx.Select("config_id").Limit(1).Find(&id).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return id
	}
	return id
}

// SelectValueByKey 通过Key查询Value
func (r SysConfig) SelectValueByKey(configKey string) string {
	if configKey == "" {
		return ""
	}
	tx := db.DB("").Model(&model.SysConfig{})
	tx.Where("config_key = ? and del_flag = '0'", configKey)
	// 查询数据
	var configValue string = ""
	if err := tx.Select("config_value").Limit(1).Find(&configValue).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return configValue
	}
	return configValue
}
