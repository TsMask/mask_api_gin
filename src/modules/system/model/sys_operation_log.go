package model

// SysOperationLog 系统操作日志表（用 create_by 字符串软关联 sys_user.user_name）
type SysOperationLog struct {
	ID                int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`          // 操作ID
	CreateBy          string `json:"create_by" gorm:"column:create_by"`                     // 创建者
	CreateTime        int64  `json:"create_time" gorm:"column:create_time"`                 // 创建时间
	Title             string `json:"title" gorm:"column:title"`                             // 模块标题
	OperationType     string `json:"operation_type" gorm:"column:operation_type"`           // 业务类型（0其它 1新增 2修改 3删除 4授权 5导出 6导入 7强退 8清空数据）
	OperationUrl      string `json:"operation_url" gorm:"column:operation_url"`             // 请求URL
	OperationMethod   string `json:"operation_method" gorm:"column:operation_method"`       // 请求方式
	OperationIp       string `json:"operation_ip" gorm:"column:operation_ip"`               // 主机地址
	OperationLocation string `json:"operation_location" gorm:"column:operation_location"`   // 操作地点
	OperationParam    string `json:"operation_param" gorm:"column:operation_param"`         // 请求参数
	OperationFunction string `json:"operation_function" gorm:"column:operation_function"`   // 方法函数方法
	OperationCostTime int64  `json:"operation_cost_time" gorm:"column:operation_cost_time"` // 消耗时间（毫秒）
	OperationStatus   string `json:"operation_status" gorm:"column:operation_status"`       // 操作状态（0异常 1正常）
	OperationMsg      string `json:"operation_msg" gorm:"column:operation_msg"`             // 操作消息
}

// TableName 表名称
func (*SysOperationLog) TableName() string {
	return "sys_operation_log"
}
