package channels

import (
	"fmt"

	"github.com/jiansoft/chopper/concurrency/core"
	"github.com/jiansoft/chopper/concurrency/fiber"
)

type channelSubscription struct {
	identifyId string
	fiber      fiber.IFiber
	receiver   core.Task
}

func (c *channelSubscription) init(fiber fiber.IFiber, task core.Task) *channelSubscription {
	c.fiber = fiber
	c.receiver = task
	c.identifyId = fmt.Sprintf("%p-%p", &c, &task)
	return c
}

func NewChannelSubscription(fiber fiber.IFiber, task core.Task) *channelSubscription {
	return new(channelSubscription).init(fiber, task)
}

func (c *channelSubscription) OnMessageOnProducerThread(msg ...interface{}) {
	c.fiber.Enqueue(c.receiver.Func, msg...)
}

//實作 IProducerThreadSubscriber.Subscriptions
func (c *channelSubscription) Subscriptions() core.ISubscriptionRegistry {
	return c.fiber.(core.ISubscriptionRegistry)
}

//實作 IProducerThreadSubscriber.ReceiveOnProducerThread
func (c *channelSubscription) ReceiveOnProducerThread(msg ...interface{}) {
	c.OnMessageOnProducerThread(msg...)
}

func (c *channelSubscription) Dispose() {
	c.fiber.(core.ISubscriptionRegistry).DeregisterSubscription(c)
}

func (c *channelSubscription) IdentifyId() string {
	return c.identifyId
}
