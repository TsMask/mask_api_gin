package service

import (
	"mask_api_gin/src/framework/token"
	"mask_api_gin/src/modules/monitor/model"
)

// NewSysUserOnline 服务层实例化
var NewSysUserOnline = &SysUserOnline{}

// SysUserOnline 在线用户 服务层处理
type SysUserOnline struct{}

// UserInfoToUserOnline 在线用户信息
func (s SysUserOnline) UserInfoToUserOnline(info token.UserInfo) model.SysUserOnline {
	if info.UserId <= 0 {
		return model.SysUserOnline{}
	}

	sysUserOnline := model.SysUserOnline{
		TokenID:       info.DeviceId,
		UserName:      info.User.UserName,
		LoginIp:       info.LoginIp,
		LoginLocation: info.LoginLocation,
		Browser:       info.Browser,
		OS:            info.OS,
		LoginTime:     info.LoginTime,
	}
	if info.User.DeptId > 0 && info.User.Dept != nil {
		sysUserOnline.DeptName = info.User.Dept.DeptName
	}
	return sysUserOnline
}
