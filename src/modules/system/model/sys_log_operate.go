package model

// SysLogOperate 系统操作日志表
type SysLogOperate struct {
	ID             string `json:"id" gorm:"column:id;primary_key"`               // 操作ID
	Title          string `json:"title" gorm:"column:title"`                     // 模块标题
	BusinessType   string `json:"businessType" gorm:"column:business_type"`      // 业务类型（0其它 1新增 2修改 3删除 4授权 5导出 6导入 7强退 8清空数据）
	OperaUrl       string `json:"operaUrl" gorm:"column:opera_url"`              // 请求URL
	OperaUrlMethod string `json:"operaUrlMethod" gorm:"column:opera_url_method"` // 请求方式
	OperaIp        string `json:"operaIp" gorm:"column:opera_ip"`                // 主机地址
	OperaLocation  string `json:"operaLocation" gorm:"column:opera_location"`    // 操作地点
	OperaParam     string `json:"operaParam" gorm:"column:opera_param"`          // 请求参数
	OperaMsg       string `json:"operaMsg" gorm:"column:opera_msg"`              // 操作消息
	OperaMethod    string `json:"operaMethod" gorm:"column:opera_method"`        // 方法名称
	OperaBy        string `json:"operaBy" gorm:"column:opera_by"`                // 操作人员
	OperaTime      int64  `json:"operaTime" gorm:"column:opera_time"`            // 操作时间
	StatusFlag     string `json:"statusFlag" gorm:"column:status_flag"`          // 操作状态（0异常 1正常）
	CostTime       int64  `json:"costTime" gorm:"column:cost_time"`              // 消耗时间（毫秒）
}

// TableName 表名称
func (*SysLogOperate) TableName() string {
	return "sys_log_operate"
}
