package service

import (
	db "mask_api_gin/src/framework/data_source"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/demo/model"
)

// NewDemoORMService 测试ORM信息
// https://gorm.io/zh_CN/docs/query.html
var NewDemoORMService = DemoORMService{}

type DemoORMService struct{}

// SelectPage 分页查询
func (s *DemoORMService) SelectPage(query map[string]any) (map[string]any, error) {
	var (
		pageSize int64
		pageNum  int64
		title    string
	)
	if v, ok := query["pageSize"]; ok && v != "" {
		pageSize = parse.Number(v)
	}
	if v, ok := query["pageNum"]; ok && v != "" {
		pageNum = parse.Number(v)
	}
	if v, ok := query["title"]; ok && v != "" {
		title = v.(string)
	}

	// 检查分页条件
	if pageSize < 0 || pageSize > 50 {
		pageSize = 0
	}
	if pageNum < 1 || pageNum > 50 {
		pageNum = 10
	}

	// 条件判断
	where := &model.DemoORM{}
	if title != "" {
		where.Title = title
	}

	// 执行查询记录总数
	var total int64
	totalResult := db.DB("").Model(&model.DemoORM{}).Where(where).Count(&total)
	if total == 0 || totalResult.Error != nil {
		return map[string]any{
			"total": total,
			"rows":  []model.DemoORM{},
		}, totalResult.Error
	}

	// 执行查询记录
	var rows []model.DemoORM
	rowsResult := db.DB("").Where(where).Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&rows)
	if rowsResult.Error != nil {
		return map[string]any{
			"total": total,
			"rows":  []model.DemoORM{},
		}, rowsResult.Error
	}

	return map[string]any{
		"total": total,
		"rows":  rows,
	}, nil
}

// SelectList 查询集合
func (s *DemoORMService) SelectList(demoORM model.DemoORM) ([]model.DemoORM, error) {

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

// SelectById 通过ID查询
func (s *DemoORMService) SelectById(id string) (model.DemoORM, error) {
	var result model.DemoORM

	err := db.DB("").First(&result, id).Error
	if err != nil {
		return result, err
	}

	return result, nil
}

// Insert 新增
func (s *DemoORMService) Insert(demoORM model.DemoORM) (model.DemoORM, error) {
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
	result.UpdateBy = demoORM.UpdateBy
	result.UpdateTime = demoORM.UpdateTime
	result.Remark = demoORM.Remark
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
