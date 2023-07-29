package cron

// https://blog.csdn.net/zjbyough/article/details/113853582
// https://mp.weixin.qq.com/s/Ak7RBv1NuS-VBeDNo8_fww
// 还有优化空间 TODO
import (
	"errors"
	"mask_api_gin/src/modules/monitor/model"
	"sync"

	"github.com/robfig/cron/v3"
)

// Options 调度任务处理接收参数
type Options struct {
	// 触发执行cron重复多次
	Repeat bool
	// 定时任务调度表记录信息
	SysJob model.SysJob
}

// 定义内部调度任务实例
var c *cron.Cron

// 互斥锁
var mutex sync.Mutex

// 日志收集
var clog = cronLog{}

// InitNew 初始调度任务实例
func InitNew() {
	c = cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.Recover(clog), // 添加默认的 panic 恢复逻辑
			cron.DelayIfStillRunning(clog),
			cron.SkipIfStillRunning(clog),
		),
	)

	// 实例启动
	c.Start()
}

// 队列
type queue struct {
	// 队列名称
	Name string
	// 任务函数
	JobFunc jobFunc
	// 处理器任务列表
	Processor *[]job
}

// 任务函数
type jobFunc func(options Options) interface{}

// 处理器任务
type job struct {
	// 任务ID
	JobID string
	// 执行任务ID
	EntryID cron.EntryID
}

// 定义内部队列列表
var queueList []queue = make([]queue, 0)

// QueueList 队列列表
func QueueList() []queue {
	return queueList
}

// AddQueue 添加队列
func AddQueue(name string, jobFun jobFunc) {
	// 检查是否已有队列名称
	for _, queue := range queueList {
		if queue.Name == name {
			return
		}
	}
	// 没有则正常添加到队列
	queueList = append(queueList, queue{
		Name:      name,
		JobFunc:   jobFun,
		Processor: &[]job{},
	})
}

// GetQueue 获取队列
func GetQueue(name string) (queue, error) {
	for _, queue := range queueList {
		if queue.Name == name {
			return queue, nil
		}
	}
	return queue{}, errors.New("未找到队列")
}

// GetJob 队列上的任务
func (q queue) GetJob(jobId string) job {
	for _, v := range *q.Processor {
		if v.JobID == jobId {
			return v
		}
	}
	return job{}
}

// RunJob 队列上运行任务
func (q queue) RunJob(options Options, jobId string, spec string) {
	cmd := func() {
		mutex.Lock()         // 加锁
		defer mutex.Unlock() // 解锁
		// panics 异常收集
		defer func() {
			if r := recover(); r != nil {
				clog.Error(nil, "failed", options, r)
			}
		}()
		// 执行任务函数
		result := q.JobFunc(options)
		// 记录完成结果
		clog.Completed(options, result)
	}

	// 不是重复任务
	if !options.Repeat {
		entryID, err := c.AddFunc(spec, cmd)
		if err != nil {
			clog.Error(err, "添加任务失败！")
		}
		ce := c.Entry(entryID)
		if ce.Valid() {
			go ce.Job.Run()
		}
		c.Remove(entryID)
		return
	}

	// 检查是否有任务实例
	var ce cron.Entry
	for _, v := range *q.Processor {
		if v.JobID == jobId {
			ce = c.Entry(v.EntryID)
			break
		}
	}

	// 移除存在的任务实例
	if ce.Valid() {
		// 移除任务实例
		c.Remove(ce.ID)
		// 过滤去掉原先的实例ID
		filtered := make([]job, 0)
		for _, v := range *q.Processor {
			if v.EntryID != ce.ID {
				filtered = append(filtered, v)
			}
		}
		*q.Processor = filtered
	}

	// 添加任务到 cron 中，并设置执行时间规则
	entryID, err := c.AddFunc(spec, cmd)
	if err != nil {
		clog.Error(err, "添加任务失败！")
	}

	// 新增新任务实例ID
	*q.Processor = append(*q.Processor, job{
		JobID:   jobId,
		EntryID: entryID,
	})
}

// RemoveJob 队列上移除任务
func (q queue) RemoveJob(jobId string) bool {
	// 检查是否有任务实例
	var ce cron.Entry
	for _, v := range *q.Processor {
		if v.JobID == jobId {
			ce = c.Entry(v.EntryID)
			break
		}
	}

	// 移除存在的任务实例
	if ce.Valid() {
		// 移除任务实例
		c.Remove(ce.ID)
		// 过滤去掉原先的实例ID
		filtered := make([]job, 0)
		for _, v := range *q.Processor {
			if v.EntryID != ce.ID {
				filtered = append(filtered, v)
			}
		}
		*q.Processor = filtered
		return true
	}
	return false
}
