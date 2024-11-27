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
func (s SysLogLogin) FindByPage(query map[string]string) ([]model.SysLogLogin, int64) {
	return s.SysLogLogin.SelectByPage(query)
}

// Insert 新增信息
func (s SysLogLogin) Insert(userName, status, msg string, ilobArr [4]string) string {
	sysSysLogLogin := model.SysLogLogin{
		LoginIp:       ilobArr[0],
		LoginLocation: ilobArr[1],
		OS:            ilobArr[2],
		Browser:       ilobArr[3],
		UserName:      userName,
		StatusFlag:    status,
		Msg:           msg,
	}
	return s.SysLogLogin.Insert(sysSysLogLogin)
}

// Clean 清空系统登录日志
func (s SysLogLogin) Clean() error {
	return s.SysLogLogin.Clean()
}
