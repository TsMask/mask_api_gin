package model

// RegisterBody 用户注册对象
type RegisterBody struct {
	// Username 用户名
	Username string `json:"username"`

	// Password 用户密码
	Password string `json:"password"`

	// ConfirmPassword 用户确认密码
	ConfirmPassword string `json:"confirmPassword"`

	// Code 验证码
	Code string `json:"code"`

	// UUID 验证码唯一标识
	UUID string `json:"uuid"`

	// UserType 标记用户类型
	UserType string `json:"userType"`
}
