package model

// SysRoleDept 角色和部门关联对象 sys_role_dept
type SysRoleDept struct {
	RoleID string `json:"roleId"` // 角色ID
	DeptID string `json:"deptId"` // 部门ID
}
