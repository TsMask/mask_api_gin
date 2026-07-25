package cron

import "github.com/robfig/cron/v3"

// 任务执行状态
const (
	// Waiting 已受理，等待执行
	Waiting = iota
	// Active 执行中
	Active
	// Completed 执行完成
	Completed
	// Failed 执行失败
	Failed
)

// JobOptions 任务参数信息
type JobOptions struct {
	ID   string // 执行任务编号
	Cron string // 重复任务 cron 表达式，空字符串表示立即执行一次

	// Concurrent 是否允许并发执行。
	//
	// false（默认）：同 ID 任务上次尚未执行完成时，本次触发将被跳过；
	// true：允许上一次未完成即再次执行。
	Concurrent bool
}

// queueJob 队列中的任务条目（以指针形式存入 map，使状态可写回共享实例）
type queueJob struct {
	id         string       // 任务编号（JobOptions.ID）
	data       any          // 执行任务时传入的参数
	cronExpr   string       // cron 表达式，空字符串表示一次性任务
	concurrent bool         // 是否允许并发执行
	cid        cron.EntryID // robfig/cron 排程 ID（一次性任务为 0）
	status     int          // 任务执行状态
	timestamp  int64        // 本次执行开始时间（毫秒）
}

// jobRunner 适配 robfig/cron 的 Job 接口，每次触发都从队列中按 ID 查找最新任务条目执行。
type jobRunner struct {
	queue *Queue
	id    string
}

// Run 实现 cron.Job
func (r *jobRunner) Run() {
	r.queue.runOnce(r.id)
}
