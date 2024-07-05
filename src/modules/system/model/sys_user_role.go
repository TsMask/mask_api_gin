package model

// SysUserRole 用户和角色关联对象 sys_user_role
type SysUserRole struct {
	UserID string `json:"userId"` // 用户ID
	RoleID string `json:"roleId"` // 角色ID
}
