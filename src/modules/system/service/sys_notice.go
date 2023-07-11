package service

import "mask_api_gin/src/modules/system/model"

// ISysNotice 公告 服务层接口
type ISysNotice interface {
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
