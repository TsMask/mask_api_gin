package model

// SysRoleMenu 角色和菜单关联表
type SysRoleMenu struct {
	RoleId int64 `json:"roleId" gorm:"column:role_id"` // 角色ID
	MenuId int64 `json:"menuId" gorm:"column:menu_id"` // 菜单ID
}

// TableName 表名称
func (*SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
