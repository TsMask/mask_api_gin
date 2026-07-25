package job

import (
	"encoding/json"
	"time"

	"mask_api_gin/src/internal/infra/cron"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/modules/job/testerror"
	"mask_api_gin/src/modules/job/testjob"
	"mask_api_gin/src/modules/job/testprogress"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/service"
)

// Setup 模块注册
func Setup() {
	// 初始化定时任务处理
	cron.CreateQueue("test_job", testjob.NewProcessor)
	cron.CreateQueue("test_progress", testprogress.NewProcessor)
	cron.CreateQueue("test_error", testerror.NewProcessor)
}

// RegisterOnJobLog 注入任务执行日志持久化钩子。
//
// 将引擎发出的执行事件映射为调度任务日志并落库；引擎本身不依赖任何业务模型。
// 需在调度任务初始化（Reset）之前调用。
func RegisterOnJobLog() {
	cron.OnJobLog = func(record cron.JobLogRecord) {
		// 还原业务负载
		sysJob, ok := record.Data.(model.SysJob)
		if !ok {
			return
		}

		// 任务配置为不记录日志则跳过
		if sysJob.JobSaveLog == "" || sysJob.JobSaveLog == constants.STATUS_NO {
			return
		}

		// 结果信息 key 的名称
		resultName := "failed"
		if record.Status == constants.STATUS_YES {
			resultName = "completed"
		}

		// 结果信息序列化字符串
		jsonByte, _ := json.Marshal(map[string]any{
			"cron":          sysJob.JobCron,
			"invoke_target": sysJob.JobInvokeTarget,
			"result":        resultName,
			"message":       record.Result,
		})
		jobMsg := string(jsonByte)
		if len(jobMsg) > 500 {
			jobMsg = jobMsg[:500]
		}

		// 计算执行耗时并写入日志
		duration := time.Since(time.UnixMilli(record.Timestamp))
		service.NewSysJobLog.Insert(model.SysJobLog{
			JobName:         sysJob.JobName,
			JobGroup:        sysJob.JobGroup,
			JobInvokeTarget: sysJob.JobInvokeTarget,
			JobTargetParams: sysJob.JobTargetParams,
			LogStatus:       record.Status,
			LogMsg:          jobMsg,
			LogCostTime:     duration.Milliseconds(),
			CreateTime:      time.Now().UnixMilli(),
		})
	}
}
