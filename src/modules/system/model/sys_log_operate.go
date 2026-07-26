package model

// SysLogOperate 系统操作日志表（用 create_by 字符串软关联 sys_user.user_name）
type SysLogOperate struct {
	ID              int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`      // 操作ID
	CreateBy        string `json:"create_by" gorm:"column:create_by"`                 // 创建者
	CreateTime      int64  `json:"create_time" gorm:"column:create_time"`             // 创建时间
	Title           string `json:"title" gorm:"column:title"`                         // 模块标题
	OperateType     string `json:"operate_type" gorm:"column:operate_type"`           // 业务类型（0其它 1新增 2修改 3删除 4授权 5导出 6导入 7强退 8清空数据）
	OperateUrl      string `json:"operate_url" gorm:"column:operate_url"`             // 请求URL
	OperateMethod   string `json:"operate_method" gorm:"column:operate_method"`       // 请求方式
	OperateIp       string `json:"operate_ip" gorm:"column:operate_ip"`               // 主机地址
	OperateLocation string `json:"operate_location" gorm:"column:operate_location"`   // 操作地点
	OperateParam    string `json:"operate_param" gorm:"column:operate_param"`         // 请求参数
	OperateFunction string `json:"operate_function" gorm:"column:operate_function"`   // 方法函数方法
	OperateCostTime int64  `json:"operate_cost_time" gorm:"column:operate_cost_time"` // 消耗时间（毫秒）
	OperateStatus   string `json:"operate_status" gorm:"column:operate_status"`       // 操作状态（0异常 1正常）
	OperateMsg      string `json:"operate_msg" gorm:"column:operate_msg"`             // 操作消息
}

// TableName 表名称
func (*SysLogOperate) TableName() string {
	return "sys_log_operate"
}
