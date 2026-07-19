package repository

import (
	"fmt"

	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/modules/system/model"
)

// NewSysLogOperate 实例化数据层
var NewSysLogOperate = &SysLogOperate{}

// SysLogOperate 操作日志表 数据层处理
type SysLogOperate struct{}

// sysLogOperateSortMapping 操作日志排序字段映射
var sysLogOperateSortMapping = map[string]string{
	"id":              "id",
	"createBy":        "create_by",
	"createTime":      "create_time",
	"title":           "title",
	"operateType":     "operate_type",
	"operateCostTime": "operate_cost_time",
	"operateStatus":   "operate_status",
}

// SelectByPage 分页查询集合
func (r SysLogOperate) SelectByPage(query map[string]string, dataScopeSQL string) ([]model.SysLogOperate, int64, error) {
	tx := db.DB("").Model(&model.SysLogOperate{})
	// 查询条件拼接
	if v, ok := query["title"]; ok && v != "" {
		tx = tx.Where("title like ?", v+"%")
	}
	if v, ok := query["operateType"]; ok && v != "" {
		tx = tx.Where("operate_type = ?", v)
	}
	if v, ok := query["createBy"]; ok && v != "" {
		tx = tx.Where("create_by like ?", v+"%")
	}
	if v, ok := query["operateIp"]; ok && v != "" {
		tx = tx.Where("operate_ip like ?", v+"%")
	}
	if v, ok := query["operateStatus"]; ok && v != "" {
		tx = tx.Where("operate_status = ?", v)
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
	if dataScopeSQL != "" {
		dataScopeSQL = fmt.Sprintf("select distinct user_name from sys_user where %s", dataScopeSQL)
		tx = tx.Where(fmt.Sprintf("create_by in ( %s )", dataScopeSQL))
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysLogOperate{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return rows, total, nil
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
		return rows, total, err
	}
	return rows, total, nil
}

// Insert 新增信息
func (r SysLogOperate) Insert(sysLogOperate model.SysLogOperate) (int64, error) {
	// 执行插入
	if err := db.DB("").Create(&sysLogOperate).Error; err != nil {
		return 0, err
	}
	return sysLogOperate.ID, nil
}

// Clean 清空信息
func (r SysLogOperate) Clean() (int64, error) {
	tx := db.DB("").Where("1 = 1").Delete(&model.SysLogOperate{})
	if err := tx.Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}
