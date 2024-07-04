package service

import (
	"mask_api_gin/src/framework/vo"
	"mask_api_gin/src/modules/monitor/model"
)

// ISysUserOnlineService 在线用户 服务层接口
type ISysUserOnlineService interface {
	// LoginUserToUserOnline 在线用户信息
	LoginUserToUserOnline(loginUser vo.LoginUser) model.SysUserOnline
}
