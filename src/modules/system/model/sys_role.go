package model

// SysRole 角色信息表（用户 M:N 角色 sys_user_role；角色 M:N 菜单 sys_role_menu）
// role_key='system' 为系统指定超级管理员，鉴权中间件对其免校验
type SysRole struct {
	ID                    int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`                    // 角色ID
	CreateBy              string `json:"create_by" gorm:"column:create_by"`                               // 创建者
	CreateTime            int64  `json:"create_time" gorm:"column:create_time"`                           // 创建时间
	UpdateBy              string `json:"update_by" gorm:"column:update_by"`                               // 更新者
	UpdateTime            int64  `json:"update_time" gorm:"column:update_time"`                           // 更新时间
	Remark                string `json:"remark" gorm:"column:remark"`                                     // 备注
	RoleName              string `json:"role_name" gorm:"column:role_name" binding:"required"`            // 角色名称
	RoleKey               string `json:"role_key" gorm:"column:role_key" binding:"required"`              // 角色键值
	RoleSort              int64  `json:"role_sort" gorm:"column:role_sort"`                               // 显示顺序
	RoleStatus            string `json:"role_status" gorm:"column:role_status"`                           // 角色状态（0停用 1正常）
	RoleDel               string `json:"-" gorm:"column:role_del"`                                        // 删除标记（0存在 1删除）
	RoleDataScope         string `json:"role_data_scope" gorm:"column:role_data_scope"`                   // 数据范围（1全部 2自定 3本部门 4本部门及以下 5仅本人）
	RoleMenuCheckStrictly string `json:"role_menu_check_strictly" gorm:"column:role_menu_check_strictly"` // 菜单树选择项是否关联显示（0不关联 1互相）
	RoleDeptCheckStrictly string `json:"role_dept_check_strictly" gorm:"column:role_dept_check_strictly"` // 部门树选择项是否关联显示（0不关联 1互相）

	// ====== 非数据库字段属性 ======

	MenuIds []int64 `json:"menu_ids,omitempty" gorm:"-"` // 菜单组
	DeptIds []int64 `json:"dept_ids,omitempty" gorm:"-"` // 部门组（数据权限）
}

// TableName 表名称
func (*SysRole) TableName() string {
	return "sys_role"
}
