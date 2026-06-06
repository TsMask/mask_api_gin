package repository

import (
	"fmt"
	"time"

	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"
)

// NewSysLogOperate 实例化数据层
var NewSysLogOperate = &SysLogOperate{}

// SysLogOperate 操作日志表 数据层处理
type SysLogOperate struct{}

// sysLogOperateSortMapping 操作日志排序字段映射
var sysLogOperateSortMapping = map[string]string{
	"id":           "id",
	"title":        "title",
	"businessType": "business_type",
	"operaBy":      "opera_by",
	"operaTime":    "opera_time",
	"costTime":     "cost_time",
	"statusFlag":   "status_flag",
}

// SelectByPage 分页查询集合
func (r SysLogOperate) SelectByPage(query map[string]string, dataScopeSQL string) ([]model.SysLogOperate, int64) {
	tx := db.DB("").Model(&model.SysLogOperate{})
	// 查询条件拼接
	if v, ok := query["title"]; ok && v != "" {
		tx = tx.Where("title like ?", v+"%")
	}
	if v, ok := query["businessType"]; ok && v != "" {
		tx = tx.Where("business_type = ?", v)
	}
	if v, ok := query["operaBy"]; ok && v != "" {
		tx = tx.Where("opera_by like ?", v+"%")
	}
	if v, ok := query["operaIp"]; ok && v != "" {
		tx = tx.Where("opera_ip like ?", v+"%")
	}
	if v, ok := query["statusFlag"]; ok && v != "" {
		tx = tx.Where("status_flag = ?", v)
	}
	if v, ok := query["beginTime"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s000", v)
			tx = tx.Where("opera_time >= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("opera_time >= ?", v)
		}
	}
	if v, ok := query["endTime"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s999", v)
			tx = tx.Where("opera_time <= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("opera_time <= ?", v)
		}
	}
	if dataScopeSQL != "" {
		dataScopeSQL = fmt.Sprintf("select distinct user_name from sys_user where %s", dataScopeSQL)
		tx = tx.Where(fmt.Sprintf("opera_by in ( %s )", dataScopeSQL))
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysLogOperate{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return rows, total
	}

	// 查询数据分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序
	orders := db.SortParse(query["sort"], sysLogOperateSortMapping)
	tx = db.SortOrDefault(tx, orders, "id desc")

	// 查询数据
	err := tx.Find(&rows).Error
	if err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows, total
	}
	return rows, total
}

// Insert 新增信息
func (r SysLogOperate) Insert(sysLogOperate model.SysLogOperate) int64 {
	if sysLogOperate.OperaBy != "" {
		sysLogOperate.OperaTime = time.Now().UnixMilli()
	}
	// 执行插入
	if err := db.DB("").Create(&sysLogOperate).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
	}
	return sysLogOperate.ID
}

// Clean 清空信息
func (r SysLogOperate) Clean() int64 {
	tx := db.DB("").Where("1 = 1").Delete(&model.SysLogOperate{})
	if err := tx.Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}
