package model

// SysUser 用户信息表
type SysUser struct {
	UserId     string `json:"userId" gorm:"column:user_id;primaryKey;type:int;autoIncrement"` // 用户ID
	DeptId     string `json:"deptId" gorm:"column:dept_id"`                                   // 部门ID
	UserName   string `json:"userName" gorm:"column:user_name"`                               // 用户账号
	Email      string `json:"email" gorm:"column:email"`                                      // 用户邮箱
	Phone      string `json:"phone" gorm:"column:phone"`                                      // 手机号码
	NickName   string `json:"nickName" gorm:"column:nick_name"`                               // 用户昵称
	Sex        string `json:"sex" gorm:"column:sex"`                                          // 用户性别（0未选择 1男 2女）
	Avatar     string `json:"avatar" gorm:"column:avatar"`                                    // 头像地址
	Passwd     string `json:"-" gorm:"column:passwd"`                                         // 密码
	StatusFlag string `json:"statusFlag" gorm:"column:status_flag"`                           // 账号状态（0停用 1正常）
	DelFlag    string `json:"-" gorm:"column:del_flag"`                                       // 删除标记（0存在 1删除）
	LoginIp    string `json:"loginIp" gorm:"column:login_ip"`                                 // 最后登录IP
	LoginTime  int64  `json:"loginTime" gorm:"column:login_time"`                             // 最后登录时间
	CreateBy   string `json:"createBy" gorm:"column:create_by"`                               // 创建者
	CreateTime int64  `json:"createTime" gorm:"column:create_time"`                           // 创建时间
	UpdateBy   string `json:"updateBy" gorm:"column:update_by"`                               // 更新者
	UpdateTime int64  `json:"updateTime" gorm:"column:update_time"`                           // 更新时间
	Remark     string `json:"remark" gorm:"column:remark"`                                    // 备注

	// ====== 非数据库字段属性 ======

	Dept    *SysDept   `json:"dept" gorm:"-"`              // 部门对象
	Roles   []*SysRole `json:"roles" gorm:"-"`             // 角色对象组
	RoleIds []string   `json:"roleIds,omitempty" gorm:"-"` // 角色组
	PostIds []string   `json:"postIds,omitempty" gorm:"-"` // 岗位组
}

// TableName 表名称
func (*SysUser) TableName() string {
	return "sys_user"
}
