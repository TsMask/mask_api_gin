package repository

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/demo/model"
)

// NewDemoORM 实例化数据层
var NewDemoORM = &DemoORM{}

// DemoORM 测试ORM表 数据层处理
type DemoORM struct{}

// demoORMSortMapping 测试ORM排序字段映射
var demoORMSortMapping = map[string]string{
	"id":         "id",
	"title":      "title",
	"ormType":    "orm_type",
	"statusFlag": "status_flag",
	"createBy":   "create_by",
	"createTime": "create_time",
	"updateTime": "update_time",
}

// SelectByPage 分页查询
func (r DemoORM) SelectByPage(query map[string]string) ([]model.DemoORM, int64) {
	tx := db.DB("").Model(&model.DemoORM{})
	if v, ok := query["title"]; ok && v != "" {
		tx = tx.Where("title like ?", v+"%")
	}
	if v, ok := query["statusFlag"]; ok && v != "" {
		tx = tx.Where("status_flag = ?", v)
	}

	var total int64 = 0
	if err := tx.Count(&total).Error; err != nil || total <= 0 {
		return []model.DemoORM{}, total
	}

	rows := []model.DemoORM{}
	pageNum, pageSize := db.PageNumSize(query["pageNum"], query["pageSize"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序
	orders := db.SortParse(query["sort"], demoORMSortMapping)
	tx = db.SortOrDefault(tx, orders, "id desc")

	if err := tx.Find(&rows).Error; err != nil {
		return rows, total
	}
	return rows, total
}

// Select 查询集合
func (r DemoORM) Select(demoORM model.DemoORM) []model.DemoORM {
	tx := db.DB("").Model(&model.DemoORM{})
	if demoORM.Title != "" {
		tx = tx.Where("title like ?", demoORM.Title+"%")
	}
	if demoORM.StatusFlag != "" {
		tx = tx.Where("status_flag = ?", demoORM.StatusFlag)
	}

	rows := []model.DemoORM{}
	if err := tx.Find(&rows).Error; err != nil {
		return rows
	}
	return rows
}

// SelectById 通过ID查询
func (r DemoORM) SelectById(id int64) model.DemoORM {
	item := model.DemoORM{}
	if id <= 0 {
		return item
	}
	tx := db.DB("").Model(&model.DemoORM{})
	tx = tx.Where("id = ?", id)
	if err := tx.Limit(1).Find(&item).Error; err != nil {
		logger.Errorf("query find err => %v", err.Error())
	}
	return item
}

// Insert 新增
func (r DemoORM) Insert(demoORM model.DemoORM) int64 {
	if err := db.DB("").Create(&demoORM).Error; err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
	}
	return demoORM.Id
}

// Update 更新
func (r DemoORM) Update(demoORM model.DemoORM) int64 {
	if demoORM.Id <= 0 {
		return 0
	}
	tx := db.DB("").Model(&model.DemoORM{})
	tx = tx.Where("id = ?", demoORM.Id)
	if err := tx.Omit("id", "create_by", "create_time").Updates(demoORM).Error; err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// DeleteByIds 批量删除
func (r DemoORM) DeleteByIds(ids []int64) int64 {
	if len(ids) <= 0 {
		return 0
	}
	tx := db.DB("").Where("id in ?", ids)
	if err := tx.Delete(&model.DemoORM{}).Error; err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}

// Clean 清空表
func (r DemoORM) Clean() int64 {
	tx := db.DB("").Where("1 = 1").Delete(&model.DemoORM{})
	if err := tx.Error; err != nil {
		logger.Errorf("clean err => %v", err.Error())
		return 0
	}
	return tx.RowsAffected
}
