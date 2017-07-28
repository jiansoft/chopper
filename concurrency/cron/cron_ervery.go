package cron

import (
	"time"

	"github.com/jiansoft/chopper/concurrency/fiber"
)

var everySchedulerExecutor = NewJobEverySchedulerExecutor()

type jobEverySchedulerExecutor struct {
	fiber fiber.IFiber
}

func NewJobEverySchedulerExecutor() *jobEverySchedulerExecutor {
	return new(jobEverySchedulerExecutor).init()
}

func (c *jobEverySchedulerExecutor) init() *jobEverySchedulerExecutor {
	c.fiber = fiber.NewGoroutineMulti()
	c.fiber.Start()
	return c
}

func EverySunday() *Job {
	return newCronWeekday(time.Sunday)
}

func EveryMonday() *Job {
	return newCronWeekday(time.Monday)
}

func EveryTuesday() *Job {
	return newCronWeekday(time.Tuesday)
}

func EveryWednesday() *Job {
	return newCronWeekday(time.Wednesday)
}

func EveryThursday() *Job {
	return newCronWeekday(time.Thursday)
}

func EveryFriday() *Job {
	return newCronWeekday(time.Friday)
}

func EverySaturday() *Job {
	return newCronWeekday(time.Saturday)
}

func newCronWeekday(weekday time.Weekday) *Job {
	c := NewJob(1, everySchedulerExecutor.fiber)
	c.unit = "weeks"
	c.weekday = weekday
	return c
}

func newCronHour() *Job {
	c := NewJob(1, everySchedulerExecutor.fiber)
	c.unit = "hour"
	return c
}

func EveryHour() *Job {
	return newCronHour()
}
