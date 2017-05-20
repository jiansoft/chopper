package fiber

import (
	"sync"

	"github.com/jiansoft/chopper/concurrency/core"
	"github.com/jiansoft/chopper/system"
)

type FiberThread struct {
	scheduler      core.IScheduler
	executor       core.IExecutor
	tasks          []core.Task
	swapTasks      []core.Task
	executionState executionState
	lock           *sync.Mutex
	cond           *sync.Cond
	subscriptions  *core.Disposer
}

func (f *FiberThread) init() *FiberThread {
	f.executionState = created
	f.subscriptions = core.NewDisposer()
	f.swapTasks = []core.Task{}
	f.scheduler = core.NewScheduler(f)
	f.executor = core.NewDefaultExecutor()
	f.tasks = []core.Task{}
	f.lock = new(sync.Mutex)
	f.cond = sync.NewCond(f.lock)
	return f
}

func NewFiberThread() *FiberThread {
	return new(FiberThread).init()
}

func (f *FiberThread) Start() {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.executionState == running {
		return
	}
	f.executionState = running
	f.scheduler.Start()
	go func() {
		for f.executeNextBatch() {
		}
		//log.Infof("exit goroutine")
	}()
}
func (f *FiberThread) Stop() {
	f.lock.Lock()
	//f.running = false
	f.executionState = stopped
	f.cond.Broadcast()
	f.lock.Unlock()
	f.scheduler.Stop()
}

func (f *FiberThread) Dispose() {
	f.lock.Lock()
	//f.running = false
	f.executionState = stopped
	f.cond.Broadcast()
	f.lock.Unlock()
	f.scheduler.Dispose()
	f.subscriptions.Dispose()
}

func (f *FiberThread) Enqueue(taskFun interface{}, params ...interface{}) {
	f.EnqueueWithTask(core.Task{PaddingFunc: taskFun, FuncParams: params})
}

func (f *FiberThread) EnqueueWithTask(task core.Task) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.tasks = append(f.tasks, task)
	//喚醒等待中的 Goroutine
	f.cond.Broadcast()
}

func (f FiberThread) Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d system.IDisposable) {
	return f.scheduler.Schedule(firstInMs, taskFun, params...)
}

func (f FiberThread) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d system.IDisposable) {
	return f.scheduler.ScheduleOnInterval(firstInMs, regularInMs, taskFun, params...)
}

/*implement ISubscriptionRegistry.RegisterSubscription */
func (f *FiberThread) RegisterSubscription(toAdd system.IDisposable) {
	f.subscriptions.Add(toAdd)
}

/*implement ISubscriptionRegistry.DeregisterSubscription */
func (f *FiberThread) DeregisterSubscription(toRemove system.IDisposable) {
	f.subscriptions.Remove(toRemove)
}

func (f *FiberThread) NumSubscriptions() int {
	return f.subscriptions.Count()
}

func (f *FiberThread) executeNextBatch() bool {
	toExecute, ok := f.dequeueAll()
	if ok {
		f.executor.ExecuteTasks(toExecute)
	}
	return ok
}

func (f *FiberThread) dequeueAll() ([]core.Task, bool) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if !f.readyToDequeue() {
		return nil, false
	}
	f.swapTasks, f.tasks = f.tasks, f.swapTasks
	f.tasks = f.tasks[:0]
	return f.swapTasks, true
}

func (f *FiberThread) readyToDequeue() bool {
	//若貯列中已沒有任務要執行時，將 Goroutine 進入等侯訊號的狀態
	tmp := f.executionState == running
	for len(f.tasks) == 0 && tmp {
		f.cond.Wait()
	}
	return tmp
}
