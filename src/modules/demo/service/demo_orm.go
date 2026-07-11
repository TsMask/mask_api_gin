package service

import (
	"time"

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
	return s.demoORMRepository.SelectByPage(query)
}

// Find 查询集合
func (s DemoORMService) Find(demoORM model.DemoORM) []model.DemoORM {
	return s.demoORMRepository.Select(demoORM)
}

// FindById 通过ID查询
func (s DemoORMService) FindById(id int64) model.DemoORM {
	return s.demoORMRepository.SelectById(id)
}

// Insert 新增
func (s DemoORMService) Insert(demoORM model.DemoORM) int64 {
	demoORM.CreateBy = "system"
	demoORM.CreateTime = time.Now().UnixMilli()
	demoORM.UpdateBy = demoORM.CreateBy
	demoORM.UpdateTime = demoORM.CreateTime
	return s.demoORMRepository.Insert(demoORM)
}

// Update 更新
func (s DemoORMService) Update(demoORM model.DemoORM) int64 {
	if demoORM.Id <= 0 {
		return 0
	}
	item := s.demoORMRepository.SelectById(demoORM.Id)
	if item.Id <= 0 {
		return 0
	}
	item.Title = demoORM.Title
	item.OrmType = demoORM.OrmType
	item.OrmStatus = demoORM.OrmStatus
	item.Remark = demoORM.Remark
	item.UpdateBy = "system"
	item.UpdateTime = time.Now().UnixMilli()
	return s.demoORMRepository.Update(item)
}

// DeleteByIds 批量删除
func (s DemoORMService) DeleteByIds(ids []int64) int64 {
	return s.demoORMRepository.DeleteByIds(ids)
}

// Clean 清空测试ORM表
func (s DemoORMService) Clean() int64 {
	return s.demoORMRepository.Clean()
}
