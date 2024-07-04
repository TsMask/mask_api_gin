package model

// SysJob 调度任务信息表 sys_job
type SysJob struct {
	JobID          string `json:"jobId"`                             // 任务ID
	JobName        string `json:"jobName" binding:"required"`        // 任务名称
	JobGroup       string `json:"jobGroup" binding:"required"`       // 任务组名
	InvokeTarget   string `json:"invokeTarget" binding:"required"`   // 调用目标字符串
	TargetParams   string `json:"targetParams"`                      // 调用目标传入参数
	CronExpression string `json:"cronExpression" binding:"required"` // cron执行表达式
	MisfirePolicy  string `json:"misfirePolicy"`                     // 计划执行错误策略（1立即执行 2执行一次 3放弃执行）
	Concurrent     string `json:"concurrent"`                        // 是否并发执行（0禁止 1允许）
	Status         string `json:"status"`                            // 任务状态（0暂停 1正常）
	SaveLog        string `json:"saveLog"`                           // 是否记录任务日志（0不记录 1记录）
	CreateBy       string `json:"createBy"`                          // 创建者
	CreateTime     int64  `json:"createTime"`                        // 创建时间
	UpdateBy       string `json:"updateBy"`                          // 更新者
	UpdateTime     int64  `json:"updateTime"`                        // 更新时间
	Remark         string `json:"remark"`                            // 备注
}
