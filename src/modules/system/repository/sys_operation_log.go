package repository

import (
	"fmt"
	"strings"

	"mask_api_gin/src/internal/infra/database"
	"mask_api_gin/src/modules/system/model"
)

// NewSysOperationLog 实例化数据层
var NewSysOperationLog = &SysOperationLog{}

// SysOperationLog 操作日志表 数据层处理
type SysOperationLog struct{}

// sysOperationLogSortMapping 操作日志排序字段映射
var sysOperationLogSortMapping = map[string]string{
	"id":                  "id",
	"create_by":           "create_by",
	"create_time":         "create_time",
	"title":               "title",
	"operation_type":      "operation_type",
	"operation_cost_time": "operation_cost_time",
	"operation_status":    "operation_status",
}

// SelectByPage 分页查询集合
func (r SysOperationLog) SelectByPage(query map[string]string, dataScopeSQL string) ([]model.SysOperationLog, int64, error) {
	tx := database.DB("").Model(&model.SysOperationLog{})
	// 查询条件拼接
	if v, ok := query["title"]; ok && v != "" {
		tx = tx.Where("title like ?", v+"%")
	}
	if v, ok := query["operation_type"]; ok && v != "" {
		operationTypes := strings.Split(v, ",")
		tx = tx.Where("operation_type in ?", operationTypes)
	}
	if v, ok := query["create_by"]; ok && v != "" {
		tx = tx.Where("create_by like ?", v+"%")
	}
	if v, ok := query["operation_ip"]; ok && v != "" {
		tx = tx.Where("operation_ip like ?", v+"%")
	}
	if v, ok := query["operation_status"]; ok && v != "" {
		tx = tx.Where("operation_status = ?", v)
	}
	if v, ok := query["create_time_begin"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s000", v)
			tx = tx.Where("create_time >= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("create_time >= ?", v)
		}
	}
	if v, ok := query["create_time_end"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s999", v)
			tx = tx.Where("create_time <= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("create_time <= ?", v)
		}
	}
	if dataScopeSQL != "" {
		dataScopeSQL = fmt.Sprintf("select distinct user_name from sys_user where %s", dataScopeSQL)
		tx = tx.Where(fmt.Sprintf("create_by in ( %s )", dataScopeSQL))
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysOperationLog{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil {
		return rows, total, err
	}
	if total <= 0 {
		return rows, total, nil
	}

	// 查询数据分页
	pageNum, pageSize := database.PageNumSize(query["page_num"], query["page_size"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序
	orders := database.SortParse(query["sort"], sysOperationLogSortMapping)
	tx = database.SortOrDefault(tx, orders, "id desc")

	// 查询数据
	err := tx.Find(&rows).Error
	if err != nil {
		return rows, total, err
	}
	return rows, total, nil
}

// Insert 新增信息
func (r SysOperationLog) Insert(sysOperationLog model.SysOperationLog) (int64, error) {
	// 执行插入
	if err := database.DB("").Create(&sysOperationLog).Error; err != nil {
		return 0, err
	}
	return sysOperationLog.ID, nil
}

// Clean 清空信息
func (r SysOperationLog) Clean() (int64, error) {
	tx := database.DB("").Where("1 = 1").Delete(&model.SysOperationLog{})
	if err := tx.Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}
