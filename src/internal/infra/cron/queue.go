package cron

import (
	"fmt"
	"sync"
	"time"

	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/logger"
)

// QueueProcessor 队列处理函数接口
type QueueProcessor interface {
	// Execute 实际执行函数
	//   - any 返回有效值，记录为成功
	//   - error 非 nil，记录为失败
	Execute(data any) (any, error)
}

// Queue 任务队列
type Queue struct {
	Name      string // 队列名
	processor QueueProcessor

	mu   sync.Mutex
	jobs map[string]*queueJob
}

// RunJob 运行任务，data 为传入给处理器的数据。
//
// 返回任务状态：Waiting 表示已受理（立即执行或已排程），
// Failed 表示排程失败（引擎未启动或 cron 表达式非法）。
func (q *Queue) RunJob(data any, opts JobOptions) int {
	q.mu.Lock()
	cr := currentCron()

	// 同 ID 任务先清理：撤销历史 cron 排程并移除条目
	if prev, ok := q.jobs[opts.ID]; ok {
		if prev.cid != 0 && cr != nil {
			cr.Remove(prev.cid)
		}
		delete(q.jobs, opts.ID)
	}

	job := &queueJob{
		id:         opts.ID,
		data:       data,
		cronExpr:   opts.Cron,
		concurrent: opts.Concurrent,
		status:     Waiting,
	}

	// 一次性任务：入队后立即异步执行（不依赖引擎是否启动）
	if opts.Cron == "" {
		q.jobs[opts.ID] = job
		q.mu.Unlock()
		go q.runOnce(opts.ID)
		return Waiting
	}

	// 重复任务：交由 robfig/cron 排程
	if cr == nil {
		logger.Errorf("cron 引擎未启动，无法排程重复任务 id=%s expr=%s", opts.ID, opts.Cron)
		job.status = Failed
		q.jobs[opts.ID] = job
		q.mu.Unlock()
		return Failed
	}
	cid, err := cr.AddJob(opts.Cron, &jobRunner{queue: q, id: opts.ID})
	if err != nil {
		logger.Errorf("cron AddJob 失败 id=%s expr=%s err=%v", opts.ID, opts.Cron, err)
		job.status = Failed
		q.jobs[opts.ID] = job
		q.mu.Unlock()
		return Failed
	}
	job.cid = cid
	q.jobs[opts.ID] = job
	q.mu.Unlock()
	return Waiting
}

// RemoveJob 按 ID 移除任务并撤销其 cron 排程
func (q *Queue) RemoveJob(id string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	cr := currentCron()
	prev, ok := q.jobs[id]
	if !ok {
		return false
	}
	if prev.cid != 0 && cr != nil {
		cr.Remove(prev.cid)
	}
	delete(q.jobs, id)
	return true
}

// runOnce 执行一次任务。
//
// 仅在 map 访问与状态流转时持锁；processor.Execute 在锁外执行，
// 避免阻塞同队列其他任务与 cron 触发。
func (q *Queue) runOnce(id string) {
	q.mu.Lock()
	job, ok := q.jobs[id]
	if !ok {
		q.mu.Unlock()
		return
	}
	// 非并发模式：上一次仍在执行则跳过本次触发
	if job.status == Active && !job.concurrent {
		q.mu.Unlock()
		return
	}
	job.status = Active
	job.timestamp = time.Now().UnixMilli()
	data := job.data
	timestamp := job.timestamp
	opts := JobOptions{ID: job.id, Cron: job.cronExpr, Concurrent: job.concurrent}
	q.mu.Unlock()

	// 锁外执行实际任务（含 panic 恢复）
	result, err := q.execute(data)

	status := constants.STATUS_YES
	if err != nil {
		status = constants.STATUS_NO
	}

	// 写回结果状态，并清理已完成的一次性任务
	q.mu.Lock()
	if cur, ok := q.jobs[id]; ok && cur == job {
		if err != nil {
			job.status = Failed
		} else {
			job.status = Completed
		}
		if job.cid == 0 {
			delete(q.jobs, id)
		}
	}
	q.mu.Unlock()

	// 锁外派发执行事件（未注入 sink 时丢弃）
	record := JobLogRecord{
		Timestamp: timestamp,
		Status:    status,
		Data:      data,
		Opts:      opts,
	}
	if status == constants.STATUS_NO {
		record.Result = err.Error()
	} else {
		record.Result = result
	}
	emitJobLog(record)
}

// execute 执行任务，将 panic 转为 error 返回
func (q *Queue) execute(data any) (result any, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = fmt.Errorf("%v", r)
		}
	}()
	return q.processor.Execute(data)
}
