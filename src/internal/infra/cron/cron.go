package cron

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"
)

// maxConcurrentJobs 限制整个进程内同时执行的任务数量。
const maxConcurrentJobs = 64

// engineState 表示一代调度引擎；重建后旧实例不共享 EntryID 空间。
type engineState struct {
	scheduler *cron.Cron
}

var (
	// engineMu 保护 currentEngine、stopping 与 queueMap。
	engineMu sync.RWMutex
	// currentEngine 为 nil 表示当前没有可用调度引擎。
	currentEngine *engineState
	// stopping 为 true 后不再接受新的任务执行。
	stopping bool
	// queueMap 已注册的任务队列，进程内驻留，不随引擎启停重置。
	queueMap = map[string]*Queue{}

	// engineWG 跟踪已接受执行的任务，停机时统一排空。
	engineWG sync.WaitGroup
	// executeSem 限制同时执行的任务数量；派满时拒绝新执行。
	executeSem = make(chan struct{}, maxConcurrentJobs)
)

// lifecycleMu 串行化启停，避免 Stop 等待期间 Start 重开接收门禁。
var lifecycleMu sync.Mutex

// StartCron 启动（或重建）调度引擎。
//
// 重建前会等待旧调度与已接受任务执行完成，避免同 ID 任务在新旧引擎间重叠。
func StartCron() {
	lifecycleMu.Lock()
	defer lifecycleMu.Unlock()

	engineMu.Lock()
	var stopped context.Context
	if old := currentEngine; old != nil {
		stopping = true
		currentEngine = nil
		stopped = old.scheduler.Stop()
	}
	engineMu.Unlock()

	if stopped != nil {
		<-stopped.Done()
	}
	engineWG.Wait()

	engineMu.Lock()
	currentEngine = &engineState{scheduler: cron.New(cron.WithSeconds())}
	currentEngine.scheduler.Start()
	stopping = false
	engineMu.Unlock()
}

// StopCron 停止调度引擎，并等待已接受的任务及任务日志执行结束。
//
// 停止后 queueMap 仍在；再次 StartCron 后由调用方按需重新排程任务。
func StopCron() {
	lifecycleMu.Lock()
	defer lifecycleMu.Unlock()

	engineMu.Lock()
	var stopped context.Context
	if old := currentEngine; old != nil {
		stopping = true
		currentEngine = nil
		stopped = old.scheduler.Stop()
	}
	engineMu.Unlock()

	if stopped != nil {
		<-stopped.Done()
	}
	engineWG.Wait()
}

// scheduleCronJob 在引擎读锁内完成排程，避免旧引擎快照与停机操作竞争。
func scheduleCronJob(expr string, q *Queue, id string) (cron.EntryID, *engineState, error) {
	engineMu.RLock()
	defer engineMu.RUnlock()

	if stopping || currentEngine == nil {
		return 0, nil, ErrCronStopped
	}
	engine := currentEngine
	entryID, err := engine.scheduler.AddJob(expr, &jobRunner{queue: q, id: id, engine: engine})
	return entryID, engine, err
}

// removeCronJob 只从创建任务的引擎移除，避免跨引擎 EntryID 冲突。
func removeCronJob(engine *engineState, entryID cron.EntryID) {
	engineMu.RLock()
	defer engineMu.RUnlock()

	if engine != nil && engine == currentEngine {
		engine.scheduler.Remove(entryID)
	}
}

// dispatch 统一接受并启动任务执行；返回 false 表示停机或执行容量已满。
func dispatch(engine *engineState, q *Queue, id string) error {
	engineMu.Lock()
	defer engineMu.Unlock()

	if stopping || (engine != nil && engine != currentEngine) {
		return ErrCronStopped
	}
	select {
	case executeSem <- struct{}{}:
		engineWG.Add(1)
	default:
		return ErrExecutorFull
	}

	go func() {
		defer func() {
			<-executeSem
			engineWG.Done()
		}()
		q.runOnce(id)
	}()
	return nil
}

// CreateQueue 创建队列并注册处理器
func CreateQueue(name string, processor QueueProcessor) *Queue {
	engineMu.Lock()
	defer engineMu.Unlock()

	queue := &Queue{
		Name:      name,
		processor: processor,
		jobs:      map[string]*queueJob{},
		running:   map[string]int{},
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
	for name := range queueMap {
		keys = append(keys, name)
	}
	return keys
}
