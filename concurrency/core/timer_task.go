package core

import (
	"fmt"
	"time"
	// "log"
)

type timerTask struct {
	identifyId         string
	fiber              ISchedulerRegistry
	firstIntervalInMs  int64
	intervalInMs       int64
	firstIntervalTimer *time.Ticker
	intervalTimer      *time.Ticker
	task               Task
	cancelled          bool
}

//初始化
func (t *timerTask) init(registry ISchedulerRegistry, task Task, firstIntervalInMs int64, intervalInMs int64) *timerTask {
	t.fiber = registry
	t.task = task
	t.firstIntervalInMs = firstIntervalInMs
	t.intervalInMs = intervalInMs
	t.identifyId = fmt.Sprintf("%p-%p", &t, &task)
	return t
}

func newTimerTask(registry ISchedulerRegistry, task Task, firstIntervalInMs int64, intervalInMs int64) *timerTask {
	return new(timerTask).init(registry, task, firstIntervalInMs, intervalInMs)
}

func (t *timerTask) Dispose() {
	t.cancelled = true
	t.fiber.Remove(t)
}

func (t timerTask) IdentifyId() string {
	return t.identifyId
}

func (t *timerTask) schedule() {
	if t.firstIntervalInMs <= 0 {
		t.doFirstIntervalSchedule()
	} else {
		t.firstIntervalTimer = time.NewTicker(time.Duration(t.firstIntervalInMs) * time.Millisecond)
		go func() {
			//for {
			select {
			case <-t.firstIntervalTimer.C:
				t.firstIntervalTimer.Stop()
				if !t.cancelled {
					t.doFirstIntervalSchedule()
				}
			}
			//    break
			//}
			//log.Printf("firstIntervalTimer exit")
		}()
	}
}

func (t *timerTask) doFirstIntervalSchedule() {
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
