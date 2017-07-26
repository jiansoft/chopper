package cron

import (
	"time"

	"github.com/jiansoft/chopper/concurrency/fiber"
)

var cronWeekdayExecutor = NewCronWeekdayExecutor()

type CronWeekdayExecutor struct {
	fiber fiber.IFiber
}

func NewCronWeekdayExecutor() *CronWeekdayExecutor {
	return new(CronWeekdayExecutor).init()
}

func (c *CronWeekdayExecutor) init() *CronWeekdayExecutor {
	c.fiber = fiber.NewGoroutineMulti()
	c.fiber.Start()
	return c
}

func newCronWeekday(weekday time.Weekday) *CronScheduler {
	c := NewCron(1, cronWeekdayExecutor.fiber)
	c.unit = "weeks"
	c.weekday = weekday
	return c
}

func EverySunday() *CronScheduler {
	return newCronWeekday(time.Sunday)
}

func EveryMonday() *CronScheduler {
	return newCronWeekday(time.Monday)
}

func EveryTuesday() *CronScheduler {
	return newCronWeekday(time.Tuesday)
}

func EveryWednesday() *CronScheduler {
	return newCronWeekday(time.Wednesday)
}

func EveryThursday() *CronScheduler {
	return newCronWeekday(time.Thursday)
}

func EveryFriday() *CronScheduler {
	return newCronWeekday(time.Friday)
}

func EverySaturday() *CronScheduler {
	return newCronWeekday(time.Saturday)
}


func newCronHour() *CronScheduler {
    c := NewCron(1, cronWeekdayExecutor.fiber)
    c.unit = "hour"
    return c
}

func EveryHour() *CronScheduler {
    return newCronHour()
}
