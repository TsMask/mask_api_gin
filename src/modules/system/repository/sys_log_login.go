package repository

import (
	"fmt"

	"mask_api_gin/src/internal/infra/database"
	"mask_api_gin/src/modules/system/model"
)

// NewSysLogLogin 实例化数据层
var NewSysLogLogin = &SysLogLogin{}

// SysLogLoginRepository 系统登录访问表 数据层处理
type SysLogLogin struct{}

// sysLogLoginSortMapping 登录日志排序字段映射
var sysLogLoginSortMapping = map[string]string{
	"id":              "id",
	"create_time":     "create_time",
	"login_ip":        "login_ip",
	"login_user_name": "login_user_name",
	"login_status":    "login_status",
}

// SelectByPage 分页查询集合
func (r SysLogLogin) SelectByPage(query map[string]string, dataScopeSQL string) ([]model.SysLogLogin, int64, error) {
	tx := database.DB("").Model(&model.SysLogLogin{})
	// 查询条件拼接
	if v, ok := query["login_ip"]; ok && v != "" {
		tx = tx.Where("login_ip like ?", v+"%")
	}
	if v, ok := query["login_user_name"]; ok && v != "" {
		tx = tx.Where("login_user_name like ?", v+"%")
	}
	if v, ok := query["login_status"]; ok && v != "" {
		tx = tx.Where("login_status = ?", v)
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
		tx = tx.Where(fmt.Sprintf("login_user_name in ( %s )", dataScopeSQL))
	}

	// 查询结果
	var total int64 = 0
	rows := []model.SysLogLogin{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return rows, total, nil
	}

	// 查询数据分页
	pageNum, pageSize := database.PageNumSize(query["page_num"], query["page_size"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序
	orders := database.SortParse(query["sort"], sysLogLoginSortMapping)
	tx = database.SortOrDefault(tx, orders, "id desc")

	// 查询数据
	err := tx.Find(&rows).Error
	if err != nil {
		return rows, total, err
	}
	return rows, total, nil
}

// Insert 新增信息 返回新增的数据ID
func (r SysLogLogin) Insert(sysLogLogin model.SysLogLogin) (int64, error) {
	// 执行插入
	if err := database.DB("").Create(&sysLogLogin).Error; err != nil {
		return 0, err
	}
	return sysLogLogin.ID, nil
}

// Clean 清空信息
func (r SysLogLogin) Clean() (int64, error) {
	tx := database.DB("").Where("1 = 1").Delete(&model.SysLogLogin{})
	if err := tx.Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}
