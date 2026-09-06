package model

// PasswordResetBody 重置密码对象
type PasswordResetBody struct {
	Username        string `json:"username" binding:"required"`         // 用户账号
	Code            string `json:"code" binding:"required"`             // 找回密码验证码
	UUID            string `json:"uuid" binding:"required"`             // 找回密码验证码唯一标识
	Password        string `json:"password" binding:"required"`         // 新密码
	ConfirmPassword string `json:"confirm_password" binding:"required"` // 确认新密码
}

// PasswordResetTarget 密码找回目标
type PasswordResetTarget struct {
	UserID   int64  // 用户ID
	UserName string // 用户账号
	Channel  string // 联系方式类型
	Contact  string // 联系方式
}
