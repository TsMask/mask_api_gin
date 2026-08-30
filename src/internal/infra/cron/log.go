package cron

import (
	"fmt"
	"sync/atomic"
	"time"
)

// JobLogRecord 任务执行事件记录
//
// Data 为 RunJob 传入的原始负载，业务侧自行类型断言。
type JobLogRecord struct {
	Timestamp int64      // 执行开始时间（毫秒）
	Status    string     // constants.STATUS_YES 成功 / constants.STATUS_NO 失败
	Result    any        // 成功为处理器返回值，失败为错误信息
	Data      any        // RunJob 传入的原始 data
	Opts      JobOptions // 任务参数信息
}

// JobLogSink 任务执行日志持久化钩子
type JobLogSink func(JobLogRecord)

// onJobLog 使用原子指针保证运行期任务读取与启动期注册不会产生数据竞争。
var onJobLog atomic.Pointer[JobLogSink]

// EngineError 表示调度引擎自身的受理或回调异常，不表示处理器业务失败。
type EngineError struct {
	Timestamp int64  // 错误发生时间（毫秒）
	Queue     string // 队列名称
	ID        string // 任务编号
	Err       error  // 具体错误
}

// EngineErrorSink 引擎错误回调钩子，由业务层决定记录或上报方式。
type EngineErrorSink func(EngineError)

// onEngineError 使用原子指针保证运行期读取与启动期注册不会产生数据竞争。
var onEngineError atomic.Pointer[EngineErrorSink]

// SetOnJobLog 注入任务执行日志持久化钩子；nil 表示丢弃执行事件。
func SetOnJobLog(sink JobLogSink) {
	if sink == nil {
		onJobLog.Store(nil)
		return
	}
	onJobLog.Store(&sink)
}

// SetOnEngineError 注入引擎错误回调；nil 表示丢弃引擎错误。
func SetOnEngineError(sink EngineErrorSink) {
	if sink == nil {
		onEngineError.Store(nil)
		return
	}
	onEngineError.Store(&sink)
}

// emitEngineError 派发引擎错误到注入的 sink。
func emitEngineError(queueName, id string, err error) {
	sinkPtr := onEngineError.Load()
	if sinkPtr == nil {
		return
	}

	errorRecord := EngineError{
		Timestamp: time.Now().UnixMilli(),
		Queue:     queueName,
		ID:        id,
		Err:       err,
	}
	(*sinkPtr)(errorRecord)
}

// emitJobLog 派发执行事件到注入的 sink；钩子 panic 会转为引擎错误外抛。
func emitJobLog(queueName string, record JobLogRecord) {
	sinkPtr := onJobLog.Load()
	if sinkPtr == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			emitEngineError(queueName, record.Opts.ID, fmt.Errorf("cron job log sink panic: %v", r))
		}
	}()
	(*sinkPtr)(record)
}
