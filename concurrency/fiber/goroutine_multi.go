package fiber

import (
	"sync"

	"github.com/jiansoft/chopper/concurrency/core"
	"github.com/jiansoft/chopper/system"
)

type GoroutineMulti struct {
	queue          core.IQueue
	scheduler      core.IScheduler
	executor       core.IExecutor
	executionState executionState
	lock           *sync.Mutex
	subscriptions  *core.Disposer
	flushPending   bool
}

func (g *GoroutineMulti) init() *GoroutineMulti {
	g.queue = core.NewDefaultQueue()
	g.executionState = created
	g.scheduler = core.NewScheduler(g)
	g.executor = core.NewDefaultExecutor()
	g.subscriptions = core.NewDisposer()
	g.lock = new(sync.Mutex)
	return g
}

func NewGoroutineMulti() *GoroutineMulti {
	return new(GoroutineMulti).init()
}

func (g *GoroutineMulti) Start() {
	if g.executionState == running {
		return
	}
	g.executionState = running
	g.scheduler.Start()
	g.Enqueue(func() {})
}

func (g *GoroutineMulti) Stop() {
	g.executionState = stopped
}

func (g *GoroutineMulti) Dispose() {
	g.Stop()
	g.scheduler.Dispose()
	g.subscriptions.Dispose()
	g.queue.DequeueAll()
}

func (g *GoroutineMulti) Enqueue(taskFun interface{}, params ...interface{}) {
	g.EnqueueWithTask(core.Task{Func: taskFun, Params: params})
}

func (g *GoroutineMulti) EnqueueWithTask(task core.Task) {
	if g.executionState != running {
		return
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	g.queue.Enqueue(task)
	if g.flushPending {
		return
	}
	g.run(g.flush)
	g.flushPending = true
}

func (g *GoroutineMulti) Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d system.IDisposable) {
	return g.scheduler.Schedule(firstInMs, taskFun, params...)
}

func (g *GoroutineMulti) ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d system.IDisposable) {
	return g.scheduler.ScheduleOnInterval(firstInMs, regularInMs, taskFun, params...)
}

/*implement ISubscriptionRegistry.RegisterSubscription */
func (g *GoroutineMulti) RegisterSubscription(toAdd system.IDisposable) {
	g.subscriptions.Add(toAdd)
}

/*implement ISubscriptionRegistry.DeregisterSubscription */
func (g *GoroutineMulti) DeregisterSubscription(toRemove system.IDisposable) {
	g.subscriptions.Remove(toRemove)
}

func (g GoroutineMulti) NumSubscriptions() int {
	return g.subscriptions.Count()
}

func (g *GoroutineMulti) flush() {
	toDoTasks, ok := g.queue.DequeueAll()
	if !ok {
		g.flushPending = false
		return
	}
	g.executor.ExecuteTasks(toDoTasks)
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.queue.Count() > 0 {
		//It has new task enqueue when clear tasks
		g.run(g.flush)
	} else {
		//Task is empty
		g.flushPending = false
	}
}

func (g GoroutineMulti) run(taskFun interface{}) {
	go core.Task{Func: taskFun}.Execute()
}
