package token

import systemModel "mask_api_gin/src/modules/system/model"

// UserInfo 系统用户令牌信息对象
type UserInfo struct {
	DeviceId      string              `json:"deviceId"`      // 用户设备标识
	UserId        int64               `json:"userId"`        // 用户ID
	DeptId        int64               `json:"deptId"`        // 部门ID
	LoginTime     int64               `json:"loginTime"`     // 登录时间时间戳
	ExpireTime    int64               `json:"expireTime"`    // 过期时间时间戳
	LoginIp       string              `json:"loginIp"`       // 登录IP地址 x.x.x.x
	LoginLocation string              `json:"loginLocation"` // 登录地点 xx xx
	Browser       string              `json:"browser"`       // 浏览器类型
	OS            string              `json:"os"`            // 操作系统
	Permissions   []string            `json:"permissions"`   // 权限列表
	User          systemModel.SysUser `json:"user"`          // 用户信息
}
