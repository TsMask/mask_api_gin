package service

// 账号注册操作处理 服务层接口
type IRegister interface {
	// ValidateCaptcha 校验验证码
	ValidateCaptcha(code, uuid string) error

	// ByUserName 账号注册
	ByUserName(username, password, userType string) (string, error)
}
