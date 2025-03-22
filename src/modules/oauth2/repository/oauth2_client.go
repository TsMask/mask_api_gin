package repository

import (
	"time"

	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/oauth2/model"
)

// NewOauth2Client 实例化数据层
var NewOauth2Client = &Oauth2Client{}

// Oauth2Client 用户授权第三方应用表 数据层处理
type Oauth2Client struct{}

// SelectByPage 分页查询集合
func (r Oauth2Client) SelectByPage(query map[string]string) ([]model.Oauth2Client, int64) {
	tx := db.DB("").Model(&model.Oauth2Client{})
	tx = tx.Where("del_flag = 0")
	// 查询条件拼接
	if v, ok := query["clientId"]; ok && v != "" {
		tx = tx.Where("client_id = ?", v)
	}
	if v, ok := query["title"]; ok && v != "" {
		tx = tx.Where("title like ?", v+"%")
	}

	// 查询结果
	var total int64 = 0
	rows := []model.Oauth2Client{}

	// 查询数量为0直接返回
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return rows, total
	}

	// 查询数据分页
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)
	err := tx.Find(&rows).Error
	if err != nil {
		return rows, total
	}
	return rows, total
}

// Select 查询集合
func (r Oauth2Client) Select(param model.Oauth2Client) []model.Oauth2Client {
	tx := db.DB("").Model(&model.Oauth2Client{})
	// 查询条件拼接
	if param.ClientId != "" {
		tx = tx.Where("client_id = ?", param.ClientId)
	}
	if param.Title != "" {
		tx = tx.Where("title like ?", param.Title+"%")
	}

	// 查询数据
	rows := []model.Oauth2Client{}
	if err := tx.Find(&rows).Error; err != nil {
		return rows
	}
	return rows
}

// SelectByIds 通过ID查询信息
func (r Oauth2Client) SelectByIds(ids []int64) []model.Oauth2Client {
	rows := []model.Oauth2Client{}
	if len(ids) <= 0 {
		return rows
	}
	tx := db.DB("").Model(&model.Oauth2Client{})
	tx = tx.Where("del_flag = 0")
	// 构建查询条件
	tx = tx.Where("id in ?", ids)
	// 查询数据
	if err := tx.Find(&rows).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
		return rows
	}
	return rows
}

// Insert 新增信息 返回新增数据ID
func (r Oauth2Client) Insert(param model.Oauth2Client) int64 {
	if param.CreateBy != "" {
		ms := time.Now().UnixMilli()
		param.UpdateBy = param.CreateBy
		param.UpdateTime = ms
		param.CreateTime = ms
		param.DelFlag = "0"
	}
	// 执行插入
	if err := db.DB("").Create(&param).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
	}
	return param.Id
}

// Update 修改信息 返回受影响行数
func (r Oauth2Client) Update(param model.Oauth2Client) int64 {
	if param.Id <= 0 {
		return 0
	}
	if param.UpdateBy != "" {
		param.UpdateTime = time.Now().UnixMilli()
	}
	tx := db.DB("").Model(&model.Oauth2Client{})
	// 构建查询条件
	tx = tx.Where("id = ?", param.Id)
	tx = tx.Omit("id", "create_by", "create_time")
	// 执行更新
	if err := tx.Updates(param).Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByIds 批量删除信息 返回受影响行数
func (r Oauth2Client) DeleteByIds(ids []int64) int64 {
	if len(ids) <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.Oauth2Client{})
	// 构建查询条件
	tx = tx.Where("id in ?", ids)
	// 执行更新删除标记
	if err := tx.Update("del_flag", "1").Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// SelectByClientId 通过clientId查询
func (r Oauth2Client) SelectByClientId(clientId string) model.Oauth2Client {
	item := model.Oauth2Client{}
	if clientId == "" {
		return item
	}
	tx := db.DB("").Model(&model.Oauth2Client{})
	tx = tx.Where("del_flag = 0")
	// 构建查询条件
	tx = tx.Where("client_id = ?", clientId)
	// 查询数据
	if err := tx.Find(&item).Error; err != nil {
		return item
	}
	return item
}
