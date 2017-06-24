package channels

import (
	"github.com/jiansoft/chopper/collections"
	"github.com/jiansoft/chopper/concurrency/core"
	"github.com/jiansoft/chopper/concurrency/fiber"
	"github.com/jiansoft/chopper/system"
)

type channel struct {
	subscribers collections.ConcurrentMap
}

func (c *channel) init() *channel {
	c.subscribers = collections.NewConcurrentMap()
	return c
}

func NewChannel() *channel {
	return new(channel).init()
}

func (c *channel) Subscribe(fiber fiber.IFiber, taskFun interface{}, params ...interface{}) system.IDisposable {
	subscription := NewChannelSubscription(fiber, core.Task{Func: taskFun, Params: params})
	return c.SubscribeOnProducerThreads(subscription)
}

func (c *channel) SubscribeOnProducerThreads(subscriber IProducerThreadSubscriber) system.IDisposable {
	job := core.Task{Func: subscriber.ReceiveOnProducerThread}
	return c.subscribeOnProducerThreads(job, subscriber.Subscriptions())
}

func (c *channel) subscribeOnProducerThreads(subscriber core.Task, fiber core.ISubscriptionRegistry) system.IDisposable {
	unsubscriber := NewUnsubscriber(subscriber, c, fiber)
	//將訂閱者的方法註冊到 IFiber內 ，當 IFiber.Dispose()時，同步將訂閱的方法移除
	fiber.RegisterSubscription(unsubscriber)
	//放到Channel 內的貯列，當 Chanel.Publish 時發布給訂閱的方法
	c.subscribers.Set(unsubscriber.IdentifyId(), unsubscriber)
	return unsubscriber
}

func (c *channel) Publish(msg ...interface{}) {
	for _, val := range c.subscribers.Items() {
		val.(*unsubscriber).fiber.(fiber.IFiber).Enqueue(val.(*unsubscriber).receiver.Func, msg...)
	}
}

func (c *channel) unsubscribe(disposable system.IDisposable) {
	//Remove the subscriber
	if val, ok := c.subscribers.Get(disposable.IdentifyId()); ok {
		val.(*unsubscriber).fiber.DeregisterSubscription(val.(*unsubscriber))
		c.subscribers.Remove(disposable.IdentifyId())
	}
}

func (c *channel) NumSubscribers() int {
	return c.subscribers.Count()
}
