package cron

import (
	"sync"

	"github.com/robfig/cron/v3"
)

// 调度引擎实例（nil 表示已停止）
var c *cron.Cron

// engineMu 保护 c 与 queueMap 的并发访问
var engineMu sync.RWMutex

// queueMap 已注册的任务队列（按名称索引）。
// 进程内常驻，不随引擎启停重置——处理器只在启动时注册一次。
var queueMap = map[string]*Queue{}

// StartCron 启动（或重启）调度引擎。
//
// 仅管理 *cron.Cron 实例，可在运行期反复调用；已注册的队列与处理器不受影响。
func StartCron() {
	engineMu.Lock()
	defer engineMu.Unlock()
	if c != nil {
		c.Stop()
	}
	c = cron.New(cron.WithSeconds())
	c.Start()
}

// StopCron 停止调度引擎：暂停所有 cron 排程，保留已注册的队列与处理器。
//
// 停止后 queueMap 仍在；再次 StartCron 即可恢复，由调用方按需重新排程任务。
func StopCron() {
	engineMu.Lock()
	defer engineMu.Unlock()
	if c != nil {
		c.Stop()
		c = nil
	}
}

// currentCron 在读锁下返回当前引擎实例（可能为 nil）。
//
// 返回的是快照指针：即便随后 StopCron 将 c 置 nil，已取出的实例仍可安全调用
// （robfig/cron 对已停止实例的 AddJob/Remove 为安全的空操作）。
func currentCron() *cron.Cron {
	engineMu.RLock()
	defer engineMu.RUnlock()
	return c
}

// CreateQueue 创建队列并注册处理器
func CreateQueue(name string, processor QueueProcessor) *Queue {
	engineMu.Lock()
	defer engineMu.Unlock()
	queue := &Queue{
		Name:      name,
		processor: processor,
		jobs:      map[string]*queueJob{},
	}
	queueMap[name] = queue
	return queue
}

// GetQueue 通过名称获取队列实例
func GetQueue(name string) *Queue {
	engineMu.RLock()
	defer engineMu.RUnlock()
	return queueMap[name]
}

// QueueNames 获取已注册的队列名称
func QueueNames() []string {
	engineMu.RLock()
	defer engineMu.RUnlock()
	keys := make([]string, 0, len(queueMap))
	for k := range queueMap {
		keys = append(keys, k)
	}
	return keys
}
