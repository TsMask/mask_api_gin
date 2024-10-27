package model

// SysLogOperate 系统操作日志表
type SysLogOperate struct {
	OperaId       string `json:"opera_id" gorm:"column:opera_id;primary_key"` // 日志主键
	Title         string `json:"title" gorm:"column:title"`                   // 模块标题
	BusinessType  string `json:"business_type" gorm:"column:business_type"`   // 业务类型（0其它 1新增 2修改 3删除 4授权 5导出 6导入 7强退 8清空数据）
	Method        string `json:"method" gorm:"column:method"`                 // 方法名称
	RequestMethod string `json:"request_method" gorm:"column:request_method"` // 请求方式
	OperatorType  string `json:"operator_type" gorm:"column:operator_type"`   // 操作人员类别（0其它 1后台用户 2手机端用户）
	OperaName     string `json:"opera_name" gorm:"column:opera_name"`         // 操作人员
	DeptName      string `json:"dept_name" gorm:"column:dept_name"`           // 部门名称
	OperaUrl      string `json:"opera_url" gorm:"column:opera_url"`           // 请求URL
	OperaIp       string `json:"opera_ip" gorm:"column:opera_ip"`             // 主机地址
	OperaLocation string `json:"opera_location" gorm:"column:opera_location"` // 操作地点
	OperaParam    string `json:"opera_param" gorm:"column:opera_param"`       // 请求参数
	OperaMsg      string `json:"opera_msg" gorm:"column:opera_msg"`           // 操作消息
	Status        string `json:"status" gorm:"column:status"`                 // 操作状态（0异常 1正常）
	OperaTime     int64  `json:"opera_time" gorm:"column:opera_time"`         // 操作时间
	CostTime      int64  `json:"cost_time" gorm:"column:cost_time"`           // 消耗时间（毫秒）
}

// TableName 表名称
func (*SysLogOperate) TableName() string {
	return "sys_log_operate"
}
