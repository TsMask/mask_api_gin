package repository

import (
	"fmt"
	"strings"
	"time"

	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/system/model"
)

// NewSysLogLogin 实例化数据层
var NewSysLogLogin = &SysLogLogin{}

// SysLogLoginRepository 系统登录访问表 数据层处理
type SysLogLogin struct{}

// SelectByPage 分页查询集合
func (r SysLogLogin) SelectByPage(query map[string]string, dataScopeSQL string) ([]model.SysLogLogin, int64) {
	tx := db.DB("").Model(&model.SysLogLogin{})
	// 查询条件拼接
	if v, ok := query["loginIp"]; ok && v != "" {
		tx = tx.Where("login_ip like ?", v+"%")
	}
	if v, ok := query["userName"]; ok && v != "" {
		tx = tx.Where("user_name like ?", v+"%")
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
	if dataScopeSQL != "" {
		dataScopeSQL = fmt.Sprintf("select distinct user_name from sys_user where %s", dataScopeSQL)
		tx = tx.Where(fmt.Sprintf("user_name in ( %s )", dataScopeSQL))
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
	// 排序
	sortByStr, sortOk := query["sortBy"]
	sortOrderStr, orderOk := query["sortOrder"]
	if sortOk && sortByStr != "" && orderOk && sortOrderStr != "" {
		sortByArr := strings.Split(sortByStr, ",")
		sortOrderArr := strings.Split(sortOrderStr, ",")
		for i := range sortByArr {
			sortBy := sortByArr[i]
			sortOrder := ""
			if i < len(sortOrderArr) {
				sortOrder = sortOrderArr[i]
			}
			// 排序字段
			sort := "id"
			switch sortBy {
			case "loginIp":
				sort = "login_ip"
			case "userName":
				sort = "user_name"
			case "createTime":
				sort = "create_time"
			}
			// 排序方式
			order := "ASC"
			if strings.HasPrefix(sortOrder, "asc") {
				order = "ASC"
			} else if strings.HasPrefix(sortOrder, "desc") {
				order = "DESC"
			}
			tx = tx.Order(fmt.Sprintf("%s %s", sort, order))
		}
	} else {
		tx = tx.Order("id desc")
	}
	// 查询数据
	err := tx.Find(&rows).Error
	if err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows, total
	}
	return rows, total
}

// Insert 新增信息 返回新增的数据ID
func (r SysLogLogin) Insert(sysLogLogin model.SysLogLogin) int64 {
	sysLogLogin.LoginTime = time.Now().UnixMilli()
	// 执行插入
	if err := db.DB("").Create(&sysLogLogin).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
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
