package testjob

import (
	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/modules/monitor/model"
)

var NewProcessor = &Processor{}

// Processor 队列任务处理
type Processor struct{}

func (p Processor) Execute(data any) (any, error) {
	sysJob := data.(model.SysJob)
	logger.Infof("任务ID %d", sysJob.ID)

	// 返回结果，用于记录执行结果
	result := map[string]any{
		"jobName":      sysJob.JobName,
		"invokeTarget": sysJob.JobInvokeTarget,
		"targetParams": sysJob.JobTargetParams,
	}
	return result, nil
}
