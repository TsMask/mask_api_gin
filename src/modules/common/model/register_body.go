package model

// RegisterBody 用户注册对象
type RegisterBody struct {
	// Username 用户名
	Username string `json:"username" binding:"required"`

	// Password 用户密码
	Password string `json:"password" binding:"required"`

	// ConfirmPassword 用户确认密码
	ConfirmPassword string `json:"confirmPassword" binding:"required"`

	// Code 验证码
	Code string `json:"code" binding:"required"`

	// UUID 验证码唯一标识
	UUID string `json:"uuid" binding:"required"`

	// UserType 标记用户类型
	UserType string `json:"userType"`
}
