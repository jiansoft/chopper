package core

import "github.com/jiansoft/chopper/system"

type Scheduler struct {
	fiber       IExecutionContext
	running     bool
	isDispose   bool
	disposabler *Disposer
}

func (s *Scheduler) init(executionState IExecutionContext) *Scheduler {
	s.fiber = executionState
	s.running = true
	s.disposabler = NewDisposer()
	return s
}

func NewScheduler(executionState IExecutionContext) *Scheduler {
	return new(Scheduler).init(executionState)
}

func (s *Scheduler) Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d system.IDisposable) {
	if firstInMs <= 0 {
		pendingAction := newPendingTask(Task{PaddingFunc: taskFun, FuncParams: params})
		s.fiber.Enqueue(pendingAction.Execute)
		return pendingAction
	}
	pending := newTimerTask(s, Task{PaddingFunc: taskFun, FuncParams: params}, firstInMs, -1)
	//log.Infof("addPending pending:%p",&pending)
	s.addPending(pending)
	return pending
}

func (s *Scheduler) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d system.IDisposable) {
	pending := newTimerTask(s, Task{PaddingFunc: taskFun, FuncParams: params}, firstInMs, regularInMs)
	//log.Infof("addPending pending:%p",&pending)
	s.addPending(pending)
	return pending
}

//實作 ISchedulerRegistry.Enqueue
func (s *Scheduler) Enqueue(taskFun interface{}, params ...interface{}) {
	//s.fiber.Enqueue(taskFun, params...)
	s.fiber.EnqueueWithTask(Task{FuncParams: params, PaddingFunc: taskFun})
}

//實作 ISchedulerRegistry.Remove
func (s *Scheduler) Remove(d system.IDisposable) {
	s.fiber.Enqueue(s.disposabler.Remove, d)
}

func (s *Scheduler) Start() {
	s.running = true
	s.isDispose = false
}

func (s *Scheduler) Stop() {
	s.running = false
}

func (s *Scheduler) Dispose() {
	s.Stop()
	s.isDispose = true
	s.disposabler.Dispose()
}

func (s *Scheduler) addPending(pending *timerTask) {
	s.fiber.Enqueue(func() {
		if s.isDispose {
			return
		}
		s.disposabler.Add(pending)
		pending.schedule()
	})
}
