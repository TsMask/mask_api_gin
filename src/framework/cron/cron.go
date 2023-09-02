package cron

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// 定义内部调度任务实例
var c *cron.Cron

// 任务队列
var queueMap map[string]Queue

// StartCron 启动调度任务实例
func StartCron() {
	queueMap = make(map[string]Queue)
	c = cron.New(cron.WithSeconds())
	c.Start()
}

// StopCron 停止调度任务实例
func StopCron() {
	c.Stop()
}

// CreateQueue 创建队列注册处理器
func CreateQueue(name string, processor QueueProcessor) Queue {
	queue := Queue{
		Name:      name,
		Processor: processor,
		Job:       &[]*QueueJob{},
	}
	queueMap[name] = queue
	return queue
}

// GetQueue 通过名称获取队列实例
func GetQueue(name string) Queue {
	if v, ok := queueMap[name]; ok {
		return v
	}
	return Queue{}
}

// QueueNames 获取注册的队列名称
func QueueNames() []string {
	keys := make([]string, 0, len(queueMap))
	for k := range queueMap {
		keys = append(keys, k)
	}
	return keys
}

// Queue 任务队列
type Queue struct {
	Name      string // 队列名
	Processor QueueProcessor
	Job       *[]*QueueJob
}

// QueueProcessor 队列处理函数接口
type QueueProcessor interface {
	// Execute 实际执行函数
	Execute(data any) any
}

// RunJob 运行任务，data是传入的数据
func (q *Queue) RunJob(data any, opts JobOptions) int {
	job := &QueueJob{
		Status:         Waiting,
		Data:           data,
		Opts:           opts,
		queueName:      q.Name,
		queueProcessor: &q.Processor,
	}

	// 非重复任务立即执行
	if opts.Cron == "" {
		one := job.runJob(false)
		newLog.Info("RunJob", one.cid, opts.JobId, one.Status)
		if one.Status == Waiting || one.Status == Completed {
			go one.Run()
		}
		return one.Status
	}

	// 移除已存的任务ID
	q.RemoveJob(opts.JobId)

	// 添加新任务
	cid, err := c.AddJob(opts.Cron, job)
	if err != nil {
		newLog.Error(err, "err")
		job.Status = Failed
	}
	job.cid = cid
	*q.Job = append(*q.Job, job)
	newLog.Info("RunJob", cid, opts.JobId)
	return job.Status
}

// RemoveJob 移除任务
func (q *Queue) RemoveJob(jobId string) bool {
	// 移除已存的任务ID
	for i, v := range *q.Job {
		if jobId == v.Opts.JobId {
			newLog.Info("RemoveJob", v.cid, jobId)
			c.Remove(v.cid)
			// 从切片 jobs 中删除指定索引位置的元素
			jobs := *q.Job
			jobs = append(jobs[:i], jobs[i+1:]...)
			*q.Job = jobs
			return true
		}
	}
	return false
}

// Status 任务执行状态
const (
	Waiting = iota
	Active
	Completed
	Failed
)

// JobOptions 任务参数信息
type JobOptions struct {
	JobId string // 执行任务编号
	Cron  string // 重复任务cron表达式
}

// QueueJob 队列内部执行任务
type QueueJob struct {
	Status    int   // 任务执行状态
	Timestamp int64 // 执行时间
	Data      any   // 执行任务时传入的参数
	Opts      JobOptions

	cid cron.EntryID // 执行ID

	queueName      string //队列名
	queueProcessor *QueueProcessor
}

// runJob 通过队列名获取当前执行ID任务
func (job *QueueJob) runJob(repeat bool) *QueueJob {
	q := GetQueue(job.queueName)
	for _, v := range *q.Job {
		if repeat {
			if v.Opts.JobId == job.Opts.JobId {
				return v
			}
		} else {
			if v.Opts.JobId == job.Opts.JobId && v.cid == 0 {
				return v
			}
		}
	}
	*q.Job = append(*q.Job, job)
	return job
}

// Run 实现的接口函数
func (s QueueJob) Run() {
	// 检查当前任务
	job := s.runJob(s.cid != 0)

	// Active 状态不执行
	if Active == job.Status {
		return
	}

	// panics 异常收集
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			job.Status = Failed
			newLog.Error(err, "failed", job)
		}
	}()

	// 开始执行
	job.Status = Active
	job.Timestamp = time.Now().UnixMilli()
	newLog.Info("run", job.cid, job.Opts.JobId)

	// 获取队列处理器接口实现
	processor := *job.queueProcessor
	result := processor.Execute(job.Data)
	job.Status = Completed
	newLog.Completed(result, "completed", job)
}
