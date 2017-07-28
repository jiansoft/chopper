package cron

import (
	"github.com/jiansoft/chopper/concurrency/fiber"
)

var delaySchedulerExecutor = NewJobDelaySchedulerExecutor()

type jobDelaySchedulerExecutor struct {
	fiber fiber.IFiber
}

func NewJobDelaySchedulerExecutor() *jobDelaySchedulerExecutor {
	return new(jobDelaySchedulerExecutor).init()
}

func (c *jobDelaySchedulerExecutor) init() *jobDelaySchedulerExecutor {
	c.fiber = fiber.NewGoroutineMulti()
	c.fiber.Start()
	return c
}

func (c *jobDelaySchedulerExecutor) Delay(delayInMs int64) *Job {
	return newCronDelay(delayInMs)
}

func Delay(delayInMs int64) *Job {
	return delaySchedulerExecutor.Delay(delayInMs)
}

func newCronDelay(delayInMs int64) *Job {
	c := NewJob(delayInMs, delaySchedulerExecutor.fiber)
	c.unit = "delay"
	return c
}
