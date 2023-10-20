package service

import (
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// 实例化服务层 SysLogLoginImpl 结构体
var NewSysLogLoginImpl = &SysLogLoginImpl{
	sysLogLoginService: repository.NewSysLogLoginImpl,
}

// SysLogLoginImpl 系统登录访问 服务层处理
type SysLogLoginImpl struct {
	// 系统登录访问信息
	sysLogLoginService repository.ISysLogLogin
}

// SelectSysLogLoginPage 分页查询系统登录日志集合
func (s *SysLogLoginImpl) SelectSysLogLoginPage(query map[string]any) map[string]any {
	return s.sysLogLoginService.SelectSysLogLoginPage(query)
}

// SelectSysLogLoginList 查询系统登录日志集合
func (s *SysLogLoginImpl) SelectSysLogLoginList(sysSysLogLogin model.SysLogLogin) []model.SysLogLogin {
	return s.sysLogLoginService.SelectSysLogLoginList(sysSysLogLogin)
}

// InsertSysLogLogin 新增系统登录日志
func (s *SysLogLoginImpl) InsertSysLogLogin(sysSysLogLogin model.SysLogLogin) string {
	return s.sysLogLoginService.InsertSysLogLogin(sysSysLogLogin)
}

// DeleteSysLogLoginByIds 批量删除系统登录日志
func (s *SysLogLoginImpl) DeleteSysLogLoginByIds(loginIds []string) int64 {
	return s.sysLogLoginService.DeleteSysLogLoginByIds(loginIds)
}

// CleanSysLogLogin 清空系统登录日志
func (s *SysLogLoginImpl) CleanSysLogLogin() error {
	return s.sysLogLoginService.CleanSysLogLogin()
}

// CreateSysLogLogin 创建系统登录日志
func (s *SysLogLoginImpl) CreateSysLogLogin(userName, status, msg string, ilobArgs ...string) string {
	sysSysLogLogin := model.SysLogLogin{
		IPAddr:        ilobArgs[0],
		LoginLocation: ilobArgs[1],
		OS:            ilobArgs[2],
		Browser:       ilobArgs[3],
		UserName:      userName,
		Status:        status,
		Msg:           msg,
	}
	return s.InsertSysLogLogin(sysSysLogLogin)
}
