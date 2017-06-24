package cron

import (
	"fmt"

	"github.com/jiansoft/chopper/concurrency/core"
	"github.com/jiansoft/chopper/concurrency/fiber"
	"github.com/jiansoft/chopper/system"
)

var cronDelayExecutor = NewCronDelayExecutor()

type CronDelayExecutor struct {
	fiber fiber.IFiber
}

func NewCronDelayExecutor() *CronDelayExecutor {
	return new(CronDelayExecutor)
}

func (c *CronDelayExecutor) init() *CronDelayExecutor {
	c.fiber = fiber.NewGoroutineMulti()
	c.fiber.Start()
	return c
}

func (c *CronDelayExecutor) Delay(delayInMs int64) *CronDelay {
	return newCronDelay(delayInMs, c.fiber)
}

type CronDelay struct {
	identifyId   string
	fiber        fiber.IFiber
	task         core.Task
	taskDisposer system.IDisposable
	delayInMs    int64
}

func (c *CronDelay) init(delayInMs int64, fiber fiber.IFiber) *CronDelay {
	c.delayInMs = delayInMs
	c.fiber = fiber
	c.identifyId = fmt.Sprintf("%p-%p", &c, &fiber)
	return c
}

func newCronDelay(delayInMs int64, fiber fiber.IFiber) *CronDelay {
	return new(CronDelay).init(delayInMs, fiber)
}

func (c CronDelay) IdentifyId() string {
	return c.identifyId
}

func (c *CronDelay) Dispose() {
	c.taskDisposer.Dispose()
	c.fiber = nil
}

func Delay(delayInMs int64) *CronDelay {
	return cronDelayExecutor.Delay(delayInMs)
}

func (c *CronDelay) Do(taskFun interface{}, params ...interface{}) system.IDisposable {
	c.task = core.NewTask(taskFun, params...)
    c.taskDisposer = cronSchedulerExecutor.fiber.Schedule(c.delayInMs, c.task.Execute)
	return c
}
