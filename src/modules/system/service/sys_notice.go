package service

import (
	"fmt"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysNotice 实例化服务层
var NewSysNotice = &SysNotice{
	sysNoticeRepository: repository.NewSysNotice,
}

// SysNotice 公告 服务层处理
type SysNotice struct {
	sysNoticeRepository *repository.SysNotice // 公告服务
}

// FindByPage 分页查询列表数据
func (s SysNotice) FindByPage(query map[string]any) ([]model.SysNotice, int64) {
	return s.sysNoticeRepository.SelectByPage(query)
}

// Find 查询列表数据
func (s SysNotice) Find(sysNotice model.SysNotice) []model.SysNotice {
	return s.sysNoticeRepository.Select(sysNotice)
}

// FindById 通过ID查询信息
func (s SysNotice) FindById(noticeId string) model.SysNotice {
	if noticeId == "" {
		return model.SysNotice{}
	}
	configs := s.sysNoticeRepository.SelectByIds([]string{noticeId})
	if len(configs) > 0 {
		return configs[0]
	}
	return model.SysNotice{}
}

// Insert 新增信息
func (s SysNotice) Insert(sysNotice model.SysNotice) string {
	return s.sysNoticeRepository.Insert(sysNotice)
}

// Update 修改信息
func (s SysNotice) Update(sysNotice model.SysNotice) int64 {
	return s.sysNoticeRepository.Update(sysNotice)
}

// DeleteByIds 批量删除信息
func (s SysNotice) DeleteByIds(noticeIds []string) (int64, error) {
	// 检查是否存在
	notices := s.sysNoticeRepository.SelectByIds(noticeIds)
	if len(notices) <= 0 {
		return 0, fmt.Errorf("没有权限访问公告信息数据！")
	}
	for _, notice := range notices {
		// 检查是否为已删除
		if notice.DelFlag == "1" {
			return 0, fmt.Errorf(notice.NoticeId + " 公告信息已经删除！")
		}
	}
	if len(notices) == len(noticeIds) {
		rows := s.sysNoticeRepository.DeleteByIds(noticeIds)
		return rows, nil
	}
	return 0, fmt.Errorf("删除公告信息失败！")
}
