package service

import (
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
)

// 实例化服务层 SysLogininforImpl 结构体
var NewSysLogininforImpl = &SysLogininforImpl{
	sysLogininforService: repository.NewSysLogininforImpl,
}

// SysLogininforImpl 系统登录访问 服务层处理
type SysLogininforImpl struct {
	// 系统登录访问信息
	sysLogininforService repository.ISysLogininfor
}

// SelectLogininforPage 分页查询系统登录日志集合
func (s *SysLogininforImpl) SelectLogininforPage(query map[string]string) map[string]interface{} {
	return s.sysLogininforService.SelectLogininforPage(query)
}

// SelectLogininforList 查询系统登录日志集合
func (s *SysLogininforImpl) SelectLogininforList(sysLogininfor model.SysLogininfor) []model.SysLogininfor {
	return s.sysLogininforService.SelectLogininforList(sysLogininfor)
}

// InsertLogininfor 新增系统登录日志
func (s *SysLogininforImpl) InsertLogininfor(sysLogininfor model.SysLogininfor) string {
	return s.sysLogininforService.InsertLogininfor(sysLogininfor)
}

// DeleteLogininforByIds 批量删除系统登录日志
func (s *SysLogininforImpl) DeleteLogininforByIds(infoIds []string) int64 {
	return s.sysLogininforService.DeleteLogininforByIds(infoIds)
}

// CleanLogininfor 清空系统登录日志
func (s *SysLogininforImpl) CleanLogininfor() error {
	return s.sysLogininforService.CleanLogininfor()
}

// NewLogininfor 生成系统登录日志
func (s *SysLogininforImpl) NewLogininfor(userName, status, msg string, ilobArgs ...string) string {
	sysLogininfor := model.SysLogininfor{
		IPAddr:        ilobArgs[0],
		LoginLocation: ilobArgs[1],
		OS:            ilobArgs[2],
		Browser:       ilobArgs[3],
		UserName:      userName,
		Status:        status,
		Msg:           msg,
	}
	return s.InsertLogininfor(sysLogininfor)
}
