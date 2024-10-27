package model

// SysJob 调度任务调度表
type SysJob struct {
	JobId          string `json:"job_id" gorm:"column:job_id"`                                      // 任务ID
	JobName        string `json:"job_name" gorm:"column:job_name" binding:"required"`               // 任务名称
	JobGroup       string `json:"job_group" gorm:"column:job_group" binding:"required"`             // 任务组名
	InvokeTarget   string `json:"invoke_target" gorm:"column:invoke_target" binding:"required"`     // 调用目标字符串
	TargetParams   string `json:"target_params" gorm:"column:target_params"`                        // 调用目标传入参数
	CronExpression string `json:"cron_expression" gorm:"column:cron_expression" binding:"required"` // cron执行表达式
	MisfirePolicy  string `json:"misfire_policy" gorm:"column:misfire_policy"`                      // 计划执行错误策略（1立即执行 2执行一次 3放弃执行）
	Concurrent     string `json:"concurrent" gorm:"column:concurrent"`                              // 是否并发执行（0禁止 1允许）
	Status         string `json:"status" gorm:"column:status"`                                      // 任务状态（0暂停 1正常）
	SaveLog        string `json:"save_log" gorm:"column:save_log"`                                  // 是否记录任务日志（0不记录 1记录）
	CreateBy       string `json:"create_by" gorm:"column:create_by"`                                // 创建者
	CreateTime     int64  `json:"create_time" gorm:"column:create_time"`                            // 创建时间
	UpdateBy       string `json:"update_by" gorm:"column:update_by"`                                // 更新者
	UpdateTime     int64  `json:"update_time" gorm:"column:update_time"`                            // 更新时间
	Remark         string `json:"remark" gorm:"column:remark"`                                      // 备注
}

// TableName 表名称
func (*SysJob) TableName() string {
	return "sys_job"
}
