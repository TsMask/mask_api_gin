package model

// SysLogininfor 系统访问记录表 sys_logininfor
type SysLogininfor struct {
	// 访问ID
	InfoID string `json:"infoId"`
	// 用户账号
	UserName string `json:"userName"`
	// 登录IP地址
	IPAddr string `json:"ipaddr"`
	// 登录地点
	LoginLocation string `json:"loginLocation"`
	// 浏览器类型
	Browser string `json:"browser"`
	// 操作系统
	OS string `json:"os"`
	// 登录状态（0失败 1成功）
	Status string `json:"status"`
	// 提示消息
	Msg string `json:"msg"`
	// 访问时间
	LoginTime int64 `json:"loginTime"`
}
