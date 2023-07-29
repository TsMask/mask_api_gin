package model

// SysJobLog 定时任务调度日志表 sys_job_log
type SysJobLog struct {
	// 日志序号
	JobLogID string `json:"jobLogId"`
	// 任务名称
	JobName string `json:"jobName"`
	// 任务组名
	JobGroup string `json:"jobGroup"`
	// 调用目标字符串
	InvokeTarget string `json:"invokeTarget"`
	// 调用目标传入参数
	TargetParams string `json:"targetParams"`
	// 日志信息
	JobMsg string `json:"jobMsg"`
	// 执行状态（0失败 1正常）
	Status string `json:"status"`
	// 创建时间
	CreateTime int64 `json:"createTime"`
}
