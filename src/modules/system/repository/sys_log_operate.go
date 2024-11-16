package repository

import (
	"fmt"
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/framework/utils/date"
	"mask_api_gin/src/modules/system/model"
	"time"
)

// NewSysLogOperate 实例化数据层
var NewSysLogOperate = &SysLogOperate{}

// SysLogOperateRepository 操作日志表 数据层处理
type SysLogOperate struct{}

// SelectByPage 分页查询集合
func (r SysLogOperate) SelectByPage(query map[string]any) ([]model.SysLogOperate, int64) {
	tx := db.DB("").Model(&model.SysLogOperate{})
	// 查询条件拼接
	if v, ok := query["title"]; ok && v != "" {
		tx = tx.Where("title like concat(?, '%')", v)
	}
	if v, ok := query["businessType"]; ok && v != "" {
		tx = tx.Where("business_type = ?", v)
	}
	if v, ok := query["operaBy"]; ok && v != "" {
		tx = tx.Where("opera_by like concat(?, '%')", v)
	}
	if v, ok := query["operaIp"]; ok && v != "" {
		tx = tx.Where("opera_ip like concat(?, '%')", v)
	}
	if v, ok := query["statusFlag"]; ok && v != "" {
		tx = tx.Where("status_flag = ?", v)
	}
	if v, ok := query["beginTime"]; ok && v != "" {
		tx = tx.Where("opera_time >= ?", v)
	}
	if v, ok := query["endTime"]; ok && v != "" {
		tx = tx.Where("opera_time <= ?", v)
	}
	if v, ok := query["params[beginTime]"]; ok && v != "" {
		beginDate := date.ParseStrToDate(fmt.Sprint(v), date.YYYY_MM_DD)
		tx = tx.Where("opera_time >= ?", beginDate.UnixMilli())
	}
	if v, ok := query["params[endTime]"]; ok && v != "" {
		endDate := date.ParseStrToDate(fmt.Sprint(v), date.YYYY_MM_DD)
		tx = tx.Where("opera_time <= ?", endDate.UnixMilli())
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
	err := tx.Limit(pageSize).Offset(pageSize * pageNum).Find(&rows).Error
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
func (r SysLogOperate) Clean() error {
	sql := "truncate table sys_log_operate"
	_, err := db.ExecDB("", sql, []any{})
	return err
}
