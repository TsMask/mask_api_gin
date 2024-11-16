package model

// SysJobLog 调度任务调度日志表
type SysJobLog struct {
	LogId        int64  `json:"logId" gorm:"column:log_id;primary_key"`   // 任务日志ID
	JobName      string `json:"jobName" gorm:"column:job_name"`           // 任务名称
	JobGroup     string `json:"jobGroup" gorm:"column:job_group"`         // 任务组名
	InvokeTarget string `json:"invokeTarget" gorm:"column:invoke_target"` // 调用目标字符串
	TargetParams string `json:"targetParams" gorm:"column:target_params"` // 调用目标传入参数
	JobMsg       string `json:"jobMsg" gorm:"column:job_msg"`             // 日志信息
	StatusFlag   string `json:"statusFlag" gorm:"column:status_flag"`     // 执行状态（0失败 1正常）
	CreateTime   int64  `json:"createTime" gorm:"column:create_time"`     // 创建时间
	CostTime     int64  `json:"costTime" gorm:"column:cost_time"`         // 消耗时间（毫秒）
}

// TableName 表名称
func (*SysJobLog) TableName() string {
	return "sys_job_log"
}
