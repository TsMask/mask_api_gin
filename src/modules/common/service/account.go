package service

import "mask_api_gin/src/framework/vo"

// 账号身份操作服务 服务层接口
type IAccount interface {
	// ValidateCaptcha 校验验证码
	ValidateCaptcha(username, code, uuid string) error

	// LoginByUsername 登录生成token
	LoginByUsername(username, password string) (vo.LoginUser, error)

	// ClearLoginRecordCache 清除错误记录次数
	ClearLoginRecordCache(loginName string) bool

	// RoleAndMenuPerms 角色和菜单数据权限
	RoleAndMenuPerms(userId string, isAdmin bool) ([]string, []string)

	// RouteMenus 前端路由所需要的菜单
	RouteMenus(userId string, isAdmin bool) []vo.Router
}
