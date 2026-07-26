package model

// SysRoleMenu 角色和菜单关联表
type SysRoleMenu struct {
	RoleId int64 `json:"role_id" gorm:"column:role_id;primaryKey"` // 角色ID
	MenuId int64 `json:"menu_id" gorm:"column:menu_id;primaryKey"` // 菜单ID
}

// TableName 表名称
func (*SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
