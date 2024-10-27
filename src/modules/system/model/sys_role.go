package model

// SysRole 角色信息表
type SysRole struct {
	RoleId            string `json:"role_id" gorm:"column:role_id;primary_key"`             // 角色ID
	RoleName          string `json:"role_name" gorm:"column:role_name" binding:"required"`  // 角色名称
	RoleKey           string `json:"role_key" gorm:"column:role_key" binding:"required"`    // 角色键值
	RoleSort          int64  `json:"role_sort" gorm:"column:role_sort"`                     // 显示顺序
	DataScope         string `json:"data_scope" gorm:"column:data_scope"`                   // 数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限 5：仅本人数据权限）
	MenuCheckStrictly string `json:"menu_check_strictly" gorm:"column:menu_check_strictly"` // 菜单树选择项是否关联显示（0：父子不互相关联显示 1：父子互相关联显示）
	DeptCheckStrictly string `json:"dept_check_strictly" gorm:"column:dept_check_strictly"` // 部门树选择项是否关联显示（0：父子不互相关联显示 1：父子互相关联显示 ）
	Status            string `json:"status" gorm:"column:status"`                           // 角色状态（0停用 1正常）
	DelFlag           string `json:"del_flag" gorm:"column:del_flag"`                       // 删除标志（0存在 1删除）
	CreateBy          string `json:"create_by" gorm:"column:create_by"`                     // 创建者
	CreateTime        int64  `json:"create_time" gorm:"column:create_time"`                 // 创建时间
	UpdateBy          string `json:"update_by" gorm:"column:update_by"`                     // 更新者
	UpdateTime        int64  `json:"update_time" gorm:"column:update_time"`                 // 更新时间
	Remark            string `json:"remark" gorm:"column:remark"`                           // 备注

	// ====== 非数据库字段属性 ======

	MenuIds []string `json:"menuIds,omitempty" gorm:"-"` // 菜单组
	DeptIds []string `json:"deptIds,omitempty" gorm:"-"` // 部门组（数据权限）
}

// TableName 表名称
func (*SysRole) TableName() string {
	return "sys_role"
}
