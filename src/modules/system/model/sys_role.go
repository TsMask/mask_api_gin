package model

// SysRole 角色对象 sys_role
type SysRole struct {
	// 角色ID
	RoleID string `json:"roleId"`
	// 角色名称
	RoleName string `json:"roleName" binding:"required"`
	// 角色键值
	RoleKey string `json:"roleKey" binding:"required"`
	// 显示顺序
	RoleSort int `json:"roleSort"`
	// 数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限 5：仅本人数据权限）
	DataScope string `json:"dataScope"`
	// 菜单树选择项是否关联显示（0：父子不互相关联显示 1：父子互相关联显示）
	MenuCheckStrictly string `json:"menuCheckStrictly"`
	// 部门树选择项是否关联显示（0：父子不互相关联显示 1：父子互相关联显示）
	DeptCheckStrictly string `json:"deptCheckStrictly"`
	// 角色状态（0停用 1正常）
	Status string `json:"status"`
	// 删除标志（0代表存在 1代表删除）
	DelFlag string `json:"delFlag"`
	// 创建者
	CreateBy string `json:"createBy"`
	// 创建时间
	CreateTime int64 `json:"createTime"`
	// 更新者
	UpdateBy string `json:"updateBy"`
	// 更新时间
	UpdateTime int64 `json:"updateTime"`
	// 备注
	Remark string `json:"remark"`

	// ====== 非数据库字段属性 ======

	// 菜单组
	MenuIds []string `json:"menuIds,omitempty"`
	// 部门组（数据权限）
	DeptIds []string `json:"deptIds,omitempty"`
	// 角色菜单权限
	Permissions []string `json:"permissions,omitempty"`
	// 用户是否存在此角色标识 默认不存在
	Flag bool `json:"flag,omitempty"`
}
