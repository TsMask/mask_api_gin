package model

// SysUser 用户对象 sys_user
type SysUser struct {
	UserID     string `json:"userId"`                      // 用户ID
	DeptID     string `json:"deptId"`                      // 部门ID
	UserName   string `json:"userName" binding:"required"` // 用户账号
	NickName   string `json:"nickName" binding:"required"` // 用户昵称
	UserType   string `json:"userType"`                    // 用户类型（sys系统用户）
	Email      string `json:"email"`                       // 用户邮箱
	Phone      string `json:"phone"`                       // 手机号码
	Sex        string `json:"sex"`                         // 用户性别（0未知 1男 2女）
	Avatar     string `json:"avatar"`                      // 头像地址
	Password   string `json:"-"`                           // 密码
	Status     string `json:"status"`                      // 账号状态（0停用 1正常）
	DelFlag    string `json:"delFlag"`                     // 删除标志（0存在 1删除）
	LoginIP    string `json:"loginIp"`                     // 最后登录IP
	LoginDate  int64  `json:"loginDate"`                   // 最后登录时间
	CreateBy   string `json:"createBy"`                    // 创建者
	CreateTime int64  `json:"createTime"`                  // 创建时间
	UpdateBy   string `json:"updateBy"`                    // 更新者
	UpdateTime int64  `json:"updateTime"`                  // 更新时间
	Remark     string `json:"remark"`                      // 备注

	// ====== 非数据库字段属性 ======

	Dept    SysDept   `json:"dept,omitempty" binding:"structonly"` // 部门对象
	Roles   []SysRole `json:"roles"`                               // 角色对象组
	RoleID  string    `json:"roleId,omitempty"`                    // 角色ID
	RoleIDs []string  `json:"roleIds,omitempty"`                   // 角色组
	PostIDs []string  `json:"postIds,omitempty"`                   // 岗位组
}
