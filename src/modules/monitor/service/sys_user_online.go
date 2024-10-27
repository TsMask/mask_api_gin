package service

import (
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/monitor/model"
)

// NewSysUserOnline 服务层实例化
var NewSysUserOnline = &SysUserOnline{}

// SysUserOnline 在线用户 服务层处理
type SysUserOnline struct{}

// LoginUserToUserOnline 在线用户信息
func (s SysUserOnline) LoginUserToUserOnline(loginUser vo.LoginUser) model.SysUserOnline {
	if loginUser.UserID == "" {
		return model.SysUserOnline{}
	}

	sysUserOnline := model.SysUserOnline{
		TokenID:       loginUser.UUID,
		UserName:      loginUser.User.UserName,
		IPAddr:        loginUser.IPAddr,
		LoginLocation: loginUser.LoginLocation,
		Browser:       loginUser.Browser,
		OS:            loginUser.OS,
		LoginTime:     loginUser.LoginTime,
	}
	if loginUser.User.DeptId != "" {
		sysUserOnline.DeptName = loginUser.User.Dept.DeptName
	}
	return sysUserOnline
}
