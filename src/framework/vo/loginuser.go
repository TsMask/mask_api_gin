package vo

import systemModel "mask_api_gin/src/modules/system/model"

// LoginUser 登录用户身份权限信息对象
type LoginUser struct {
	UserId        int64               `json:"userId"`        // 用户ID
	DeptId        int64               `json:"deptId"`        // 部门ID
	UUID          string              `json:"uuid"`          // 用户唯一标识
	LoginTime     int64               `json:"loginTime"`     // 登录时间时间戳
	ExpireTime    int64               `json:"expireTime"`    // 过期时间时间戳
	LoginIp       string              `json:"loginIp"`       // 登录IP地址 x.x.x.x
	LoginLocation string              `json:"loginLocation"` // 登录地点 xx xx
	Browser       string              `json:"browser"`       // 浏览器类型
	OS            string              `json:"os"`            // 操作系统
	Permissions   []string            `json:"permissions"`   // 权限列表
	User          systemModel.SysUser `json:"user"`          // 用户信息
}
