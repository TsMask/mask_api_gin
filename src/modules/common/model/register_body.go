package model

// RegisterBody 用户注册对象
type RegisterBody struct {
	Username        string `json:"username" binding:"required"`        // 用户名
	Password        string `json:"password" binding:"required"`        // 用户密码
	ConfirmPassword string `json:"confirmPassword" binding:"required"` // 用户确认密码
	Code            string `json:"code"`                               // 验证码
	UUID            string `json:"uuid"`                               // 验证码唯一标识
	UserType        string `json:"userType"`                           // 标记用户类型
}
