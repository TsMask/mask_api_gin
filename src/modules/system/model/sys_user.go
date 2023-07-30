package model

// SysUser 用户对象 sys_user
type SysUser struct {
	// 用户ID
	UserID string `json:"userId"`
	// 部门ID
	DeptID string `json:"deptId"`
	// 用户账号
	UserName string `json:"userName" binding:"required"`
	// 用户昵称
	NickName string `json:"nickName" binding:"required"`
	// 用户类型（sys系统用户）
	UserType string `json:"userType"`
	// 用户邮箱
	Email string `json:"email"`
	// 手机号码
	PhoneNumber string `json:"phonenumber"`
	// 用户性别（0未知 1男 2女）
	Sex string `json:"sex"`
	// 头像地址
	Avatar string `json:"avatar"`
	// 密码
	Password string `json:"-"`
	// 帐号状态（0停用 1正常）
	Status string `json:"status"`
	// 删除标志（0代表存在 1代表删除）
	DelFlag string `json:"delFlag"`
	// 最后登录IP
	LoginIP string `json:"loginIp"`
	// 最后登录时间
	LoginDate int64 `json:"loginDate"`
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

	// 部门对象
	Dept SysDept `json:"dept,omitempty" binding:"structonly"`
	// 角色对象组
	Roles []SysRole `json:"roles"`
	// 角色ID
	RoleID string `json:"roleId,omitempty"`
	// 角色组
	RoleIDs []string `json:"roleIds,omitempty"`
	// 岗位组
	PostIDs []string `json:"postIds,omitempty"`
}
