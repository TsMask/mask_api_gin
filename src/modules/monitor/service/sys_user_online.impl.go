package service

import (
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/monitor/model"
)

// 实例化服务层 SysUserOnlineImpl 结构体
var NewSysUserOnlineImpl = &SysUserOnlineImpl{}

// SysUserOnlineImpl 在线用户 数据层处理
type SysUserOnlineImpl struct{}

// LoginUserToUserOnline 设置在线用户信息
func (r *SysUserOnlineImpl) LoginUserToUserOnline(loginUser vo.LoginUser) model.SysUserOnline {
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
