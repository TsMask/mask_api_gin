package model

// SysRoleDept 角色和部门关联表
type SysRoleDept struct {
	RoleId string `json:"role_id" gorm:"column:role_id"` // 角色ID
	DeptId string `json:"dept_id" gorm:"column:dept_id"` // 部门ID
}

// TableName 表名称
func (*SysRoleDept) TableName() string {
	return "sys_role_dept"
}
