package model

// LoginBody 用户登录对象
type LoginBody struct {
	// Username 用户名
	Username string `json:"username"`

	// Password 用户密码
	Password string `json:"password"`

	// Code 验证码
	Code string `json:"code"`

	// UUID 验证码唯一标识
	UUID string `json:"uuid"`
}
