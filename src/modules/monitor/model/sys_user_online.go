package model

// SysUserOnline 当前在线会话对象
type SysUserOnline struct {
	TokenID       string `json:"token_id"`       // 会话编号
	DeptName      string `json:"dept_name"`      // 部门名称
	UserName      string `json:"user_name"`      // 用户名称
	LoginIp       string `json:"login_ip"`       // 登录IP地址
	LoginLocation string `json:"login_location"` // 登录地址
	Browser       string `json:"browser"`        // 浏览器类型
	OS            string `json:"os"`             // 操作系统
	LoginTime     int64  `json:"login_time"`     // 登录时间
}
