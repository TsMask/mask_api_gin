package examples

import (
	"time"

	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/modules/job/model"
)

// NewProgressProcessor demo: loops and reports progress, used to verify progress tracking.
var NewProgressProcessor = &ProgressProcessor{}

// ProgressProcessor progress sample job processor.
type ProgressProcessor struct {
	progress int // task progress
	count    int // run count
}

func (p *ProgressProcessor) Execute(data any) (any, error) {
	logger.Infof("执行 %d 次，上次进度：%d", p.count, p.progress)
	p.count++

	sysJob := data.(model.SysJob)
	logger.Infof("任务ID %d", sysJob.ID)

	// simulate work that reports progress
	i := 0
	p.progress = i
	for i < 20 {
		progress := p.progress
		logger.Infof("jonId: %d => 任务进度：%d", sysJob.ID, progress)
		time.Sleep(time.Second * 2)
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
