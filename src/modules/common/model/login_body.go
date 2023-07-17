package model

// LoginBody 用户登录对象
type LoginBody struct {
	// Username 用户名
	Username string `json:"username" binding:"required"`

	// Password 用户密码
	Password string `json:"password" binding:"required"`

	// Code 验证码
	Code string `json:"code" binding:"required"`

	// UUID 验证码唯一标识
	UUID string `json:"uuid" binding:"required"`
}
