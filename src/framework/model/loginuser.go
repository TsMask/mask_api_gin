package model

import systemModel "mask_api_gin/src/modules/system/model"

// LoginUser 登录用户身份权限信息对象
type LoginUser struct {
	// UserID 用户ID
	UserID string `json:"userId"`

	// DeptID 部门ID
	DeptID string `json:"deptId"`

	// UUID 用户唯一标识
	UUID string `json:"uuid"`

	// LoginTime 登录时间时间戳
	LoginTime int64 `json:"loginTime"`

	// ExpireTime 过期时间时间戳
	ExpireTime int64 `json:"expireTime"`

	// IPAddr 登录IP地址 x.x.x.x
	IPAddr string `json:"ipaddr"`

	// LoginLocation 登录地点 xx xx
	LoginLocation string `json:"loginLocation"`

	// Browser 浏览器类型
	Browser string `json:"browser"`

	// OS 操作系统
	OS string `json:"os"`

	// Permissions 权限列表
	Permissions []string `json:"permissions"`

	// User 用户信息
	User systemModel.SysUser `json:"user"`
}
