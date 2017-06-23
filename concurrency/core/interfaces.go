package core

import (
	"github.com/jiansoft/chopper/system"
)

//Allows for the registration and deregistration of subscriptions /*The IFiber has implemented*/
type ISubscriptionRegistry interface {
	//Register subscription to be unsubcribed from when the scheduler is disposed.
	RegisterSubscription(system.IDisposable)
	//Deregister a subscription.
	DeregisterSubscription(system.IDisposable)
}

type IExecutionContext interface {
	Enqueue(taskFun interface{}, params ...interface{})
	EnqueueWithTask(task Task)
}
