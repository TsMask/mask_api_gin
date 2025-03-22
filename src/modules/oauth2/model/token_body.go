package model

// TokenBody 获取访问令牌参数
type TokenBody struct {
	ClientId     string `json:"clientId" binding:"required"`                                         // 申请应用时获得的client_id
	ClientSecret string `json:"clientSecret" binding:"required"`                                     // 申请应用时分配的secret
	GrantType    string `json:"grantType" binding:"required,oneof=authorization_code refresh_token"` // 请求的类型，此处的值固定为 authorization_code/refresh_token
	Code         string `json:"code"`                                                                // 授权拿到的code值
	RefreshToken string `json:"refreshToken"`                                                        // 刷新令牌
}
