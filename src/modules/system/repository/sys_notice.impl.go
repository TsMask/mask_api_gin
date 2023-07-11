package repository

import "mask_api_gin/src/modules/system/model"

// SysNoticeImpl 通知公告表 数据层处理
var SysNoticeImpl = &sysNoticeImpl{
	selectSql: "",
}

type sysNoticeImpl struct {
	// 查询视图对象SQL
	selectSql string
}

// SelectNoticePage 分页查询公告列表
func (r *sysNoticeImpl) SelectNoticePage(query map[string]string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectNoticeList 查询公告列表
func (r *sysNoticeImpl) SelectNoticeList(sysNotice model.SysNotice) []model.SysNotice {
	return []model.SysNotice{}
}

// SelectNoticeById 查询公告信息
func (r *sysNoticeImpl) SelectNoticeById(noticeId string) model.SysNotice {
	return model.SysNotice{}
}

// InsertNotice 新增公告
func (r *sysNoticeImpl) InsertNotice(sysNotice model.SysNotice) string {
	return ""
}

// UpdateNotice 修改公告
func (r *sysNoticeImpl) UpdateNotice(sysNotice model.SysNotice) int {
	return 0
}

// DeleteNoticeByIds 批量删除公告信息
func (r *sysNoticeImpl) DeleteNoticeByIds(noticeIds []string) int {
	return 0
}
