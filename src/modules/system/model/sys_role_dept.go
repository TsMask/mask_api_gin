package model

// SysRoleDept 角色和部门关联表
type SysRoleDept struct {
	RoleId int64 `json:"role_id" gorm:"column:role_id;primaryKey"` // 角色ID
	DeptId int64 `json:"dept_id" gorm:"column:dept_id;primaryKey"` // 部门ID
}

// TableName 表名称
func (*SysRoleDept) TableName() string {
	return "sys_role_dept"
}
