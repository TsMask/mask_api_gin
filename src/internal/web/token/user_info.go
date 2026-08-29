package token

import "encoding/json"

// UserInfo 系统用户令牌信息对象
type UserInfo struct {
	SessionId     string          `json:"session_id"`     // 会话唯一标识
	DeviceId      string          `json:"device_id"`      // 用户设备标识
	UserId        int64           `json:"user_id"`        // 用户ID
	DeptId        int64           `json:"dept_id"`        // 部门ID
	LoginIp       string          `json:"login_ip"`       // 登录IP地址 x.x.x.x
	LoginLocation string          `json:"login_location"` // 登录地点 xx xx
	LoginBrowser  string          `json:"login_browser"`  // 浏览器类型
	LoginOS       string          `json:"login_os"`       // 操作系统
	LoginTime     int64           `json:"login_time"`     // 登录时间时间戳
	ExpireTime    int64           `json:"expire_time"`    // 过期时间时间戳
	Permissions   []string        `json:"permissions"`    // 权限列表
	UserName      string          `json:"user_name"`      // login account name (top-level, for web reads)
	RoleKeys      []string        `json:"role_keys"`      // role keys (top-level, for web reads)
	User          json.RawMessage `json:"user"`           // 用户信息
}
