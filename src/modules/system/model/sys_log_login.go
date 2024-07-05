package model

// SysLogLogin 系统登录日志表 sys_log_login
type SysLogLogin struct {
	LoginID       string `json:"loginId"`       // 登录ID
	UserName      string `json:"userName"`      // 用户账号
	IPAddr        string `json:"ipaddr"`        // 登录IP地址
	LoginLocation string `json:"loginLocation"` // 登录地点
	Browser       string `json:"browser"`       // 浏览器类型
	OS            string `json:"os"`            // 操作系统
	Status        string `json:"status"`        // 登录状态（0失败 1成功）
	Msg           string `json:"msg"`           // 提示消息
	LoginTime     int64  `json:"loginTime"`     // 访问时间
}
