package service

// IRegisterService 账号注册操作 服务层接口
type IRegisterService interface {
	// ValidateCaptcha 校验验证码
	ValidateCaptcha(code, uuid string) error

	// ByUserName 账号注册
	ByUserName(username, password, userType string) (string, error)
}
