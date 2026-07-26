package model

// SysJobLog 调度任务调度日志表（按 job_name+job_group 归属，冗余调用信息便于历史追溯）
type SysJobLog struct {
	ID              int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`                         // 任务日志ID
	CreateTime      int64  `json:"create_time" gorm:"column:create_time"`                                // 创建时间
	JobName         string `json:"job_name" gorm:"column:job_name" binding:"required"`                   // 任务名称
	JobGroup        string `json:"job_group" gorm:"column:job_group" binding:"required"`                 // 任务组名
	JobInvokeTarget string `json:"job_invoke_target" gorm:"column:job_invoke_target" binding:"required"` // 调用目标字符串
	JobTargetParams string `json:"job_target_params" gorm:"column:job_target_params"`                    // 调用目标传入参数
	LogMsg          string `json:"log_msg" gorm:"column:log_msg"`                                        // 日志信息
	LogStatus       string `json:"log_status" gorm:"column:log_status"`                                  // 执行状态（0失败 1正常）
	LogCostTime     int64  `json:"log_cost_time" gorm:"column:log_cost_time"`                            // 消耗时间（毫秒）
}

// TableName 表名称
func (*SysJobLog) TableName() string {
	return "sys_job_log"
}
