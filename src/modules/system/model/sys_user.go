package model

// SysUser 用户信息表
type SysUser struct {
	UserId     string `json:"user_id" gorm:"column:user_id;primary_key"`            // 用户ID
	DeptId     string `json:"dept_id" gorm:"column:dept_id"`                        // 部门ID
	UserName   string `json:"user_name" gorm:"column:user_name" binding:"required"` // 用户账号
	NickName   string `json:"nick_name" gorm:"column:nick_name" binding:"required"` // 用户昵称
	UserType   string `json:"user_type" gorm:"column:user_type"`                    // 用户类型（sys系统用户）
	Email      string `json:"email" gorm:"column:email"`                            // 用户邮箱
	Phone      string `json:"phone" gorm:"column:phone"`                            // 手机号码
	Sex        string `json:"sex" gorm:"column:sex"`                                // 用户性别（0未知 1男 2女）
	Avatar     string `json:"avatar" gorm:"column:avatar"`                          // 头像地址
	Password   string `json:"password" gorm:"column:password"`                      // 密码
	Status     string `json:"status" gorm:"column:status"`                          // 账号状态（0停用 1正常）
	DelFlag    string `json:"del_flag" gorm:"column:del_flag"`                      // 删除标志（0存在 1删除）
	LoginIp    string `json:"login_ip" gorm:"column:login_ip"`                      // 最后登录IP
	LoginDate  int64  `json:"login_date" gorm:"column:login_date"`                  // 最后登录时间
	CreateBy   string `json:"create_by" gorm:"column:create_by"`                    // 创建者
	CreateTime int64  `json:"create_time" gorm:"column:create_time"`                // 创建时间
	UpdateBy   string `json:"update_by" gorm:"column:update_by"`                    // 更新者
	UpdateTime int64  `json:"update_time" gorm:"column:update_time"`                // 更新时间
	Remark     string `json:"remark" gorm:"column:remark"`                          // 备注

	// ====== 非数据库字段属性 ======

	Dept    SysDept   `json:"dept,omitempty" gorm:"-" binding:"structonly"` // 部门对象
	Roles   []SysRole `json:"roles" gorm:"-"`                               // 角色对象组
	RoleIds []string  `json:"role_ids,omitempty" gorm:"-"`                  // 角色组
	PostIds []string  `json:"post_ids,omitempty" gorm:"-"`                  // 岗位组
}

// TableName 表名称
func (*SysUser) TableName() string {
	return "sys_user"
}
