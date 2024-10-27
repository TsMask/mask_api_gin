package model

// SysJobLog 调度任务调度日志表
type SysJobLog struct {
	JobLogId     string `json:"job_log_id" gorm:"column:job_log_id"`       // 任务日志ID
	JobName      string `json:"job_name" gorm:"column:job_name"`           // 任务名称
	JobGroup     string `json:"job_group" gorm:"column:job_group"`         // 任务组名
	InvokeTarget string `json:"invoke_target" gorm:"column:invoke_target"` // 调用目标字符串
	TargetParams string `json:"target_params" gorm:"column:target_params"` // 调用目标传入参数
	JobMsg       string `json:"job_msg" gorm:"column:job_msg"`             // 日志信息
	Status       string `json:"status" gorm:"column:status"`               // 执行状态（0失败 1正常）
	CreateTime   int64  `json:"create_time" gorm:"column:create_time"`     // 创建时间
	CostTime     int64  `json:"cost_time" gorm:"column:cost_time"`         // 消耗时间（毫秒）
}

// TableName 表名称
func (*SysJobLog) TableName() string {
	return "sys_job_log"
}
