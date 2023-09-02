package bar

import (
	"mask_api_gin/src/framework/cron"
	"mask_api_gin/src/framework/logger"
	"time"
)

var NewProcessor = &BarProcessor{
	progress: 0,
	count:    0,
}

// bar 队列任务处理
type BarProcessor struct {
	// 任务进度
	progress int
	// 执行次数
	count int
}

func (s *BarProcessor) Execute(data any) any {
	logger.Infof("执行 %d 次，上次进度： %d ", s.count, s.progress)
	s.count++

	options := data.(cron.JobData)
	sysJob := options.SysJob
	logger.Infof("重复 %v 任务ID %s", options.Repeat, sysJob.JobID)

	// 实现任务处理逻辑
	i := 0
	s.progress = i
	for i < 5 {
		// 获取任务进度
		progress := s.progress
		logger.Infof("jonId: %s => 任务进度：%d", sysJob.JobID, progress)
		// 延迟响应
		time.Sleep(time.Second * 2)
		// 程序中途执行错误
		if i == 3 {
			// arr := [1]int{1}
			// arr[i] = 3
			// fmt.Println(arr)
			// return "i = 3"
			panic("程序中途执行错误")
		}
		i++
		// 改变任务进度
		s.progress = i
	}

	// 返回结果，用于记录执行结果
	return map[string]any{
		"repeat":       options.Repeat,
		"jobName":      sysJob.JobName,
		"invokeTarget": sysJob.InvokeTarget,
		"targetParams": sysJob.TargetParams,
	}
}
