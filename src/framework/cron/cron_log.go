package cron

import (
	"encoding/json"
	constCommon "mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
	"time"
)

// cronLog 实例任务执行日志收集
var cronLog = clog{}

// clog 任务执行日志收集
type clog struct{}

// Info 任务普通信息收集
func (clog) Info(msg string, keysAndValues ...any) {
	//log.Printf("Info msg: %v ====> kv: %v", msg, keysAndValues)
}

// Error 任务异常错误收集
func (clog) Error(err error, msg string, keysAndValues ...any) {
	//log.Printf("Error: %v -> msg: %v ====> kv: %v", err, msg, keysAndValues)
	//log.Printf("k0: %v", keysAndValues[0].(QueueJob))

	// 指定的错误收集
	if msg == "failed" {
		// 任务对象
		job := keysAndValues[0].(QueueJob)

		// 读取任务信息进行保存日志
		if data, ok := job.Data.(JobData); ok {
			// 日志数据
			jobLog := jobLogData{
				Timestamp: job.Timestamp,
				Data:      data,
				Result:    err.Error(),
			}
			jobLog.SaveLog(constCommon.StatusNo)
		}
	}
}

// Completed 任务完成return的结果收集
func (clog) Completed(result any, msg string, keysAndValues ...any) {
	//log.Printf("Completed: %v -> msg: %v ====> kv: %v", result, msg, keysAndValues)
	//log.Printf("k0: %v", keysAndValues[0].(QueueJob))

	// 指定的完成收集
	if msg == "completed" {
		// 任务对象
		job := keysAndValues[0].(QueueJob)

		// 读取任务信息进行保存日志
		if data, ok := job.Data.(JobData); ok {
			// 日志数据
			jobLog := jobLogData{
				Timestamp: job.Timestamp,
				Data:      data,
				Result:    result,
			}
			jobLog.SaveLog(constCommon.StatusYes)
		}
	}
}

// jobLogData 日志记录数据
type jobLogData struct {
	Timestamp int64
	Data      JobData
	Result    any
}

// SaveLog 日志记录保存
func (jl *jobLogData) SaveLog(status string) {
	// 读取任务信息
	sysJob := jl.Data.SysJob

	// 任务日志不需要记录
	if sysJob.SaveLog == "" || sysJob.SaveLog == constCommon.StatusNo {
		return
	}

	// 结果信息key的Name
	resultName := "failed"
	if status == constCommon.StatusYes {
		resultName = "completed"
	}

	// 结果信息序列化字符串
	jsonByte, _ := json.Marshal(map[string]any{
		"cron":    jl.Data.Repeat,
		"name":    resultName,
		"message": jl.Result,
	})
	jobMsg := string(jsonByte)
	if len(jobMsg) > 500 {
		jobMsg = jobMsg[:500]
	}

	// 创建日志对象
	duration := time.Since(time.UnixMilli(jl.Timestamp))
	sysJobLog := model.SysJobLog{
		JobName:      sysJob.JobName,
		JobGroup:     sysJob.JobGroup,
		InvokeTarget: sysJob.InvokeTarget,
		TargetParams: sysJob.TargetParams,
		Status:       status,
		JobMsg:       jobMsg,
		CostTime:     duration.Milliseconds(),
	}
	// 插入数据
	repository.NewSysJobLogRepository.Insert(sysJobLog)
}

// JobData 调度任务日志收集结构体，执行任务时传入的接收参数
type JobData struct {
	// 触发执行cron重复多次
	Repeat bool
	// 定时任务调度表记录信息
	SysJob model.SysJob
}
