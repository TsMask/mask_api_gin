package model

// SysLogLogin 系统登录日志表（用 login_user_name 字符串软关联 sys_user.user_name）
type SysLogLogin struct {
	ID            int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"` // 登录ID
	CreateBy      string `json:"createBy" gorm:"column:create_by"`             // 创建者
	CreateTime    int64  `json:"createTime" gorm:"column:create_time"`         // 创建时间
	LoginUserName string `json:"loginUserName" gorm:"column:login_user_name"`  // 用户账号
	LoginIp       string `json:"loginIp" gorm:"column:login_ip"`               // 登录IP地址
	LoginLocation string `json:"loginLocation" gorm:"column:login_location"`   // 登录地点
	LoginBrowser  string `json:"loginBrowser" gorm:"column:login_browser"`     // 浏览器类型
	LoginOs       string `json:"loginOs" gorm:"column:login_os"`               // 操作系统
	LoginStatus   string `json:"loginStatus" gorm:"column:login_status"`       // 登录状态（0失败 1成功）
	LoginMsg      string `json:"loginMsg" gorm:"column:login_msg"`             // 提示消息
}

// TableName 表名称
func (*SysLogLogin) TableName() string {
	return "sys_log_login"
}
