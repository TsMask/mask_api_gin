package repository

import (
	"mask_api_gin/src/internal/infra/database"
	"mask_api_gin/src/modules/demo/model"
)

// NewDemoORM 实例化数据层
var NewDemoORM = &DemoORM{}

// DemoORM 测试ORM表 数据层处理
type DemoORM struct{}

// demoORMSortMapping 测试ORM排序字段映射
var demoORMSortMapping = map[string]string{
	"id":          "id",
	"create_time": "create_time",
	"update_time": "update_time",
	"title":       "title",
	"orm_type":    "orm_type",
	"orm_status":  "orm_status",
}

// SelectByPage 分页查询
func (r DemoORM) SelectByPage(query map[string]string) ([]model.DemoORM, int64, error) {
	tx := database.DB("").Model(&model.DemoORM{})
	if v, ok := query["title"]; ok && v != "" {
		tx = tx.Where("title like ?", v+"%")
	}
	if v, ok := query["orm_status"]; ok && v != "" {
		tx = tx.Where("orm_status = ?", v)
	}

	var total int64 = 0
	if err := tx.Count(&total).Error; err != nil {
		return nil, total, err
	}
	if total <= 0 {
		return nil, total, nil
	}

	rows := []model.DemoORM{}
	pageNum, pageSize := database.PageNumSize(query["page_num"], query["page_size"])
	tx = tx.Limit(pageSize).Offset(pageSize * pageNum)

	// 排序
	orders := database.SortParse(query["sort"], demoORMSortMapping)
	tx = database.SortOrDefault(tx, orders, "id desc")

	if err := tx.Find(&rows).Error; err != nil {
		return rows, total, err
	}
	return rows, total, nil
}

// Select 查询集合
func (r DemoORM) Select(demoORM model.DemoORM) ([]model.DemoORM, error) {
	tx := database.DB("").Model(&model.DemoORM{})
	if demoORM.Title != "" {
		tx = tx.Where("title like ?", demoORM.Title+"%")
	}
	if demoORM.OrmStatus != "" {
		tx = tx.Where("orm_status = ?", demoORM.OrmStatus)
	}

	rows := []model.DemoORM{}
	if err := tx.Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

// SelectById 通过ID查询
func (r DemoORM) SelectById(id int64) (model.DemoORM, error) {
	item := model.DemoORM{}
	if id <= 0 {
		return item, nil
	}
	tx := database.DB("").Model(&model.DemoORM{})
	tx = tx.Where("id = ?", id)
	if err := tx.Limit(1).Find(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

// Insert 新增
func (r DemoORM) Insert(demoORM model.DemoORM) (int64, error) {
	if err := database.DB("").Create(&demoORM).Error; err != nil {
		return 0, err
	}
	return demoORM.Id, nil
}

// Update 更新
func (r DemoORM) Update(demoORM model.DemoORM) (int64, error) {
	if demoORM.Id <= 0 {
		return 0, nil
	}
	tx := database.DB("").Model(&model.DemoORM{})
	tx = tx.Where("id = ?", demoORM.Id)
	if err := tx.Omit("id", "create_by", "create_time").Updates(demoORM).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// DeleteByIds 批量删除
func (r DemoORM) DeleteByIds(ids []int64) (int64, error) {
	if len(ids) <= 0 {
		return 0, nil
	}
	tx := database.DB("").Where("id in ?", ids)
	if err := tx.Delete(&model.DemoORM{}).Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}

// Clean 清空表
func (r DemoORM) Clean() (int64, error) {
	tx := database.DB("").Where("1 = 1").Delete(&model.DemoORM{})
	if err := tx.Error; err != nil {
		return 0, err
	}
	return tx.RowsAffected, nil
}
