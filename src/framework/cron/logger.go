package cron

import (
	"encoding/json"
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/logger"
	"mask_api_gin/src/modules/monitor/model"
	"mask_api_gin/src/modules/monitor/repository"
)

// cronLog 日志收集
type cronLog struct{}

// Info 任务普通信息收集
func (clog cronLog) Info(msg string, keysAndValues ...interface{}) {
	logger.Infof("Info ====> %s kv: %v", msg, keysAndValues)
}

// Error 任务异常错误收集
func (clog cronLog) Error(err error, msg string, keysAndValues ...interface{}) {
	// 不是指定的错误收集
	if msg != "failed" {
		logger.Infof("Error %v ====> %s kv: %v", err, msg, keysAndValues)
		return
	}
	// 结果信息序列化字符串
	msgMap := map[string]interface{}{
		"name":    "failed",
		"message": keysAndValues[1],
	}
	jsonByte, _ := json.Marshal(msgMap)
	jobMsg := string(jsonByte)
	if len(jobMsg) > 500 {
		jobMsg = jobMsg[:500]
	}
	// 读取任务信息创建日志对象
	if options, ok := keysAndValues[0].(Options); ok {
		sysJob := options.SysJob
		sysJobLog := model.SysJobLog{
			JobName:      sysJob.JobName,
			JobGroup:     sysJob.JobGroup,
			InvokeTarget: sysJob.InvokeTarget,
			TargetParams: sysJob.TargetParams,
			Status:       common.STATUS_YES,
			JobMsg:       jobMsg,
		}
		// 插入数据
		repository.NewSysJobLogImpl.InsertJobLog(sysJobLog)
	}
}

// Completed 任务完成return的结果收集
func (clog cronLog) Completed(options Options, result interface{}) {
	// 结果信息序列化字符串
	msgMap := map[string]interface{}{
		"name":    "completed",
		"message": result,
	}
	jsonByte, _ := json.Marshal(msgMap)
	jobMsg := string(jsonByte)
	if len(jobMsg) > 500 {
		jobMsg = jobMsg[:500]
	}
	// 读取任务信息创建日志对象
	sysJob := options.SysJob
	sysJobLog := model.SysJobLog{
		JobName:      sysJob.JobName,
		JobGroup:     sysJob.JobGroup,
		InvokeTarget: sysJob.InvokeTarget,
		TargetParams: sysJob.TargetParams,
		Status:       common.STATUS_YES,
		JobMsg:       jobMsg,
	}
	// 插入数据
	repository.NewSysJobLogImpl.InsertJobLog(sysJobLog)
}
