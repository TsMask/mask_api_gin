package model

// SysLogLogin 系统登录日志表
type SysLogLogin struct {
	LoginId       string `json:"login_id" gorm:"column:login_id"`             // 登录ID
	UserName      string `json:"user_name" gorm:"column:user_name"`           // 用户账号
	IPAddr        string `json:"ipaddr" gorm:"column:ipaddr"`                 // 登录IP地址
	LoginLocation string `json:"login_location" gorm:"column:login_location"` // 登录地点
	Browser       string `json:"browser" gorm:"column:browser"`               // 浏览器类型
	OS            string `json:"os" gorm:"column:os"`                         // 操作系统
	Status        string `json:"status" gorm:"column:status"`                 // 登录状态（0失败 1成功）
	Msg           string `json:"msg" gorm:"column:msg"`                       // 提示消息
	LoginTime     int64  `json:"login_time" gorm:"column:login_time"`         // 登录时间
}

// TableName 表名称
func (*SysLogLogin) TableName() string {
	return "sys_log_login"
}
