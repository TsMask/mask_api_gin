package cron

import (
	"encoding/json"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
	"time"
)

// 实例任务执行日志收集
var newLog = cronlog{}

// cronlog 任务执行日志收集
type cronlog struct{}

// Info 任务普通信息收集
func (s cronlog) Info(msg string, keysAndValues ...any) {
	// logger.Infof("Info msg: %v ====> kv: %v", msg, keysAndValues)

}

// Error 任务异常错误收集
func (s cronlog) Error(err error, msg string, keysAndValues ...any) {
	// logger.Errorf("Error: %v -> msg: %v ====> kv: %v", err, msg, keysAndValues)
	// logger.Errorf("k0: %v", keysAndValues[0].(*QueueJob))

	// 指定的错误收集
	if msg == "failed" {
		// 任务对象
		job := keysAndValues[0].(*QueueJob)

		// 结果信息序列化字符串
		jsonByte, _ := json.Marshal(map[string]any{
			"name":    "failed",
			"message": err.Error(),
		})
		jobMsg := string(jsonByte)
		if len(jobMsg) > 500 {
			jobMsg = jobMsg[:500]
		}

		// 读取任务信息创建日志对象
		if data, ok := job.Data.(JobData); ok {
			duration := time.Since(time.UnixMilli(job.Timestamp))
			sysJob := data.SysJob
			if sysJob.JobID == job.Opts.JobId {
				sysJobLog := model.SysJobLog{
					JobName:      sysJob.JobName,
					JobGroup:     sysJob.JobGroup,
					InvokeTarget: sysJob.InvokeTarget,
					TargetParams: sysJob.TargetParams,
					Status:       common.STATUS_NO,
					JobMsg:       jobMsg,
					CostTime:     duration.Milliseconds(),
				}
				// 插入数据
				repository.NewSysJobLogImpl.InsertJobLog(sysJobLog)
			}
		}
	}
}

// Completed 任务完成return的结果收集
func (s cronlog) Completed(result any, msg string, keysAndValues ...any) {
	// logger.Infof("Completed: %v -> msg: %v ====> kv: %v", result, msg, keysAndValues)
	// logger.Infof("k0: %v", keysAndValues[0].(*QueueJob))

	// 指定的完成收集
	if msg == "completed" {
		// 任务对象
		job := keysAndValues[0].(*QueueJob)

		// 结果信息序列化字符串
		jsonByte, _ := json.Marshal(map[string]any{
			"name":    "completed",
			"message": result,
		})
		jobMsg := string(jsonByte)
		if len(jobMsg) > 500 {
			jobMsg = jobMsg[:500]
		}

		// 读取任务信息创建日志对象
		if data, ok := job.Data.(JobData); ok {
			duration := time.Since(time.UnixMilli(job.Timestamp))
			sysJob := data.SysJob
			if sysJob.JobID == job.Opts.JobId {
				sysJobLog := model.SysJobLog{
					JobName:      sysJob.JobName,
					JobGroup:     sysJob.JobGroup,
					InvokeTarget: sysJob.InvokeTarget,
					TargetParams: sysJob.TargetParams,
					Status:       common.STATUS_YES,
					JobMsg:       jobMsg,
					CostTime:     duration.Milliseconds(),
				}
				// 插入数据
				repository.NewSysJobLogImpl.InsertJobLog(sysJobLog)
			}
		}
	}
}

// JobData 调度任务日志收集结构体，执行任务时传入的接收参数
type JobData struct {
	// 触发执行cron重复多次
	Repeat bool
	// 定时任务调度表记录信息
	SysJob model.SysJob
}
