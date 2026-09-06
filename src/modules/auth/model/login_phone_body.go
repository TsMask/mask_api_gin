package model

// LoginPhoneBody 手机号验证码登录对象
type LoginPhoneBody struct {
	Phone string `json:"phone" binding:"required"` // 手机号码
	Code  string `json:"code" binding:"required"`  // 短信验证码
	UUID  string `json:"uuid" binding:"required"`  // 短信验证码唯一标识
}
