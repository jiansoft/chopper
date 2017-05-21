package fiber

import (
	"sync"

	"github.com/jiansoft/chopper/concurrency/core"
	"github.com/jiansoft/chopper/system"
)

type FiberPool struct {
	scheduler      core.IScheduler
	executor       core.IExecutor
	tasks          []core.Task
	swapTasks      []core.Task
	executionState executionState
	lock           *sync.Mutex
	flushPending   bool
	subscriptions  *core.Disposer
}

func (f *FiberPool) init() *FiberPool {
	f.executionState = created
	f.tasks = []core.Task{}
	f.swapTasks = []core.Task{}
	f.scheduler = core.NewScheduler(f)
	f.executor = core.NewDefaultExecutor()
	f.subscriptions = core.NewDisposer()
	f.lock = new(sync.Mutex)
	return f
}

func NewFiberPool() *FiberPool {
	return new(FiberPool).init()
}

func (f *FiberPool) Start() {
	if f.executionState == running {
		return
	}
	f.executionState = running
	f.scheduler.Start()
	f.Enqueue(func() {})
}

func (f *FiberPool) Stop() {
	f.executionState = stopped
}

func (f *FiberPool) Dispose() {
	f.Stop()
	f.scheduler.Dispose()
	f.subscriptions.Dispose()
	f.clearTasks()
}

func (f *FiberPool) Enqueue(taskFun interface{}, params ...interface{}) {
	f.EnqueueWithTask(core.Task{PaddingFunc: taskFun, FuncParams: params})
}

func (f *FiberPool) EnqueueWithTask(task core.Task) {
	if f.executionState == stopped {
		return
	}

	f.lock.Lock()
	defer f.lock.Unlock()
	f.tasks = append(f.tasks, task)
	if f.executionState == created || f.flushPending {
		return
	}
	f.queue(f.flush)
	f.flushPending = true
}

func (f *FiberPool) Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d system.IDisposable) {
	return f.scheduler.Schedule(firstInMs, taskFun, params...)
}

func (f *FiberPool) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d system.IDisposable) {
	return f.scheduler.ScheduleOnInterval(firstInMs, regularInMs, taskFun, params...)
}

/*implement ISubscriptionRegistry.RegisterSubscription */
func (f *FiberPool) RegisterSubscription(toAdd system.IDisposable) {
	f.subscriptions.Add(toAdd)
}

/*implement ISubscriptionRegistry.DeregisterSubscription */
func (f *FiberPool) DeregisterSubscription(toRemove system.IDisposable) {
	f.subscriptions.Remove(toRemove)
}

func (f FiberPool) NumSubscriptions() int {
	return f.subscriptions.Count()
}

func (f *FiberPool) flush() {
	toDoTasks, ok := f.clearTasks()
	if !ok {
		return
	}
	f.executor.ExecuteTasks(toDoTasks)
	f.lock.Lock()
	if len(f.tasks) > 0 {
		//有任務在清空任務清單時又新增了任務
        //It has some new task enqueue when clear tasks
		f.queue(f.flush)
	} else {
		//Task is empty
		f.flushPending = false
	}
	f.lock.Unlock()
}

func (f *FiberPool) clearTasks() ([]core.Task, bool) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if len(f.tasks) != 0 {
		f.swapTasks, f.tasks = f.tasks, f.swapTasks
		f.tasks = f.tasks[:0]
		return f.swapTasks, true
	}
	f.flushPending = false
	return nil, false
}

func (f FiberPool) queue(taskFun interface{}) {
	go f.executor.ExecuteTask(core.Task{PaddingFunc: taskFun})
	//go core.Task{PaddingFunc: taskFun, FuncParams: params}.Execute()
}
