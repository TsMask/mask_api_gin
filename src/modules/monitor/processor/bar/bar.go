package bar

import (
	"mask_api_gin/src/framework/cron"
	"mask_api_gin/src/framework/logger"
	"time"
)

// 任务进度
var progress = 0

// Execute 队列任务处理
func Execute(options cron.Options) interface{} {
	defer func() {
		progress = 0
	}()

	sysJob := options.SysJob

	logger.Infof("重复 %v 任务ID %s", options.Repeat, sysJob.JobID)

	for progress < 10 {
		// 获取任务进度
		logger.Infof("jonId: %s => 任务进度：%d", sysJob.JobID, progress)
		// 延迟响应
		time.Sleep(2 * time.Second)
		// 程序中途执行错误
		if progress > 3 {
			panic("程序中途执行错误")
		}
		// 改变任务进度
		progress++
	}

	// 返回结果，用于记录执行结果
	return map[string]interface{}{
		"repeat":       options.Repeat,
		"jobName":      sysJob.JobName,
		"invokeTarget": sysJob.InvokeTarget,
		"targetParams": sysJob.TargetParams,
	}
}