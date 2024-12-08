package service

import (
	"mask_api_gin/src/framework/database/db"
	"mask_api_gin/src/modules/demo/model"
	"time"
)

// NewDemoORMService 实例化服务层
// https://gorm.io/zh_CN/docs/query.html
var NewDemoORMService = &DemoORMService{}

// DemoORMService 测试ORM信息 服务层处理
type DemoORMService struct{}

// FindByPage 分页查询
func (s DemoORMService) FindByPage(query map[string]string) ([]model.DemoORM, int64) {
	tx := db.DB("").Model(&model.DemoORM{})
	// 查询条件拼接
	if v, ok := query["title"]; ok && v != "" {
		tx = tx.Where("title like concat(?, '%')", v)
	}
	if v, ok := query["statusFlag"]; ok && v != "" {
		tx = tx.Where("status_flag = ?", v)
	}

	// 查询结果
	var total int64 = 0
	rows := []model.DemoORM{}

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

// Find 查询集合
func (s DemoORMService) Find(demoORM model.DemoORM) []model.DemoORM {
	tx := db.DB("").Model(&model.DemoORM{})
	// 查询条件拼接
	if demoORM.Title != "" {
		tx = tx.Where("title like concat(?, '%')", demoORM.Title)
	}
	if demoORM.StatusFlag != "" {
		tx = tx.Where("status_flag = ?", demoORM.StatusFlag)
	}

	// 查询数据
	rows := []model.DemoORM{}
	if err := tx.Find(&rows).Error; err != nil {
		return rows
	}
	return rows
}

// FindById 通过ID查询
func (s DemoORMService) FindById(id string) model.DemoORM {
	item := model.DemoORM{}
	if id == "" {
		return item
	}
	tx := db.DB("").Model(&model.DemoORM{})
	// 构建查询条件
	tx = tx.Where("id = ?", id)
	// 查询数据
	if err := tx.Find(&item).Error; err != nil {
		return item
	}
	return item
}

// Insert 新增
func (s DemoORMService) Insert(demoORM model.DemoORM) string {
	demoORM.CreateBy = "system"
	demoORM.CreateTime = time.Now().UnixMilli()
	demoORM.UpdateBy = demoORM.CreateBy
	demoORM.UpdateTime = demoORM.CreateTime
	// 执行插入
	if err := db.DB("").Create(&demoORM).Error; err != nil {
		return ""
	}
	return demoORM.Id
}

// Update 更新
func (s DemoORMService) Update(demoORM model.DemoORM) int64 {
	if demoORM.Id == "" {
		return 0
	}
	// 查询数据
	var item model.DemoORM
	err := db.DB("").First(&item, demoORM.Id).Error
	if err != nil {
		return 0
	}

	// 只改某些属性
	item.Title = demoORM.Title
	item.OrmType = demoORM.OrmType
	item.StatusFlag = demoORM.StatusFlag
	item.Remark = demoORM.Remark
	item.UpdateBy = "system"
	item.UpdateTime = time.Now().UnixMilli()
	tx := db.DB("").Model(&model.DemoORM{})
	// 构建查询条件
	tx = tx.Where("id = ?", item.Id)
	// 执行更新
	if err := tx.Omit("id", "create_by", "create_time").Updates(item).Error; err != nil {
		return 0
	}
	return tx.RowsAffected
}

// DeleteByIds 批量删除
func (s DemoORMService) DeleteByIds(ids []string) int64 {
	if len(ids) <= 0 {
		return 0
	}
	// 构建查询条件
	tx := db.DB("").Where("id in ?", ids)
	// 执行更新删除标记
	if err := tx.Delete(&model.DemoORM{}).Error; err != nil {
		return 0
	}
	return tx.RowsAffected
}

// Clean 清空测试ORM表
func (s DemoORMService) Clean() (int64, error) {
	var rows int64
	err := db.DB("").Model(&model.DemoORM{}).Count(&rows).Error
	if err != nil {
		return 0, err
	}
	// 原生SQL清空表
	db.DB("").Exec("TRUNCATE TABLE demo_orm")
	return rows, nil
}
