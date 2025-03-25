package token

// Oauth2Info 第三方客户端令牌信息对象
type Oauth2Info struct {
	DeviceId      string   `json:"deviceId"`      // 用户设备标识
	ClientId      string   `json:"clientId"`      // 客户端ID
	LoginTime     int64    `json:"loginTime"`     // 登录时间时间戳
	ExpireTime    int64    `json:"expireTime"`    // 过期时间时间戳
	LoginIp       string   `json:"loginIp"`       // 登录IP地址 x.x.x.x
	LoginLocation string   `json:"loginLocation"` // 登录地点 xx xx
	Browser       string   `json:"browser"`       // 浏览器类型
	OS            string   `json:"os"`            // 操作系统
	Scope         []string `json:"scope"`         // 权限列表
}
