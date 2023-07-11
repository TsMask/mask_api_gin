package model

// SysUserOnline 当前在线会话对象
type SysUserOnline struct {
	// 会话编号
	TokenID string `json:"tokenId"`
	// 部门名称
	DeptName string `json:"deptName"`
	// 用户名称
	UserName string `json:"userName"`
	// 登录IP地址
	IPAddr string `json:"ipaddr"`
	// 登录地址
	LoginLocation string `json:"loginLocation"`
	// 浏览器类型
	Browser string `json:"browser"`
	// 操作系统
	OS string `json:"os"`
	// 登录时间
	LoginTime int64 `json:"loginTime"`
}
