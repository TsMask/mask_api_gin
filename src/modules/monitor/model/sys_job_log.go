package model

// SysJobLog 定时任务调度日志表 sys_job_log
type SysJobLog struct {
	JobLogID     string `json:"jobLogId"`     // 日志序号
	JobName      string `json:"jobName"`      // 任务名称
	JobGroup     string `json:"jobGroup"`     // 任务组名
	InvokeTarget string `json:"invokeTarget"` // 调用目标字符串
	TargetParams string `json:"targetParams"` // 调用目标传入参数
	JobMsg       string `json:"jobMsg"`       // 日志信息
	Status       string `json:"status"`       // 执行状态（0失败 1正常）
	CreateTime   int64  `json:"createTime"`   // 创建时间
	CostTime     int64  `json:"costTime"`     // 消耗时间（毫秒）
}
