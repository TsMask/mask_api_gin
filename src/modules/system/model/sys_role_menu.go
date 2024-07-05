package model

// SysRoleMenu 角色和菜单关联对象 sys_role_menu
type SysRoleMenu struct {
	RoleID string `json:"roleId"` // 角色ID
	MenuID string `json:"menuId"` // 菜单ID
}
