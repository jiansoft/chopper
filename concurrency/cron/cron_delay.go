package cron

import (
	"github.com/jiansoft/chopper/concurrency/fiber"
)

var cronDelayExecutor = NewCronDelayExecutor()

type CronDelayExecutor struct {
	fiber fiber.IFiber
}

func NewCronDelayExecutor() *CronDelayExecutor {
	return new(CronDelayExecutor).init()
}

func (c *CronDelayExecutor) init() *CronDelayExecutor {
	c.fiber = fiber.NewGoroutineMulti()
	c.fiber.Start()
	return c
}

func Delay(delayInMs int64) *CronScheduler {
	return cronDelayExecutor.Delay(delayInMs)
}

func (c *CronDelayExecutor) Delay(delayInMs int64) *CronScheduler {
	return newCronDelay(delayInMs)
}

func newCronDelay(delayInMs int64) *CronScheduler {
	c := NewCron(delayInMs, cronDelayExecutor.fiber)
	c.unit = "delay"
	return c
}
