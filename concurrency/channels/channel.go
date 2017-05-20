package channels

import (
	"sync"

	"github.com/jiansoft/chopper/collections"
	"github.com/jiansoft/chopper/concurrency/core"
	"github.com/jiansoft/chopper/concurrency/fiber"
	"github.com/jiansoft/chopper/system"
)

type channel struct {
	subscribers collections.ConcurrentMap
	lock        *sync.Mutex
}

func (c *channel) init() *channel {
	c.subscribers = collections.NewConcurrentMap()
	c.lock = new(sync.Mutex)
	return c
}

func NewChannel() *channel {
	return new(channel).init()
}

func (c *channel) Subscribe(fiber fiber.IFiber, taskFun interface{}, params ...interface{}) system.IDisposable {
	return c.SubscribeOnProducerThreads(NewChannelSubscription(fiber, core.Task{PaddingFunc: taskFun, FuncParams: params}))
}

func (c *channel) SubscribeOnProducerThreads(subscriber IProducerThreadSubscriber) system.IDisposable {
	job := core.Task{PaddingFunc: subscriber.ReceiveOnProducerThread}
	return c.subscribeOnProducerThreads(job, subscriber.Subscriptions())
}

func (c *channel) subscribeOnProducerThreads(subscriber core.Task, fiber core.ISubscriptionRegistry) system.IDisposable {
	//再封裝一層
	unsubscriber := NewUnsubscriber(subscriber, c, fiber)
	//將訂閱者的方法註冊到 IFiber內 ，當 IFiber.Dispose()時，同步將訂閱的方法移除
	fiber.RegisterSubscription(unsubscriber)
	//放到Channel 內的貯列，當 Chanel.Publish 時發布給訂閱的方法
	c.subscribers.Set(unsubscriber.IdentifyId(), unsubscriber)
	return unsubscriber

}

func (c *channel) Publish(msg ...interface{}) {
	for _, val := range c.subscribers.Items() {
		val.(*unsubscriber).receiver.ExecuteWithParams(msg...)
	}
}

func (c *channel) Unsubscribe(disposable system.IDisposable) {
	//將訂閱的方法移除
	c.subscribers.Remove(disposable.IdentifyId())
}

func (c *channel) NumSubscribers() int {
	return c.subscribers.Count()
}
