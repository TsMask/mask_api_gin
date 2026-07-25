package service

import (
	"encoding/json"

	"mask_api_gin/src/internal/web/token"
	"mask_api_gin/src/modules/monitor/model"
	systemModel "mask_api_gin/src/modules/system/model"
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

	var userInfo systemModel.SysUser
	if len(info.User) > 0 {
		_ = json.Unmarshal(info.User, &userInfo)
	}
	sysUserOnline := model.SysUserOnline{
		TokenID:       info.DeviceId,
		UserName:      userInfo.UserName,
		LoginIp:       info.LoginIp,
		LoginLocation: info.LoginLocation,
		Browser:       info.LoginBrowser,
		OS:            info.LoginOS,
		LoginTime:     info.LoginTime,
	}
	if userInfo.DeptId > 0 && userInfo.Dept != nil {
		sysUserOnline.DeptName = userInfo.Dept.DeptName
	}
	return sysUserOnline
}
