package cron

// JobLogRecord 任务执行事件记录
//
// 引擎在任务成功或失败时发出，由业务模块注入的 OnJobLog 持久化。
// Data 为 RunJob 传入的原始负载，业务侧自行类型断言（如 monitor 的 JobData）。
type JobLogRecord struct {
	Timestamp int64      // 执行开始时间（毫秒）
	Status    string     // constants.STATUS_YES 成功 / constants.STATUS_NO 失败
	Result    any        // 成功为处理器返回值，失败为错误信息
	Data      any        // RunJob 传入的原始 data
	Opts      JobOptions // 任务参数信息
}

// OnJobLog 任务执行日志持久化钩子
//
// 由业务模块在启动时注入；为 nil 时丢弃执行事件，引擎不依赖任何业务模型。
var OnJobLog func(JobLogRecord)

// emitJobLog 派发执行事件到注入的 sink（未注入时丢弃）
func emitJobLog(record JobLogRecord) {
	if OnJobLog != nil {
		OnJobLog(record)
	}
}
