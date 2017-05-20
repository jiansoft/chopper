package fiber

import (
	"github.com/jiansoft/chopper/concurrency/core"
	"github.com/jiansoft/chopper/system"
)

type IThreadPool interface {
	Queue(t core.Task)
}

type IQueue interface {
	Enqueue(taskFun interface{}, params ...interface{})
	Run()
	Stop()
}

type IFiber interface {
	Start()
	Stop()
	Dispose()
	Enqueue(taskFun interface{}, params ...interface{})
	EnqueueWithTask(task core.Task)
	Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d system.IDisposable)
	ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d system.IDisposable)
}
