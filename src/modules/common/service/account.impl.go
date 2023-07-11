package service

import (
	monitorService "mask_api_gin/src/modules/monitor/service"
	systemService "mask_api_gin/src/modules/system/service"
)

// 账号身份操作服务 业务层处理
var AccountImpl = &accountImpl{
	sysUserService:       systemService.SysUserImpl,
	sysLogininforService: monitorService.SysLogininforImpl,
}

type accountImpl struct {
	// 用户信息服务
	sysUserService systemService.ISysUser
	// 系统登录访问信息服务
	sysLogininforService monitorService.ISysLogininfor
}

// Logout 登出清除Token
func (s *accountImpl) Logout(token string) {

}
