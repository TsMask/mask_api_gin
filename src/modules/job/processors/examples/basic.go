package examples

import (
	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/modules/job/model"
)

// NewBasicProcessor demo: returns job info immediately, used to verify the basic execution flow.
var NewBasicProcessor = &BasicProcessor{}

// BasicProcessor basic sample job processor.
type BasicProcessor struct{}

func (p *BasicProcessor) Execute(data any) (any, error) {
	sysJob := data.(model.SysJob)
	logger.Infof("任务ID %d", sysJob.ID)

	// result payload, recorded into job log
	result := map[string]any{
		"job_name":      sysJob.JobName,
		"invoke_target": sysJob.JobInvokeTarget,
		"target_params": sysJob.JobTargetParams,
	}
	return result, nil
}
