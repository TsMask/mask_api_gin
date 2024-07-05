package model

// SysRole 角色对象 sys_role
type SysRole struct {
	RoleID            string `json:"roleId"`                      // 角色ID
	RoleName          string `json:"roleName" binding:"required"` // 角色名称
	RoleKey           string `json:"roleKey" binding:"required"`  // 角色键值
	RoleSort          int    `json:"roleSort"`                    // 显示顺序
	DataScope         string `json:"dataScope"`                   // 数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限 5：仅本人数据权限）
	MenuCheckStrictly string `json:"menuCheckStrictly"`           // 菜单树选择项是否关联显示（0：父子不互相关联显示 1：父子互相关联显示）
	DeptCheckStrictly string `json:"deptCheckStrictly"`           // 部门树选择项是否关联显示（0：父子不互相关联显示 1：父子互相关联显示）
	Status            string `json:"status"`                      // 角色状态（0停用 1正常）
	DelFlag           string `json:"delFlag"`                     // 删除标志（0存在 1删除）
	CreateBy          string `json:"createBy"`                    // 创建者
	CreateTime        int64  `json:"createTime"`                  // 创建时间
	UpdateBy          string `json:"updateBy"`                    // 更新者
	UpdateTime        int64  `json:"updateTime"`                  // 更新时间
	Remark            string `json:"remark"`                      // 备注

	// ====== 非数据库字段属性 ======

	MenuIds []string `json:"menuIds,omitempty"` // 菜单组
	DeptIds []string `json:"deptIds,omitempty"` // 部门组（数据权限）
}
