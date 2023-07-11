package service

import "mask_api_gin/src/modules/monitor/model"

// ISysUserOnline 在线用户 服务层接口
type ISysUserOnline interface {
	// LoginUserToUserOnline 设置在线用户信息
	LoginUserToUserOnline(loginUser model.SysUserOnline) model.SysUserOnline
}
