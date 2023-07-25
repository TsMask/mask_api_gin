package model

// SysJob 调度任务信息表 sys_job
type SysJob struct {
	// 任务ID
	JobID string `json:"jobId"`
	// 任务名称
	JobName string `json:"jobName" binding:"required"`
	// 任务组名
	JobGroup string `json:"jobGroup" binding:"required"`
	// 调用目标字符串
	InvokeTarget string `json:"invokeTarget" binding:"required"`
	// 调用目标传入参数
	TargetParams string `json:"targetParams"`
	// cron执行表达式
	CronExpression string `json:"cronExpression" binding:"required"`
	// 计划执行错误策略（1立即执行 2执行一次 3放弃执行）
	MisfirePolicy string `json:"misfirePolicy"`
	// 是否并发执行（0禁止 1允许）
	Concurrent string `json:"concurrent"`
	// 任务状态（0暂停 1正常）
	Status string `json:"status"`
	// 创建者
	CreateBy string `json:"createBy"`
	// 创建时间
	CreateTime int64 `json:"createTime"`
	// 更新者
	UpdateBy string `json:"updateBy"`
	// 更新时间
	UpdateTime int64 `json:"updateTime"`
	// 备注
	Remark string `json:"remark"`
}
