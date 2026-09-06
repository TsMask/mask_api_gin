package model

// RegisterPhoneBody 用户手机号注册对象
type RegisterPhoneBody struct {
	Phone string `json:"phone" binding:"required"` // 用户手机号
	Code  string `json:"code" binding:"required"`  // 短信验证码
	UUID  string `json:"uuid" binding:"required"`  // 短信验证码唯一标识
}
