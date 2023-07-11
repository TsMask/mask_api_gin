package service

import (
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
)

// SysLogininforImpl 系统登录访问 业务层处理
var SysLogininforImpl = &sysLogininforImpl{
	sysLogininforService: repository.SysLogininforImpl,
}

type sysLogininforImpl struct {
	// 系统登录访问信息
	sysLogininforService repository.ISysLogininfor
}

// SelectLogininforPage 分页查询系统登录日志集合
func (s *sysLogininforImpl) SelectLogininforPage(query map[string]string) map[string]interface{} {
	return map[string]interface{}{}
}

// SelectLogininforList 查询系统登录日志集合
func (s *sysLogininforImpl) SelectLogininforList(sysLogininfor model.SysLogininfor) []model.SysLogininfor {
	return []model.SysLogininfor{}
}

// InsertLogininfor 新增系统登录日志
func (s *sysLogininforImpl) InsertLogininfor(sysLogininfor model.SysLogininfor) string {
	return ""
}

// DeleteLogininforByIds 批量删除系统登录日志
func (s *sysLogininforImpl) DeleteLogininforByIds(infoIds []string) int64 {
	return 0
}

// CleanLogininfor 清空系统登录日志
func (s *sysLogininforImpl) CleanLogininfor() error {
	return nil
}
