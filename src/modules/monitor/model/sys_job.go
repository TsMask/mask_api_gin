package model

// SysJob 调度任务调度表
type SysJob struct {
	JobId          string `json:"jobId" gorm:"column:job_id;primaryKey;type:int;autoIncrement"`    // 任务ID
	JobName        string `json:"jobName" gorm:"column:job_name" binding:"required"`               // 任务名称
	JobGroup       string `json:"jobGroup" gorm:"column:job_group" binding:"required"`             // 任务组名
	InvokeTarget   string `json:"invokeTarget" gorm:"column:invoke_target" binding:"required"`     // 调用目标字符串
	TargetParams   string `json:"targetParams" gorm:"column:target_params"`                        // 调用目标传入参数
	CronExpression string `json:"cronExpression" gorm:"column:cron_expression" binding:"required"` // cron执行表达式
	MisfirePolicy  string `json:"misfirePolicy" gorm:"column:misfire_policy"`                      // 计划执行错误策略（1立即执行 2执行一次 3放弃执行）
	Concurrent     string `json:"concurrent" gorm:"column:concurrent"`                             // 是否并发执行（0禁止 1允许）
	StatusFlag     string `json:"statusFlag" gorm:"column:status_flag"`                            // 任务状态（0暂停 1正常）
	SaveLog        string `json:"saveLog" gorm:"column:save_log"`                                  // 是否记录任务日志（0不记录 1记录）
	CreateBy       string `json:"createBy" gorm:"column:create_by"`                                // 创建者
	CreateTime     int64  `json:"createTime" gorm:"column:create_time"`                            // 创建时间
	UpdateBy       string `json:"updateBy" gorm:"column:update_by"`                                // 更新者
	UpdateTime     int64  `json:"updateTime" gorm:"column:update_time"`                            // 更新时间
	Remark         string `json:"remark" gorm:"column:remark"`                                     // 备注
}

// TableName 表名称
func (*SysJob) TableName() string {
	return "sys_job"
}
