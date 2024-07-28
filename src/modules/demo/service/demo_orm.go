package service

import (
	db "mask_api_gin/src/framework/data_source"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/demo/model"
	"time"
)

// NewDemoORMService 实例化服务层
// https://gorm.io/zh_CN/docs/query.html
var NewDemoORMService = DemoORMService{}

// DemoORMService 测试ORM信息 服务层处理
type DemoORMService struct{}

// FindByPage 分页查询
func (s *DemoORMService) FindByPage(query map[string]any) (map[string]any, error) {
	// 检查分页条件
	pageNum := int(parse.Number(query["pageNum"]))
	if pageNum < 1 || pageNum > 50 {
		pageNum = 1
	}
	pageSize := int(parse.Number(query["pageSize"]))
	if pageSize < 10 || pageSize > 50 {
		pageSize = 10
	}

	// 条件判断
	where := &model.DemoORM{}
	if v, ok := query["title"]; ok && v != "" {
		where.Title = v.(string)
	}
	if v, ok := query["status"]; ok && v != "" {
		where.Status = v.(string)
	}

	var total int64 = 0
	var rows = make([]model.DemoORM, 0)

	// 执行查询记录总数
	totalResult := db.DB("").Model(&model.DemoORM{}).Where(where).Count(&total)
	if total == 0 || totalResult.Error != nil {
		return map[string]any{
			"total": total,
			"rows":  rows,
		}, totalResult.Error
	}

	// 执行查询记录
	rowsResult := db.DB("").Where(where).Limit(pageSize).Offset(int((pageNum - 1) * pageSize)).Find(&rows)
	if rowsResult.Error != nil {
		return map[string]any{
			"total": total,
			"rows":  rows,
		}, rowsResult.Error
	}

	return map[string]any{
		"total": total,
		"rows":  rows,
	}, nil
}

// Find 查询集合
func (s *DemoORMService) Find(demoORM model.DemoORM) ([]model.DemoORM, error) {

	// 条件判断
	where := &model.DemoORM{}
	if demoORM.Title != "" {
		where.Title = demoORM.Title
	}

	var rows []model.DemoORM
	result := db.DB("").Where(where).Find(&rows)
	if result.Error != nil {
		return nil, result.Error
	}

	return rows, nil
}

// FindById 通过ID查询
func (s *DemoORMService) FindById(id string) (model.DemoORM, error) {
	var result model.DemoORM

	err := db.DB("").First(&result, id).Error
	if err != nil {
		return result, err
	}

	return result, nil
}

// Insert 新增
func (s *DemoORMService) Insert(demoORM model.DemoORM) (model.DemoORM, error) {
	demoORM.CreateBy = "system"
	demoORM.CreateTime = time.Now().UnixMilli()
	result := db.DB("").Create(&demoORM)
	if result.Error != nil {
		return demoORM, result.Error
	}

	return demoORM, nil
}

// Update 更新
func (s *DemoORMService) Update(demoORM model.DemoORM) (model.DemoORM, error) {
	var result model.DemoORM

	err := db.DB("").First(&result, demoORM.ID).Error
	if err != nil {
		return result, err
	}

	// 只改某些属性
	result.Title = demoORM.Title
	result.OrmType = demoORM.OrmType
	result.Status = demoORM.Status
	result.Remark = demoORM.Remark
	result.UpdateBy = "system"
	result.UpdateTime = time.Now().UnixMilli()
	updateResult := db.DB("").Save(&result)
	if updateResult.Error != nil {
		return result, updateResult.Error
	}

	return result, nil
}

// DeleteByIds 批量删除
func (s *DemoORMService) DeleteByIds(ids []string) int64 {
	result := db.DB("").Delete(&model.DemoORM{}, ids)
	if result.Error != nil {
		return 0
	}

	return result.RowsAffected
}

// Clean 清空测试ORM表
func (s *DemoORMService) Clean() (int64, error) {
	var rows int64
	err := db.DB("").Model(&model.DemoORM{}).Count(&rows).Error
	if err != nil {
		return 0, err
	}
	// 原生SQL清空表
	db.DB("").Exec("TRUNCATE TABLE demo_orm")
	return rows, nil
}
