package model

// SysLoginLog 系统登录日志表（用 login_user_name 字符串软关联 sys_user.user_name）
type SysLoginLog struct {
	ID            int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`  // 登录ID
	CreateBy      string `json:"create_by" gorm:"column:create_by"`             // 创建者
	CreateTime    int64  `json:"create_time" gorm:"column:create_time"`         // 创建时间
	LoginUserName string `json:"login_user_name" gorm:"column:login_user_name"` // 用户账号
	LoginIp       string `json:"login_ip" gorm:"column:login_ip"`               // 登录IP地址
	LoginLocation string `json:"login_location" gorm:"column:login_location"`   // 登录地点
	LoginBrowser  string `json:"login_browser" gorm:"column:login_browser"`     // 浏览器类型
	LoginOs       string `json:"login_os" gorm:"column:login_os"`               // 操作系统
	LoginStatus   string `json:"login_status" gorm:"column:login_status"`       // 登录状态（0失败 1成功）
	LoginMsg      string `json:"login_msg" gorm:"column:login_msg"`             // 提示消息
}

// TableName 表名称
func (*SysLoginLog) TableName() string {
	return "sys_login_log"
}
