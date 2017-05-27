package core

import (
	"fmt"
	"time"
)

type timerTask struct {
	identifyId    string
	fiber         ISchedulerRegistry
	firstInMs     int64
	intervalInMs  int64
	firstTimer    *time.Ticker
	intervalTimer *time.Ticker
	task          Task
	cancelled     bool
	chanClose     chan bool
}

//初始化
func (t *timerTask) init(fiber ISchedulerRegistry, task Task, firstInMs int64, intervalInMs int64) *timerTask {
	t.fiber = fiber
	t.task = task
	t.firstInMs = firstInMs
	t.intervalInMs = intervalInMs
	t.identifyId = fmt.Sprintf("%p-%p", &t, &task)
	t.chanClose = make(chan bool)
	return t
}

func newTimerTask(fiber ISchedulerRegistry, task Task, firstInMs int64, intervalInMs int64) *timerTask {
	return new(timerTask).init(fiber, task, firstInMs, intervalInMs)
}

func (t *timerTask) Dispose() {
	t.cancelled = true
	t.chanClose <- t.cancelled
	t.fiber.Remove(t)
}

func (t timerTask) IdentifyId() string {
	return t.identifyId
}

func (t *timerTask) schedule() {
	if t.firstInMs <= 0 {
		t.doFirstSchedule()
	} else {
		t.firstTimer = time.NewTicker(time.Duration(t.firstInMs) * time.Millisecond)
		go func() {
			select {
			case <-t.firstTimer.C:
				t.firstTimer.Stop()
				if !t.cancelled {
					t.doFirstSchedule()
				}
			case _ = <-t.chanClose:
				if nil != t.firstTimer {
					t.firstTimer.Stop()
				}
				if nil != t.intervalTimer {
					t.intervalTimer.Stop()
				}

			}
		}()
	}
}

func (t *timerTask) doFirstSchedule() {
	t.fiber.Enqueue(t.executeOnFiberThread)
	t.doIntervalSchedule()
}

func (t *timerTask) doIntervalSchedule() {
	if t.intervalInMs <= 0 {
		return
	}
	t.intervalTimer = time.NewTicker(time.Duration(t.intervalInMs) * time.Millisecond)
	go func() {
		for !t.cancelled {
			select {
			case <-t.intervalTimer.C:
				if t.cancelled {
					t.intervalTimer.Stop()
					break
				}
				t.fiber.Enqueue(t.executeOnFiberThread)
			case _ = <-t.chanClose:
				if nil != t.intervalTimer {
					t.intervalTimer.Stop()
				}
				break
			}
		}
	}()
}

func (t timerTask) executeOnFiberThread() {
	if t.cancelled {
		return
	}
	t.task.Execute()
}
