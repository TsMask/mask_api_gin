package model

// SysJob 调度任务调度表
type SysJob struct {
	ID               int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`                       // 任务ID
	CreateBy         string `json:"createBy" gorm:"column:create_by"`                                   // 创建者
	CreateTime       int64  `json:"createTime" gorm:"column:create_time"`                               // 创建时间
	UpdateBy         string `json:"updateBy" gorm:"column:update_by"`                                   // 更新者
	UpdateTime       int64  `json:"updateTime" gorm:"column:update_time"`                               // 更新时间
	Remark           string `json:"remark" gorm:"column:remark"`                                        // 备注
	JobName          string `json:"jobName" gorm:"column:job_name" binding:"required"`                  // 任务名称
	JobGroup         string `json:"jobGroup" gorm:"column:job_group" binding:"required"`                // 任务组名
	JobInvokeTarget  string `json:"jobInvokeTarget" gorm:"column:job_invoke_target" binding:"required"` // 调用目标字符串
	JobTargetParams  string `json:"jobTargetParams" gorm:"column:job_target_params"`                    // 调用目标传入参数
	JobCron          string `json:"jobCron" gorm:"column:job_cron" binding:"required"`                  // cron执行表达式
	JobMisfirePolicy string `json:"jobMisfirePolicy" gorm:"column:job_misfire_policy"`                  // 计划执行错误策略（1立即执行 2执行一次 3放弃执行）
	JobConcurrent    string `json:"jobConcurrent" gorm:"column:job_concurrent"`                         // 是否并发执行（0禁止 1允许）
	JobStatus        string `json:"jobStatus" gorm:"column:job_status"`                                 // 任务状态（0暂停 1正常）
	JobSaveLog       string `json:"jobSaveLog" gorm:"column:job_save_log"`                              // 是否记录任务日志（0不记录 1记录）
}

// TableName 表名称
func (*SysJob) TableName() string {
	return "sys_job"
}
