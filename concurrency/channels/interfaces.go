package channels

import (
	"github.com/jiansoft/chopper/concurrency/core"
	"github.com/jiansoft/chopper/concurrency/fiber"
	"github.com/jiansoft/chopper/system"
)

//var log = logging.MustGetLogger("channels")

type IPublisher interface {
	Publish(interface{})
}

//Channel subscription methods.
type ISubscriber interface {
	Subscribe(fiber fiber.IFiber, taskFun interface{}, params ...interface{}) system.IDisposable
	ClearSubscribers()
}

type IChannel interface {
	SubscribeOnProducerThreads(subscriber IProducerThreadSubscriber) system.IDisposable
}

type IRequest interface {
	Request() interface{}
	SendReply(replyMsg interface{}) bool
}

type IReply interface {
	Receive(timeoutInMs int, result *interface{}) system.IDisposable
}

type IReplySubscriber interface {
	//Action<IRequest<TR, TM>>
	Subscribe(fiber fiber.IFiber, onRequest *interface{}) system.IDisposable
}

type IQueueChannel interface {
	Subscribe(executionContext core.IExecutionContext, onMessage interface{}) system.IDisposable
	Publish(message interface{})
}
type IProducerThreadSubscriber interface {
	//Allows for the registration and deregistration of subscriptions. IFiber
	Subscriptions() core.ISubscriptionRegistry
	/*Method called from producer threads*/
	ReceiveOnProducerThread(msg ...interface{})
}
