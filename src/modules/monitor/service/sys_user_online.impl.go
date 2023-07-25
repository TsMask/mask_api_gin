package service

import (
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/monitor/model"
)

// SysUserOnlineImpl 在线用户 数据层处理
var SysUserOnlineImpl = &sysUserOnlineImpl{}

type sysUserOnlineImpl struct{}

// LoginUserToUserOnline 设置在线用户信息
func (r *sysUserOnlineImpl) LoginUserToUserOnline(loginUser vo.LoginUser) model.SysUserOnline {
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
	if loginUser.User.DeptID != "" {
		sysUserOnline.DeptName = loginUser.User.Dept.DeptName
	}
	return sysUserOnline
}
