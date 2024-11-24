package model

// SysUserRole 用户和角色关联表
type SysUserRole struct {
	UserId string `json:"userId" gorm:"column:user_id"` // 用户ID
	RoleId string `json:"roleId" gorm:"column:role_id"` // 角色ID
}

// TableName 表名称
func (*SysUserRole) TableName() string {
	return "sys_user_role"
}
