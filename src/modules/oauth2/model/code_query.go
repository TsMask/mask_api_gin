package model

// CodeQuery 重定向授权码参数
type CodeQuery struct {
	RedirectUrl string `form:"redirectUrl"  binding:"required"` // 授权回调地址
	ClientId    string `form:"clientId"  binding:"required"`    // 申请得到的客户端ID
	State       string `form:"state"  binding:"required"`       // 随机字符串，认证服务器会原封不动地返回这个值
}
