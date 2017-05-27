package core

import (
	"github.com/jiansoft/chopper/system"
)

//Allows for the registration and deregistration of subscriptions /*The IFiber has implemented*/
type ISubscriptionRegistry interface {
	//Register subscription to be unsubcribed from when the fiber is disposed.
	RegisterSubscription(system.IDisposable)
	//Deregister a subscription.
	DeregisterSubscription(system.IDisposable)
}

//Action subscriber that receives actions on producer thread.
type IProducerThreadSubscriber interface {
	ReceiveOnProducerThread(msg interface{})
	Subscriptions() ISubscriptionRegistry
}

type IExecutionContext interface {
	Enqueue(taskFun interface{}, params ...interface{})
	EnqueueWithTask(task Task)
}

type ISchedulerRegistry interface {
	Enqueue(taskFun interface{}, params ...interface{})
	Remove(d system.IDisposable)
}

type IScheduler interface {
	Schedule(firstInMs int64, taskFun interface{}, params ...interface{}) (d system.IDisposable)
	ScheduleOnInterval(firstInMs int64, regularInMs int64, taskFun interface{}, params ...interface{}) (d system.IDisposable)
	Start()
	Stop()
	Dispose()
}

type IExecutor interface {
	ExecuteTasks(t []Task)
	ExecuteTask(t Task)
}

type IQueue interface {
	Enqueue(t Task)
	DequeueAll() ([]Task, bool)
	Count() int
}
