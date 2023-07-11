package service

// 账号身份操作服务 服务层接口
type IAccount interface {
	// Logout 登出清除token
	Logout(token string)
}
