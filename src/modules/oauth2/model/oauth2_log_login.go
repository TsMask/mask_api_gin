package model

// Oauth2LogLogin 用户授权第三方应用登录日志表
type Oauth2LogLogin struct {
	ID            int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"` // 登录ID
	ClientId      string `json:"clientId" gorm:"column:client_id"`             // 应用的唯一标识
	LoginIp       string `json:"loginIp" gorm:"column:login_ip"`               // 登录IP地址
	LoginLocation string `json:"loginLocation" gorm:"column:login_location"`   // 登录地点
	Browser       string `json:"browser" gorm:"column:browser"`                // 浏览器类型
	OS            string `json:"os" gorm:"column:os"`                          // 操作系统
	StatusFlag    string `json:"statusFlag" gorm:"column:status_flag"`         // 登录状态（0失败 1成功）
	Msg           string `json:"msg" gorm:"column:msg"`                        // 提示消息
	LoginTime     int64  `json:"loginTime" gorm:"column:login_time"`           // 登录时间
}

// TableName 表名称
func (*Oauth2LogLogin) TableName() string {
	return "oauth2_log_login"
}
