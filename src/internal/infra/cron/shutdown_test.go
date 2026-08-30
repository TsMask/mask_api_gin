package cron

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func waitStopResult(t *testing.T, stopped <-chan struct{}) bool {
	t.Helper()
	select {
	case <-stopped:
		return true
	case <-time.After(30 * time.Millisecond):
		return false
	}
}

func TestStopCronWaitsForOneShotJobAndLogSink(t *testing.T) {
	StartCron()
	defer StopCron()

	proc := newProbe()
	proc.block = true
	q := CreateQueue("t_shutdown_one_shot", proc)

	sinkEntered := make(chan struct{})
	sinkRelease := make(chan struct{})
	SetOnJobLog(func(JobLogRecord) {
		sinkEntered <- struct{}{}
		<-sinkRelease
	})
	defer SetOnJobLog(nil)

	if status := q.RunJob(nil, JobOptions{ID: "shutdown-one-shot"}); status != Waiting {
		t.Fatalf("RunJob status=%d, want Waiting", status)
	}
	proc.waitEntered(t)

	stopped := make(chan struct{})
	go func() {
		StopCron()
		close(stopped)
	}()
	if waitStopResult(t, stopped) {
		t.Fatal("StopCron returned before processor finished")
	}

	close(proc.release)
	select {
	case <-sinkEntered:
	case <-time.After(time.Second):
		t.Fatal("job log sink was not called after processor finished")
	}
	if waitStopResult(t, stopped) {
		t.Fatal("StopCron returned before job log sink finished")
	}

	close(sinkRelease)
	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("StopCron did not return after job log sink finished")
	}

	// StopCron 已完成，恢复引擎供 defer 再停一次时保持测试状态干净。
	StartCron()
}

func TestStopCronWaitsForScheduledJob(t *testing.T) {
	StartCron()
	defer StopCron()

	proc := newProbe()
	proc.block = true
	q := CreateQueue("t_shutdown_scheduled", proc)
	status := q.RunJob(nil, JobOptions{
		ID:         "shutdown-scheduled",
		Cron:       "* * * * * *",
		Concurrent: false,
	})
	if status != Waiting {
		t.Fatalf("RunJob status=%d, want Waiting", status)
	}
	proc.waitEntered(t)

	stopped := make(chan struct{})
	go func() {
		StopCron()
		close(stopped)
	}()
	if waitStopResult(t, stopped) {
		t.Fatal("StopCron returned before scheduled processor finished")
	}

	close(proc.release)
	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("StopCron did not return after scheduled processor finished")
	}

	StartCron()
}

func TestRunJobRejectsDuplicateOneShotWhenNonConcurrent(t *testing.T) {
	StartCron()
	defer StopCron()

	proc := newProbe()
	proc.block = true
	q := CreateQueue("t_duplicate_one_shot", proc)
	if status := q.RunJob(nil, JobOptions{ID: "duplicate"}); status != Waiting {
		t.Fatalf("first RunJob status=%d, want Waiting", status)
	}
	proc.waitEntered(t)

	if status := q.RunJob(nil, JobOptions{ID: "duplicate"}); status != Active {
		t.Fatalf("second RunJob status=%d, want Active", status)
	}
	if got := proc.calls(); got != 1 {
		t.Fatalf("processor calls=%d, want 1", got)
	}

	close(proc.release)
	waitFor(t, func() bool {
		q.mu.Lock()
		defer q.mu.Unlock()
		return len(q.running) == 0
	}, "running count was not cleared after processor finished")
}

func TestDispatchRejectsJobWhenExecutorIsFull(t *testing.T) {
	StartCron()
	defer StopCron()

	proc := newProbe()
	proc.block = true
	q := CreateQueue("t_executor_full", proc)
	engineErrors := make(chan EngineError, 1)
	SetOnEngineError(func(event EngineError) {
		engineErrors <- event
	})
	defer SetOnEngineError(nil)
	for i := 0; i < maxConcurrentJobs; i++ {
		id := fmt.Sprintf("job-%02d", i)
		if status := q.RunJob(nil, JobOptions{ID: id, Concurrent: true}); status != Waiting {
			t.Fatalf("RunJob id=%s status=%d, want Waiting", id, status)
		}
	}
	waitFor(t, func() bool { return proc.calls() == maxConcurrentJobs }, "not all accepted jobs entered Execute")

	if status := q.RunJob(nil, JobOptions{ID: "overflow", Concurrent: true}); status != Failed {
		t.Fatalf("overflow RunJob status=%d, want Failed", status)
	}
	select {
	case event := <-engineErrors:
		if !errors.Is(event.Err, ErrExecutorFull) {
			t.Fatalf("engine error=%v, want ErrExecutorFull", event.Err)
		}
	default:
		t.Fatal("overflow rejection was not emitted to engine error sink")
	}
	if got := proc.calls(); got != maxConcurrentJobs {
		t.Fatalf("processor calls=%d, want %d", got, maxConcurrentJobs)
	}

	close(proc.release)
}

func TestEmitJobLogForwardsSinkPanicAsEngineError(t *testing.T) {
	StartCron()
	defer StopCron()

	proc := newProbe()
	q := CreateQueue("t_sink_panic", proc)
	SetOnJobLog(func(JobLogRecord) {
		panic("sink boom")
	})
	defer SetOnJobLog(nil)
	engineErrors := make(chan EngineError, 1)
	SetOnEngineError(func(event EngineError) {
		engineErrors <- event
	})
	defer SetOnEngineError(nil)

	q.RunJob(nil, JobOptions{ID: "sink-panic"})
	waitFor(t, func() bool { return proc.calls() == 1 }, "processor did not finish")
	StopCron()
	select {
	case event := <-engineErrors:
		if event.Queue != q.Name || event.ID != "sink-panic" {
			t.Fatalf("engine error queue=%s id=%s, want %s/sink-panic", event.Queue, event.ID, q.Name)
		}
	default:
		t.Fatal("sink panic was not emitted to engine error sink")
	}
	StartCron()
}
