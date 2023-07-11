package service

import (
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// SysNoticeImpl 公告 数据层处理
var SysNoticeImpl = &sysNoticeImpl{
	sysUserRepository: repository.SysNoticeImpl,
}

type sysNoticeImpl struct {
	// 用户服务
	sysUserRepository repository.ISysNotice
}

// SelectNoticePage 分页查询公告列表
func (r *sysNoticeImpl) SelectNoticePage(query map[string]string) map[string]interface{} {
	return r.sysUserRepository.SelectNoticePage(query)
}

// SelectNoticeList 查询公告列表
func (r *sysNoticeImpl) SelectNoticeList(sysNotice model.SysNotice) []model.SysNotice {
	return r.sysUserRepository.SelectNoticeList(sysNotice)
}

// SelectNoticeById 查询公告信息
func (r *sysNoticeImpl) SelectNoticeById(noticeId string) model.SysNotice {
	return r.sysUserRepository.SelectNoticeById(noticeId)
}

// InsertNotice 新增公告
func (r *sysNoticeImpl) InsertNotice(sysNotice model.SysNotice) string {
	return r.sysUserRepository.InsertNotice(sysNotice)
}

// UpdateNotice 修改公告
func (r *sysNoticeImpl) UpdateNotice(sysNotice model.SysNotice) int {
	return r.sysUserRepository.UpdateNotice(sysNotice)
}

// DeleteNoticeByIds 批量删除公告信息
func (r *sysNoticeImpl) DeleteNoticeByIds(noticeIds []string) int {
	return r.sysUserRepository.DeleteNoticeByIds(noticeIds)
}
