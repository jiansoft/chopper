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

//Action subscriber that receives actions on producer thread.
type IProducerThreadSubscriber interface {
	ReceiveOnProducerThread(msg interface{})
	Subscriptions() ISubscriptionRegistry
}

type IExecutionContext interface {
	Enqueue(taskFun interface{}, params ...interface{})
	EnqueueWithTask(task Task)
}
