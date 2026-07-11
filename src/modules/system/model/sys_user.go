package model

// SysUser 用户信息表（用户 N:1 部门 dept_id；密码 bcrypt 存储）
type SysUser struct {
	ID            int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"` // 用户ID
	CreateBy      string `json:"createBy" gorm:"column:create_by"`             // 创建者
	CreateTime    int64  `json:"createTime" gorm:"column:create_time"`         // 创建时间
	UpdateBy      string `json:"updateBy" gorm:"column:update_by"`             // 更新者
	UpdateTime    int64  `json:"updateTime" gorm:"column:update_time"`         // 更新时间
	Remark        string `json:"remark" gorm:"column:remark"`                  // 备注
	UserName      string `json:"userName" gorm:"column:user_name"`             // 用户账号
	UserNick      string `json:"userNick" gorm:"column:user_nick"`             // 用户昵称
	UserEmail     string `json:"userEmail" gorm:"column:user_email"`           // 用户邮箱
	UserPhone     string `json:"userPhone" gorm:"column:user_phone"`           // 手机号码
	UserGender    string `json:"userGender" gorm:"column:user_gender"`         // 用户性别（0未选择 1男 2女）
	UserAvatar    string `json:"userAvatar" gorm:"column:user_avatar"`         // 头像地址
	UserPassword  string `json:"-" gorm:"column:user_password"`                // 密码
	UserLoginIp   string `json:"userLoginIp" gorm:"column:user_login_ip"`      // 最后登录IP
	UserLoginTime int64  `json:"userLoginTime" gorm:"column:user_login_time"`  // 最后登录时间
	UserStatus    string `json:"userStatus" gorm:"column:user_status"`         // 账号状态（0停用 1正常）
	UserDel       string `json:"-" gorm:"column:user_del"`                     // 删除标记（0存在 1删除）
	DeptId        int64  `json:"deptId" gorm:"column:dept_id"`                 // 部门ID

	// ====== 非数据库字段属性 ======

	Dept    *SysDept   `json:"dept" gorm:"-"`              // 部门对象
	Roles   []*SysRole `json:"roles" gorm:"-"`             // 角色对象组
	RoleIds []int64    `json:"roleIds,omitempty" gorm:"-"` // 角色组
	PostIds []int64    `json:"postIds,omitempty" gorm:"-"` // 岗位组
}

// TableName 表名称
func (*SysUser) TableName() string {
	return "sys_user"
}
