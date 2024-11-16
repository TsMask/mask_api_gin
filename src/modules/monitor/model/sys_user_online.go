package model

// SysUserOnline 当前在线会话对象
type SysUserOnline struct {
	TokenID       string `json:"tokenId"`       // 会话编号
	DeptName      string `json:"deptName"`      // 部门名称
	UserName      string `json:"userName"`      // 用户名称
	LoginIp       string `json:"loginIp"`       // 登录IP地址
	LoginLocation string `json:"loginLocation"` // 登录地址
	Browser       string `json:"browser"`       // 浏览器类型
	OS            string `json:"os"`            // 操作系统
	LoginTime     int64  `json:"loginTime"`     // 登录时间
}
