package channels

import (
	"fmt"

	"github.com/jiansoft/chopper/concurrency/core"
)

type unsubscriber struct {
	identifyId    string
	channel       *channel
	receiver      core.Task
	subscriptions core.ISubscriptionRegistry //IFiber
}

func (u *unsubscriber) init(receiver core.Task, channel *channel, fiber core.ISubscriptionRegistry) *unsubscriber {
	u.identifyId = fmt.Sprintf("%p-%p", &u, &u.channel)
	u.subscriptions = fiber
	u.receiver = receiver
	u.channel = channel
	return u
}

func NewUnsubscriber(receiver core.Task, channel *channel, fiber core.ISubscriptionRegistry) *unsubscriber {
	return new(unsubscriber).init(receiver, channel, fiber)
}

func (u *unsubscriber) Dispose() {
	u.channel.Unsubscribe(u)
	//Fiber Dispose時也將訂閱的方法移除
	u.subscriptions.DeregisterSubscription(u)
}

func (u *unsubscriber) IdentifyId() string {
	return u.identifyId
}
