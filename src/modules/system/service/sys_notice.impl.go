package service

import (
	"errors"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// 实例化服务层 SysNoticeImpl 结构体
var NewSysNoticeImpl = &SysNoticeImpl{
	sysNoticeRepository: repository.NewSysNoticeImpl,
}

// SysNoticeImpl 公告 服务层处理
type SysNoticeImpl struct {
	// 公告服务
	sysNoticeRepository repository.ISysNotice
}

// SelectNoticePage 分页查询公告列表
func (r *SysNoticeImpl) SelectNoticePage(query map[string]any) map[string]any {
	return r.sysNoticeRepository.SelectNoticePage(query)
}

// SelectNoticeList 查询公告列表
func (r *SysNoticeImpl) SelectNoticeList(sysNotice model.SysNotice) []model.SysNotice {
	return r.sysNoticeRepository.SelectNoticeList(sysNotice)
}

// SelectNoticeById 查询公告信息
func (r *SysNoticeImpl) SelectNoticeById(noticeId string) model.SysNotice {
	if noticeId == "" {
		return model.SysNotice{}
	}
	configs := r.sysNoticeRepository.SelectNoticeByIds([]string{noticeId})
	if len(configs) > 0 {
		return configs[0]
	}
	return model.SysNotice{}
}

// InsertNotice 新增公告
func (r *SysNoticeImpl) InsertNotice(sysNotice model.SysNotice) string {
	return r.sysNoticeRepository.InsertNotice(sysNotice)
}

// UpdateNotice 修改公告
func (r *SysNoticeImpl) UpdateNotice(sysNotice model.SysNotice) int64 {
	return r.sysNoticeRepository.UpdateNotice(sysNotice)
}

// DeleteNoticeByIds 批量删除公告信息
func (r *SysNoticeImpl) DeleteNoticeByIds(noticeIds []string) (int64, error) {
	// 检查是否存在
	notices := r.sysNoticeRepository.SelectNoticeByIds(noticeIds)
	if len(notices) <= 0 {
		return 0, errors.New("没有权限访问公告信息数据！")
	}
	for _, notice := range notices {
		// 检查是否为已删除
		if notice.DelFlag == "1" {
			return 0, errors.New(notice.NoticeID + " 公告信息已经删除！")
		}
	}
	if len(notices) == len(noticeIds) {
		rows := r.sysNoticeRepository.DeleteNoticeByIds(noticeIds)
		return rows, nil
	}
	return 0, errors.New("删除公告信息失败！")
}
