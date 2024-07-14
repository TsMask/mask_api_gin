package service

import "mask_api_gin/src/modules/system/model"

// ISysNoticeService 公告 服务层接口
type ISysNoticeService interface {
	// FindByPage 分页查询列表数据
	FindByPage(query map[string]any) map[string]any

	// Find 查询列表数据
	Find(sysNotice model.SysNotice) []model.SysNotice

	// FindById 通过ID查询信息
	FindById(noticeId string) model.SysNotice

	// Insert 新增信息
	Insert(sysNotice model.SysNotice) string

	// Update 修改信息
	Update(sysNotice model.SysNotice) int64

	// DeleteByIds 批量删除信息
	DeleteByIds(noticeIds []string) (int64, error)
}
