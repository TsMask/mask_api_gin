package repository

import "mask_api_gin/src/modules/system/model"

// ISysNotice 通知公告表 数据层接口
type ISysNotice interface {
	// SelectNoticePage 分页查询公告列表
	SelectNoticePage(query map[string]string) map[string]interface{}

	// SelectNoticeList 查询公告列表
	SelectNoticeList(sysNotice model.SysNotice) []model.SysNotice

	// SelectNoticeById 查询公告信息
	SelectNoticeById(noticeId string) model.SysNotice

	// InsertNotice 新增公告
	InsertNotice(sysNotice model.SysNotice) string

	// UpdateNotice 修改公告
	UpdateNotice(sysNotice model.SysNotice) int

	// DeleteNoticeByIds 批量删除公告信息
	DeleteNoticeByIds(noticeIds []string) int
}
