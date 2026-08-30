package cron

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"mask_api_gin/src/internal/shared/constants"
)

// ErrCronStopped 表示调度引擎已停止，无法注册重复任务。
var ErrCronStopped = errors.New("cron engine is stopped")

// ErrExecutorFull 表示进程内并发执行槽已满，任务未被受理。
var ErrExecutorFull = errors.New("cron executor is full")

// QueueProcessor 队列处理函数接口
type QueueProcessor interface {
	// Execute 实际执行函数
	//   - any 返回有效值时，记录为成功
	//   - error 非 nil，记录为失败
	Execute(data any) (any, error)
}

// Queue 任务队列
type Queue struct {
	Name      string // 队列名称
	processor QueueProcessor

	mu      sync.Mutex
	jobs    map[string]*queueJob
	running map[string]int
}

// RunJob 运行任务，data 为传入给处理器的数据。
//
// 返回任务状态：Waiting 表示已受理，Active 表示同 ID 非并发任务正在执行，
// Failed 表示受理失败（引擎停机、表达式非法或执行容量已满）。
func (q *Queue) RunJob(data any, opts JobOptions) int {
	q.mu.Lock()

	// 一次性任务以运行计数判断重叠，不能依赖 jobs 中的旧条目是否已被替换。
	if opts.Cron == "" && q.running[opts.ID] > 0 && !opts.Concurrent {
		q.mu.Unlock()
		return Active
	}

	// 同 ID 任务先清理：撤销历史 cron 排程并移除条目。
	if prev, ok := q.jobs[opts.ID]; ok {
		removeCronJob(prev.engine, prev.cid)
		delete(q.jobs, opts.ID)
	}

	job := &queueJob{
		id:         opts.ID,
		data:       data,
		cronExpr:   opts.Cron,
		concurrent: opts.Concurrent,
		status:     Waiting,
	}

	// 一次性任务入队后立即异步执行，不依赖调度引擎是否启动。
	if opts.Cron == "" {
		q.jobs[opts.ID] = job
		q.mu.Unlock()
		dispatchErr := dispatch(nil, q, opts.ID)
		if dispatchErr == nil {
			return Waiting
		}
		q.mu.Lock()
		if cur, ok := q.jobs[opts.ID]; ok && cur == job {
			delete(q.jobs, opts.ID)
		}
		q.mu.Unlock()
		emitEngineError(q.Name, opts.ID, dispatchErr)
		return Failed
	}

	entryID, engine, err := scheduleCronJob(opts.Cron, q, opts.ID)
	if err != nil {
		job.status = Failed
		q.jobs[opts.ID] = job
		q.mu.Unlock()
		emitEngineError(q.Name, opts.ID, err)
		return Failed
	}
	job.cid = entryID
	job.engine = engine
	q.jobs[opts.ID] = job
	q.mu.Unlock()
	return Waiting
}

// RemoveJob 按 ID 移除任务并撤销其 cron 排程
func (q *Queue) RemoveJob(id string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	prev, ok := q.jobs[id]
	if !ok {
		return false
	}
	removeCronJob(prev.engine, prev.cid)
	delete(q.jobs, id)
	return true
}

// runOnce 执行一次任务。
//
// 仅在 map 访问、运行计数与状态流转时持锁；processor.Execute 在锁外执行。
func (q *Queue) runOnce(id string) {
	q.mu.Lock()
	job, ok := q.jobs[id]
	if !ok {
		q.mu.Unlock()
		return
	}
	// 非并发模式以实际运行计数为准，覆盖任务条目在执行期间被替换的情况。
	if !job.concurrent && q.running[id] > 0 {
		q.mu.Unlock()
		return
	}
	q.running[id]++
	releaseRunning := func() {
		q.mu.Lock()
		if q.running[id] <= 1 {
			delete(q.running, id)
		} else {
			q.running[id]--
		}
		q.mu.Unlock()
	}

	job.status = Active
	job.timestamp = time.Now().UnixMilli()
	data := job.data
	timestamp := job.timestamp
	opts := JobOptions{ID: job.id, Cron: job.cronExpr, Concurrent: job.concurrent}
	q.mu.Unlock()
	defer releaseRunning()

	result, err := q.execute(data)

	status := constants.STATUS_YES
	if err != nil {
		status = constants.STATUS_NO
	}

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

	record := JobLogRecord{
		Timestamp: timestamp,
		Status:    status,
		Data:      data,
		Opts:      opts,
	}
	if err != nil {
		record.Result = err.Error()
	} else {
		record.Result = result
	}
	emitJobLog(q.Name, record)
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
