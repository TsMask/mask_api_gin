package service

import "mask_api_gin/src/framework/model"

// 账号身份操作服务 服务层接口
type IAccount interface {
	// Logout 登出清除token
	Logout(token string)

	// CreateToken 创建用户登录令牌
	CreateToken(model.LoginUser) string

	// LoginByUsername 登录生成token
	LoginByUsername(username, password string) (model.LoginUser, error)

	// ValidateCaptcha 校验验证码
	ValidateCaptcha(username, code, uuid string) error

	// ClearLoginRecordCache 清除错误记录次数
	ClearLoginRecordCache(loginName string) bool
}
