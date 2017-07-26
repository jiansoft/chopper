package core

import (
	"fmt"
	"time"
)

type timerTask struct {
	identifyId   string
	scheduler    ISchedulerRegistry
	firstInMs    int64
	intervalInMs int64
	timer        *time.Ticker
	task         Task
	cancelled    bool
	firstRan     bool
}

//初始化
func (t *timerTask) init(scheduler ISchedulerRegistry, task Task, firstInMs int64, intervalInMs int64) *timerTask {
	t.scheduler = scheduler
	t.task = task
	t.firstInMs = firstInMs
	t.intervalInMs = intervalInMs
	t.identifyId = fmt.Sprintf("%p-%p", &t, &task)
	return t
}

func newTimerTask(fiber ISchedulerRegistry, task Task, firstInMs int64, intervalInMs int64) *timerTask {
	return new(timerTask).init(fiber, task, firstInMs, intervalInMs)
}

func (t *timerTask) Dispose() {
	t.cancelled = true
	if nil != t.timer {
		t.timer.Stop()
		t.timer = nil
	}
	if nil != t.scheduler {
		t.scheduler.Remove(t)
		t.scheduler = nil
	}

}

func (t timerTask) IdentifyId() string {
	return t.identifyId
}

func (t *timerTask) schedule() {
	timeInMs := t.firstInMs
	if timeInMs <= 0 {
		t.scheduler.EnqueueWithTask(t.task)
		timeInMs = t.intervalInMs
	}
	t.timer = time.NewTicker(time.Duration(timeInMs) * time.Millisecond)
	go func() {
		for !t.cancelled {
			select {
			case <-t.timer.C:
				if t.cancelled {
					break
				}
				t.scheduler.EnqueueWithTask(t.task)
				if t.intervalInMs <= 0 {
					t.Dispose()
					break
				}
				if t.firstRan {
					continue
				}
				t.timer = time.NewTicker(time.Duration(t.intervalInMs) * time.Millisecond)
				t.firstRan = true
			}
		}
	}()
}
