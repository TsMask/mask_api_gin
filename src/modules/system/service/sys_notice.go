package service

import (
	"fmt"
	"time"

	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/logger"
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
func (s SysNotice) FindByPage(query map[string]string) ([]model.SysNotice, int64) {
	rows, total, err := s.sysNoticeRepository.SelectByPage(query)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return rows, total
}

// Find 查询列表数据
func (s SysNotice) Find(sysNotice model.SysNotice) []model.SysNotice {
	rows, err := s.sysNoticeRepository.Select(sysNotice)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	return rows
}

// FindById 通过ID查询信息
func (s SysNotice) FindById(noticeId int64) model.SysNotice {
	if noticeId <= 0 {
		return model.SysNotice{}
	}
	configs, err := s.sysNoticeRepository.SelectByIds([]int64{noticeId})
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	if len(configs) > 0 {
		return configs[0]
	}
	return model.SysNotice{}
}

// Insert 新增信息
func (s SysNotice) Insert(sysNotice model.SysNotice) int64 {
	if sysNotice.CreateBy != "" {
		ms := time.Now().UnixMilli()
		sysNotice.UpdateBy = sysNotice.CreateBy
		sysNotice.UpdateTime = ms
		sysNotice.CreateTime = ms
	}
	rows, err := s.sysNoticeRepository.Insert(sysNotice)
	if err != nil {
		logger.Errorf("insert err => %v", err.Error())
		return 0
	}
	return rows
}

// Update 修改信息
func (s SysNotice) Update(sysNotice model.SysNotice) int64 {
	if sysNotice.UpdateBy != "" {
		sysNotice.UpdateTime = time.Now().UnixMilli()
	}
	rows, err := s.sysNoticeRepository.Update(sysNotice)
	if err != nil {
		logger.Errorf("update err => %v", err.Error())
		return 0
	}
	return rows
}

// DeleteByIds 批量删除信息
func (s SysNotice) DeleteByIds(noticeIds []int64, updateBy string) (int64, error) {
	// 检查是否存在
	notices, err := s.sysNoticeRepository.SelectByIds(noticeIds)
	if err != nil {
		logger.Errorf("query err => %v", err.Error())
	}
	if len(notices) <= 0 {
		return 0, fmt.Errorf("没有权限访问公告信息数据！")
	}
	if len(notices) == len(noticeIds) {
		deletableIds := make([]int64, 0)
		for _, notice := range notices {
			if notice.NoticeStatus != constants.NOTICE_STATUS_DELETED {
				deletableIds = append(deletableIds, notice.ID)
			}
		}
		if len(deletableIds) == 0 {
			return 0, fmt.Errorf("公告已是已删除状态，无需重复删除")
		}

		rows, err := s.sysNoticeRepository.DeleteByIds(deletableIds, updateBy)
		if err != nil {
			return 0, err
		}
		if rows == int64(len(deletableIds)) {
			return rows, nil
		}
	}
	return 0, fmt.Errorf("删除公告信息失败！")
}
