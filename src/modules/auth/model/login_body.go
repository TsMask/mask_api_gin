package model

// LoginBody 用户登录对象
type LoginBody struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 用户密码
	Code     string `json:"code"`                        // 验证码
	UUID     string `json:"uuid"`                        // 验证码唯一标识
}
