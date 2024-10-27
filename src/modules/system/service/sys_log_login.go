package service

import (
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/repository"
)

// NewSysLogLogin 实例化服务层
var NewSysLogLogin = &SysLogLogin{
	SysLogLogin: repository.NewSysLogLogin,
}

// SysLogLogin 系统登录日志 服务层处理
type SysLogLogin struct {
	SysLogLogin *repository.SysLogLogin // 系统登录日志信息
}

// FindByPage 分页查询列表数据
func (s SysLogLogin) FindByPage(query map[string]any) map[string]any {
	return s.SysLogLogin.SelectByPage(query)
}

// Find 查询数据
func (s SysLogLogin) Find(sysSysLogLogin model.SysLogLogin) []model.SysLogLogin {
	return s.SysLogLogin.Select(sysSysLogLogin)
}

// Insert 新增信息
func (s SysLogLogin) Insert(userName, status, msg string, ilobArr [4]string) string {
	sysSysLogLogin := model.SysLogLogin{
		IPAddr:        ilobArr[0],
		LoginLocation: ilobArr[1],
		OS:            ilobArr[2],
		Browser:       ilobArr[3],
		UserName:      userName,
		Status:        status,
		Msg:           msg,
	}
	return s.SysLogLogin.Insert(sysSysLogLogin)
}

// DeleteByIds 批量删除信息
func (s SysLogLogin) DeleteByIds(loginIds []string) int64 {
	return s.SysLogLogin.DeleteByIds(loginIds)
}

// Clean 清空系统登录日志
func (s SysLogLogin) Clean() error {
	return s.SysLogLogin.Clean()
}
