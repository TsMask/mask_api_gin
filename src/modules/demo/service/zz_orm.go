package service

import (
	"mask_api_gin/src/framework/datasource"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/demo/model"

	"gorm.io/gorm"
)

// Zzorm 测试ORM信息
// https://gorm.io/zh_CN/docs/query.html
var NewZzOrmService = &ZzOrmService{
	db: datasource.DefaultDB,
}

type ZzOrmService struct {
	// 数据库实例
	db func() *gorm.DB
}

// SelectPage 分页查询
func (s *ZzOrmService) SelectPage(query map[string]any) (map[string]any, error) {
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
	where := &model.ZzOrm{}
	if title != "" {
		where.Title = title
	}

	// 执行查询记录总数
	var total int64
	totalResult := s.db().Model(&model.ZzOrm{}).Where(where).Count(&total)
	if total == 0 || totalResult.Error != nil {
		return map[string]any{
			"total": total,
			"rows":  []model.ZzOrm{},
		}, totalResult.Error
	}

	// 执行查询记录
	var rows []model.ZzOrm
	rowsResult := s.db().Where(where).Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&rows)
	if rowsResult.Error != nil {
		return map[string]any{
			"total": total,
			"rows":  []model.ZzOrm{},
		}, rowsResult.Error
	}

	return map[string]any{
		"total": total,
		"rows":  rows,
	}, nil
}

// SelectList 查询集合
func (s *ZzOrmService) SelectList(zzOrm model.ZzOrm) ([]model.ZzOrm, error) {

	// 条件判断
	where := &model.ZzOrm{}
	if zzOrm.Title != "" {
		where.Title = zzOrm.Title
	}

	var rows []model.ZzOrm
	result := s.db().Where(where).Find(&rows)
	if result.Error != nil {
		return nil, result.Error
	}

	return rows, nil
}

// SelectById 通过ID查询
func (s *ZzOrmService) SelectById(id string) (*model.ZzOrm, error) {
	var result model.ZzOrm

	err := s.db().First(&result, id).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Insert 新增
func (s *ZzOrmService) Insert(zzOrm model.ZzOrm) (model.ZzOrm, error) {
	orm := zzOrm
	result := s.db().Create(&orm)
	if result.Error != nil {
		return orm, result.Error
	}

	return orm, nil
}

// Update 更新
func (s *ZzOrmService) Update(zzOrm model.ZzOrm) (*model.ZzOrm, error) {
	var result model.ZzOrm

	err := s.db().First(&result, zzOrm.ID).Error
	if err != nil {
		return nil, err
	}

	// 只改某些属性
	result.Title = zzOrm.Title
	result.OrmType = zzOrm.OrmType
	result.Status = zzOrm.Status
	result.UpdateBy = zzOrm.UpdateBy
	result.UpdateTime = zzOrm.UpdateTime
	result.Remark = zzOrm.Remark
	updateResult := s.db().Save(&result)
	if updateResult.Error != nil {
		return nil, updateResult.Error
	}

	return &result, nil
}

// DeleteByIds 批量删除
func (s *ZzOrmService) DeleteByIds(ids []string) int64 {
	result := s.db().Delete(&model.ZzOrm{}, ids)
	if result.Error != nil {
		return 0
	}

	return result.RowsAffected
}

// Clean 清空测试ORM表
func (s *ZzOrmService) Clean() (int64, error) {
	var rows int64
	err := s.db().Model(&model.ZzOrm{}).Count(&rows).Error
	if err != nil {
		return 0, err
	}
	// 原生SQL清空表
	s.db().Exec("TRUNCATE TABLE zz_orm")
	return rows, nil
}
