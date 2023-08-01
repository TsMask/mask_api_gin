package simple

import (
	"mask_api_gin/src/framework/cron"
	"mask_api_gin/src/framework/logger"
)

// Execute 队列任务处理
func Execute(options cron.Options) any {
	sysJob := options.SysJob

	logger.Infof("重复 %v 任务ID %s", options.Repeat, sysJob.JobID)

	// 返回结果，用于记录执行结果
	return map[string]any{
		"repeat":       options.Repeat,
		"jobName":      sysJob.JobName,
		"invokeTarget": sysJob.InvokeTarget,
		"targetParams": sysJob.TargetParams,
	}
}
