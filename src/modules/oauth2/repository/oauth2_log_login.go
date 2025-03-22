package repository

import (
	"fmt"
	"strings"
	"time"

	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/oauth2/model"
)

// NewOauth2LogLogin 实例化数据层
var NewOauth2LogLogin = &Oauth2LogLoginRepository{}

// Oauth2LogLoginRepository 用户授权第三方应用登录日志表 数据层处理
type Oauth2LogLoginRepository struct{}

// SelectByPage 分页查询集合
func (r Oauth2LogLoginRepository) SelectByPage(query map[string]string) ([]model.Oauth2LogLogin, int64) {
	tx := db.DB("").Model(&model.Oauth2LogLogin{})
	// 查询条件拼接
	if v, ok := query["loginIp"]; ok && v != "" {
		tx = tx.Where("login_ip like ?", v+"%")
	}
	if v, ok := query["clientId"]; ok && v != "" {
		tx = tx.Where("client_id = ?", v)
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
	rows := []model.Oauth2LogLogin{}

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
			sortOrder := sortOrderArr[i]
			// 排序字段
			sort := "id"
			if sortBy == "loginIp" {
				sort = "login_ip"
			} else if sortBy == "createTime" {
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
func (r Oauth2LogLoginRepository) Insert(param model.Oauth2LogLogin) int64 {
	param.LoginTime = time.Now().UnixMilli()
	// 执行插入
	if err := db.DB("").Create(&param).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
	}
	return param.ID
}

// Clean 清空信息
func (r Oauth2LogLoginRepository) Clean() int64 {
	tx := db.DB("").Delete(&model.Oauth2LogLogin{})
	if err := tx.Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}
