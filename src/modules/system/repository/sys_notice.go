package repository

import "mask_api_gin/src/modules/system/model"

// ISysNoticeRepository 通知公告表 数据层接口
type ISysNoticeRepository interface {
	// SelectByPage 分页查询集合
	SelectByPage(query map[string]any) map[string]any

	// Select 查询集合
	Select(sysNotice model.SysNotice) []model.SysNotice

	// SelectByIds 通过ID查询信息
	SelectByIds(noticeIds []string) []model.SysNotice

	// Insert 新增信息
	Insert(sysNotice model.SysNotice) string

	// Update 修改信息
	Update(sysNotice model.SysNotice) int64

	// DeleteByIds 批量删除信息
	DeleteByIds(noticeIds []string) int64
}
