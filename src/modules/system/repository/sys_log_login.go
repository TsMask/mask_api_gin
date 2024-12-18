package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"

	"fmt"
	"time"
)

// NewSysLogLogin 实例化数据层
var NewSysLogLogin = &SysLogLogin{}

// SysLogLoginRepository 系统登录访问表 数据层处理
type SysLogLogin struct{}

// SelectByPage 分页查询集合
func (r SysLogLogin) SelectByPage(query map[string]string) ([]model.SysLogLogin, int64) {
	tx := db.DB("").Model(&model.SysLogLogin{})
	// 查询条件拼接
	if v, ok := query["loginIp"]; ok && v != "" {
		tx = tx.Where("login_ip like concat(?, '%')", v)
	}
	if v, ok := query["userName"]; ok && v != "" {
		tx = tx.Where("user_name like concat(?, '%')", v)
	}
	if v, ok := query["statusFlag"]; ok && v != "" {
		tx = tx.Where("status_flag = ?", v)
	}
	if v, ok := query["beginTime"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s000", v)
			tx = tx.Where("login_time >= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("login_time >= ?", v)
		}
	}
	if v, ok := query["endTime"]; ok && v != "" {
		if len(v) == 10 {
			v = fmt.Sprintf("%s999", v)
			tx = tx.Where("login_time <= ?", v)
		} else if len(v) == 13 {
			tx = tx.Where("login_time <= ?", v)
		}
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysLogLogin{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return rows, total
	}

	// 查询数据分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)
	err := tx.Order("id desc").Find(&rows).Error
	if err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows, total
	}
	return rows, total
}

// Insert 新增信息 返回新增的数据ID
func (r SysLogLogin) Insert(sysLogLogin model.SysLogLogin) string {
	sysLogLogin.LoginTime = time.Now().UnixMilli()
	// 执行插入
	if err := db.DB("").Create(&sysLogLogin).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return ""
	}
	return sysLogLogin.ID
}

// Clean 清空信息
func (r SysLogLogin) Clean() int64 {
	tx := db.DB("").Delete(&model.SysLogLogin{})
	if err := tx.Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}
