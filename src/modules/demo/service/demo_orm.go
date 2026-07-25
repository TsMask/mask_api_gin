package service

import (
	"time"

	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/modules/demo/model"
	"mask_api_gin/src/modules/demo/repository"
)

// NewDemoORMService 实例化服务层
var NewDemoORMService = &DemoORMService{
	demoORMRepository: repository.NewDemoORM,
}

// DemoORMService 测试ORM信息 服务层处理
type DemoORMService struct {
	demoORMRepository *repository.DemoORM
}

// FindByPage 分页查询
func (s DemoORMService) FindByPage(query map[string]string) ([]model.DemoORM, int64) {
	rows, total, err := s.demoORMRepository.SelectByPage(query)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return rows, total
}

// Find 查询集合
func (s DemoORMService) Find(demoORM model.DemoORM) []model.DemoORM {
	rows, err := s.demoORMRepository.Select(demoORM)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return rows
}

// FindById 通过ID查询
func (s DemoORMService) FindById(id int64) model.DemoORM {
	item, err := s.demoORMRepository.SelectById(id)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return item
}

// Insert 新增
func (s DemoORMService) Insert(demoORM model.DemoORM) int64 {
	ms := time.Now().UnixMilli()
	demoORM.CreateBy = "system"
	demoORM.UpdateBy = demoORM.CreateBy
	demoORM.UpdateTime = ms
	demoORM.CreateTime = ms
	rows, err := s.demoORMRepository.Insert(demoORM)
	if err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
	}
	return rows
}

// Update 更新
func (s DemoORMService) Update(demoORM model.DemoORM) int64 {
	if demoORM.Id <= 0 {
		return 0
	}
	item, err := s.demoORMRepository.SelectById(demoORM.Id)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	if item.Id <= 0 {
		return 0
	}
	item.Title = demoORM.Title
	item.OrmType = demoORM.OrmType
	item.OrmStatus = demoORM.OrmStatus
	item.Remark = demoORM.Remark
	item.UpdateBy = "system"
	item.UpdateTime = time.Now().UnixMilli()
	rows, err := s.demoORMRepository.Update(item)
	if err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return rows
}

// DeleteByIds 批量删除
func (s DemoORMService) DeleteByIds(ids []int64) int64 {
	rows, err := s.demoORMRepository.DeleteByIds(ids)
	if err != nil {
		logger.Errorf("delete err => %v", err.Error())
		return 0
	}
	return rows
}

// Clean 清空测试ORM表
func (s DemoORMService) Clean() int64 {
	rows, err := s.demoORMRepository.Clean()
	if err != nil {
		logger.Errorf("clean err => %v", err.Error())
		return 0
	}
	return rows
}
