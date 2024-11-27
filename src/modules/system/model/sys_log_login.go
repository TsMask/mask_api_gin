package model

// SysLogLogin 系统登录日志表
type SysLogLogin struct {
	ID            string `json:"id" gorm:"column:id;primaryKey;type:int;autoIncrement"` // 登录ID
	UserId        int64  `json:"userId" gorm:"column:user_id"`                          // 用户ID
	UserName      string `json:"userName" gorm:"column:user_name"`                      // 用户账号
	LoginIp       string `json:"loginIp" gorm:"column:login_ip"`                        // 登录IP地址
	LoginLocation string `json:"loginLocation" gorm:"column:login_location"`            // 登录地点
	Browser       string `json:"browser" gorm:"column:browser"`                         // 浏览器类型
	OS            string `json:"os" gorm:"column:os"`                                   // 操作系统
	StatusFlag    string `json:"statusFlag" gorm:"column:status_flag"`                  // 登录状态（0失败 1成功）
	Msg           string `json:"msg" gorm:"column:msg"`                                 // 提示消息
	LoginTime     int64  `json:"loginTime" gorm:"column:login_time"`                    // 登录时间
}

// TableName 表名称
func (*SysLogLogin) TableName() string {
	return "sys_log_login"
}
