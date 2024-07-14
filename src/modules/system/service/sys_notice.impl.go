package service

import (
	"fmt"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysNotice 实例化服务层
var NewSysNotice = &SysNoticeService{
	sysNoticeRepository: repository.NewSysNotice,
}

// SysNoticeService 公告 服务层处理
type SysNoticeService struct {
	sysNoticeRepository repository.ISysNoticeRepository // 公告服务
}

// FindByPage 分页查询列表数据
func (r *SysNoticeService) FindByPage(query map[string]any) map[string]any {
	return r.sysNoticeRepository.SelectByPage(query)
}

// Find 查询列表数据
func (r *SysNoticeService) Find(sysNotice model.SysNotice) []model.SysNotice {
	return r.sysNoticeRepository.Select(sysNotice)
}

// FindById 通过ID查询信息
func (r *SysNoticeService) FindById(noticeId string) model.SysNotice {
	if noticeId == "" {
		return model.SysNotice{}
	}
	configs := r.sysNoticeRepository.SelectByIds([]string{noticeId})
	if len(configs) > 0 {
		return configs[0]
	}
	return model.SysNotice{}
}

// Insert 新增信息
func (r *SysNoticeService) Insert(sysNotice model.SysNotice) string {
	return r.sysNoticeRepository.Insert(sysNotice)
}

// Update 修改信息
func (r *SysNoticeService) Update(sysNotice model.SysNotice) int64 {
	return r.sysNoticeRepository.Update(sysNotice)
}

// DeleteByIds 批量删除信息
func (r *SysNoticeService) DeleteByIds(noticeIds []string) (int64, error) {
	// 检查是否存在
	notices := r.sysNoticeRepository.SelectByIds(noticeIds)
	if len(notices) <= 0 {
		return 0, fmt.Errorf("没有权限访问公告信息数据！")
	}
	for _, notice := range notices {
		// 检查是否为已删除
		if notice.DelFlag == "1" {
			return 0, fmt.Errorf(notice.NoticeID + " 公告信息已经删除！")
		}
	}
	if len(notices) == len(noticeIds) {
		rows := r.sysNoticeRepository.DeleteByIds(noticeIds)
		return rows, nil
	}
	return 0, fmt.Errorf("删除公告信息失败！")
}
