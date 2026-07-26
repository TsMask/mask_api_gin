package examples

import (
	"time"

	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/modules/job/model"
)

// NewErrorProcessor demo: panics mid-run, used to verify panic recovery and error logging.
var NewErrorProcessor = &ErrorProcessor{}

// ErrorProcessor error sample job processor.
type ErrorProcessor struct {
	progress int // task progress
	count    int // run count
}

func (p *ErrorProcessor) Execute(data any) (any, error) {
	logger.Infof("执行 %d 次，上次进度：%d", p.count, p.progress)
	p.count++

	sysJob := data.(model.SysJob)
	logger.Infof("任务ID %d", sysJob.ID)

	// simulate work that fails halfway
	i := 0
	p.progress = i
	for i < 5 {
		progress := p.progress
		logger.Infof("jonId: %d => 任务进度：%d", sysJob.ID, progress)
		time.Sleep(time.Second * 2)
		// panic mid-run to exercise recovery
		if i == 3 {
			panic("程序中途执行错误")
		}
		i++
		p.progress = i
	}

	result := map[string]any{
		"job_name":      sysJob.JobName,
		"invoke_target": sysJob.JobInvokeTarget,
		"target_params": sysJob.JobTargetParams,
	}
	return result, nil
}
