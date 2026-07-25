package token

import "encoding/json"

// UserInfo 系统用户令牌信息对象
type UserInfo struct {
	DeviceId      string          `json:"deviceId"`      // 用户设备标识
	UserId        int64           `json:"userId"`        // 用户ID
	DeptId        int64           `json:"deptId"`        // 部门ID
	LoginIp       string          `json:"loginIp"`       // 登录IP地址 x.x.x.x
	LoginLocation string          `json:"loginLocation"` // 登录地点 xx xx
	LoginBrowser  string          `json:"loginBrowser"`  // 浏览器类型
	LoginOS       string          `json:"loginOS"`       // 操作系统
	LoginTime     int64           `json:"loginTime"`     // 登录时间时间戳
	ExpireTime    int64           `json:"expireTime"`    // 过期时间时间戳
	Permissions   []string        `json:"permissions"`   // 权限列表
	UserName      string          `json:"userName"`      // login account name (top-level, for web reads)
	RoleKeys      []string        `json:"roleKeys"`      // role keys (top-level, for web reads)
	User          json.RawMessage `json:"user"`          // 用户信息
}
