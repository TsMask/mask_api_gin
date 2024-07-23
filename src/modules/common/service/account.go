package service

import "mask_api_gin/src/framework/vo"

// IAccountService 账号身份操作 服务层接口
type IAccountService interface {
	// ValidateCaptcha 校验验证码
	ValidateCaptcha(code, uuid string) error

	// ByUsername 登录创建用户信息
	ByUsername(username, password string) (vo.LoginUser, error)

	// UpdateLoginDateAndIP 更新登录时间和IP
	UpdateLoginDateAndIP(loginUser *vo.LoginUser) bool

	// CleanLoginRecordCache 清除错误记录次数
	CleanLoginRecordCache(username string) bool

	// RoleAndMenuPerms 角色和菜单数据权限
	RoleAndMenuPerms(userId string, isSysAdmin bool) ([]string, []string)

	// RouteMenus 前端路由所需要的菜单
	RouteMenus(userId string, isSysAdmin bool) []vo.Router
}
