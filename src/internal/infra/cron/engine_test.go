package cron

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"mask_api_gin/src/internal/shared/constants"
)

// probeProc 可控的测试处理器：可选择返回值、错误、panic 或阻塞到 release。
type probeProc struct {
	block  bool
	err    error
	panicV string

	count   atomic.Int64
	entered chan struct{}
	release chan struct{}
}

func newProbe() *probeProc {
	return &probeProc{
		entered: make(chan struct{}, 64),
		release: make(chan struct{}),
	}
}

func (p *probeProc) Execute(data any) (any, error) {
	p.count.Add(1)
	p.entered <- struct{}{}
	if p.panicV != "" {
		panic(p.panicV)
	}
	if p.err != nil {
		return nil, p.err
	}
	if p.block {
		<-p.release
	}
	return data, nil
}

func (p *probeProc) calls() int64 { return p.count.Load() }

func (p *probeProc) waitEntered(t *testing.T) {
	t.Helper()
	select {
	case <-p.entered:
	case <-time.After(time.Second):
		t.Fatal("processor 未在 1s 内进入 Execute")
	}
}

// captureSink 安装一个收集所有记录的 sink，返回读取快照的闭包。
func captureSink() (func() []JobLogRecord, func()) {
	var mu sync.Mutex
	var records []JobLogRecord
	prev := OnJobLog
	OnJobLog = func(r JobLogRecord) {
		mu.Lock()
		records = append(records, r)
		mu.Unlock()
	}
	snapshot := func() []JobLogRecord {
		mu.Lock()
		defer mu.Unlock()
		out := make([]JobLogRecord, len(records))
		copy(out, records)
		return out
	}
	restore := func() { OnJobLog = prev }
	return snapshot, restore
}

// waitFor 等待条件成立或超时。
func waitFor(t *testing.T, fn func() bool, msg string) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for !fn() {
		if time.Now().After(deadline) {
			t.Fatal(msg)
		}
		time.Sleep(2 * time.Millisecond)
	}
}

func TestRunOnceSuccessAndCleanup(t *testing.T) {
	StartCron()
	defer StopCron()
	snap, restore := captureSink()
	defer restore()

	proc := newProbe()
	q := CreateQueue("t_success", proc)
	if status := q.RunJob("hello", JobOptions{ID: "s1"}); status != Waiting {
		t.Fatalf("RunJob 返回 %d，期望 Waiting", status)
	}

	waitFor(t, func() bool { return len(snap()) == 1 }, "未收到执行事件")
	rec := snap()[0]
	if rec.Status != constants.STATUS_YES {
		t.Fatalf("status=%s 期望 YES", rec.Status)
	}
	if rec.Result != "hello" {
		t.Fatalf("result=%v 期望 hello", rec.Result)
	}

	// 一次性任务执行后应从 map 清理
	q.mu.Lock()
	_, ok := q.jobs["s1"]
	q.mu.Unlock()
	if ok {
		t.Fatal("一次性任务执行后未从 map 清理")
	}
}

func TestRunOnceFailure(t *testing.T) {
	StartCron()
	defer StopCron()
	snap, restore := captureSink()
	defer restore()

	wantErr := errors.New("boom")
	proc := newProbe()
	proc.err = wantErr
	q := CreateQueue("t_fail", proc)
	q.RunJob(nil, JobOptions{ID: "f1"})

	waitFor(t, func() bool { return len(snap()) == 1 }, "未收到执行事件")
	rec := snap()[0]
	if rec.Status != constants.STATUS_NO {
		t.Fatalf("status=%s 期望 NO", rec.Status)
	}
	if rec.Result != wantErr.Error() {
		t.Fatalf("result=%v 期望 %s", rec.Result, wantErr.Error())
	}
}

func TestRunOncePanic(t *testing.T) {
	StartCron()
	defer StopCron()
	snap, restore := captureSink()
	defer restore()

	proc := newProbe()
	proc.panicV = "kaboom"
	q := CreateQueue("t_panic", proc)
	q.RunJob(nil, JobOptions{ID: "p1"})

	waitFor(t, func() bool { return len(snap()) == 1 }, "未收到执行事件")
	rec := snap()[0]
	if rec.Status != constants.STATUS_NO {
		t.Fatalf("panic 后 status=%s 期望 NO", rec.Status)
	}
	if rec.Result != "kaboom" {
		t.Fatalf("result=%v 期望 kaboom", rec.Result)
	}
}

// 非并发模式：上一次仍在执行时，本次触发应被跳过。
func TestConcurrencySkip(t *testing.T) {
	StartCron()
	defer StopCron()

	proc := newProbe()
	proc.block = true
	q := CreateQueue("t_skip", proc)
	q.RunJob(nil, JobOptions{ID: "c1", Concurrent: false})

	proc.waitEntered(t) // 确保已进入 Active
	// 模拟并发触发：直接再跑一次（同 ID、非并发）应被跳过
	q.runOnce("c1")
	if got := proc.calls(); got != 1 {
		t.Fatalf("processor 调用 %d 次，期望 1（应跳过重叠）", got)
	}

	close(proc.release) // 放行
	waitFor(t, func() bool {
		q.mu.Lock()
		_, ok := q.jobs["c1"]
		q.mu.Unlock()
		return !ok
	}, "一次性任务执行后未从 map 清理")
}

// 并发模式：允许上一次未完成即再次执行。
func TestConcurrencyAllow(t *testing.T) {
	StartCron()
	defer StopCron()

	proc := newProbe()
	proc.block = true
	q := CreateQueue("t_allow", proc)
	q.RunJob(nil, JobOptions{ID: "c2", Concurrent: true})
	proc.waitEntered(t)

	// 并发允许：第二次触发应真正执行
	go q.runOnce("c2")
	waitFor(t, func() bool { return proc.calls() == 2 }, "第二次未执行")

	close(proc.release)
	// 两次执行完成后，一次性任务应从 map 清理
	waitFor(t, func() bool {
		q.mu.Lock()
		_, ok := q.jobs["c2"]
		q.mu.Unlock()
		return !ok
	}, "并发执行后未从 map 清理")
}

// 引擎停止后，重复任务应无法排程（返回 Failed）。
func TestStopBlocksRepeatScheduling(t *testing.T) {
	StartCron()
	StopCron()
	defer StartCron()

	proc := newProbe()
	q := CreateQueue("t_stop", proc)
	status := q.RunJob(nil, JobOptions{ID: "r1", Cron: "*/1 * * * * *"})
	if status != Failed {
		t.Fatalf("停机后排程返回 %d，期望 Failed", status)
	}
	if got := proc.calls(); got != 0 {
		t.Fatalf("停机后 processor 仍被调用 %d 次", got)
	}
}
