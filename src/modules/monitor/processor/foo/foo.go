package foo

import (
	"mask_api_gin/src/framework/cron"
	"mask_api_gin/src/framework/logger"

	"time"
)

var NewProcessor = &Processor{
	progress: 0,
	count:    0,
}

// Processor 队列任务处理
type Processor struct {
	progress int // 任务进度
	count    int // 执行次数
}

func (p *Processor) Execute(data any) (any, error) {
	logger.Infof("执行 %d 次，上次进度： %d ", p.count, p.progress)
	p.count++

	options := data.(cron.JobData)
	sysJob := options.SysJob
	logger.Infof("重复 %v 任务ID %s", options.Repeat, sysJob.JobId)

	// 实现任务处理逻辑
	i := 0
	p.progress = i
	for i < 20 {
		// 获取任务进度
		progress := p.progress
		logger.Infof("jonId: %s => 任务进度：%d", sysJob.JobId, progress)
		// 延迟响应
		time.Sleep(time.Second * 2)
		i++
		// 改变任务进度
		p.progress = i
	}

	// 返回结果，用于记录执行结果
	result := map[string]any{
		"repeat":       options.Repeat,
		"jobName":      sysJob.JobName,
		"invokeTarget": sysJob.InvokeTarget,
		"targetParams": sysJob.TargetParams,
	}
	return result, nil
}
