package service

import (
	"mask_api_gin/src/framework/token"
	"mask_api_gin/src/modules/monitor/model"
)

// NewSysUserOnline 服务层实例化
var NewSysUserOnline = &SysUserOnline{}

// SysUserOnline 在线用户 服务层处理
type SysUserOnline struct{}

// TokenInfoToUserOnline 在线用户信息
func (s SysUserOnline) TokenInfoToUserOnline(tokenInfo token.TokenInfo) model.SysUserOnline {
	if tokenInfo.UserId <= 0 {
		return model.SysUserOnline{}
	}

	sysUserOnline := model.SysUserOnline{
		TokenID:       tokenInfo.UUID,
		UserName:      tokenInfo.User.UserName,
		LoginIp:       tokenInfo.LoginIp,
		LoginLocation: tokenInfo.LoginLocation,
		Browser:       tokenInfo.Browser,
		OS:            tokenInfo.OS,
		LoginTime:     tokenInfo.LoginTime,
	}
	if tokenInfo.User.DeptId > 0 {
		sysUserOnline.DeptName = tokenInfo.User.Dept.DeptName
	}
	return sysUserOnline
}
