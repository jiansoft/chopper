package cron

import (
	"fmt"
	"time"

	"github.com/jiansoft/chopper/concurrency/core"
	"github.com/jiansoft/chopper/concurrency/fiber"
	"github.com/jiansoft/chopper/system"
)

var cronSchedulerExecutor = NewCronSchedulerExecutor()

type CronSchedulerExecutor struct {
	fiber *fiber.GoroutineMulti
}

func NewCronSchedulerExecutor() *CronSchedulerExecutor {
	return new(CronSchedulerExecutor).init()
}

func (c *CronSchedulerExecutor) init() *CronSchedulerExecutor {
	c.fiber = fiber.NewGoroutineMulti()
	c.fiber.Start()
	return c
}

func Every(interval int64) *CronScheduler {
	return cronSchedulerExecutor.Every(interval)
}

func (c *CronSchedulerExecutor) Every(interval int64) *CronScheduler {
	return newCron(interval, c.fiber)
}

type CronScheduler struct {
	fiber        fiber.IFiber
	identifyId   string
	loc          *time.Location
	task         core.Task
	taskDisposer system.IDisposable
	weekday      time.Weekday
	hour         int
	minute       int
	second       int
	unit         string
	interval     int64
}

func (c *CronScheduler) init(intervel int64, fiber fiber.IFiber) *CronScheduler {
	c.hour = -1
	c.minute = -1
	c.second = -1
	c.fiber = fiber
	c.loc = time.Local
	c.interval = intervel
	c.identifyId = fmt.Sprintf("%p-%p", &c, &fiber)
	return c
}

func newCron(intervel int64, fiber fiber.IFiber) *CronScheduler {
	return new(CronScheduler).init(intervel, fiber)
}

func (c *CronScheduler) Dispose() {
	c.taskDisposer.Dispose()
	c.fiber = nil
}

func (c CronScheduler) IdentifyId() string {
	return c.identifyId
}

func (c *CronScheduler) Days() *CronScheduler {
	c.unit = "days"
	return c
}

func (c *CronScheduler) Hour() *CronScheduler {
	c.interval = 1
	c.unit = "hour"
	return c
}

func (c *CronScheduler) Hours() *CronScheduler {
	c.unit = "hours"
	return c
}

func (c *CronScheduler) Minutes() *CronScheduler {
	c.unit = "minutes"
	return c
}

func (c *CronScheduler) Seconds() *CronScheduler {
	c.unit = "seconds"
	return c
}

func (c *CronScheduler) Sunday() *CronScheduler {
	c.interval = 1
	c.weekday = 0
	c.unit = "weeks"
	return c
}

func (c *CronScheduler) Monday() *CronScheduler {
	c.interval = 1
	c.weekday = 1
	c.unit = "weeks"
	return c
}

func (c *CronScheduler) Tuesday() *CronScheduler {
	c.interval = 1
	c.weekday = 2
	c.unit = "weeks"
	return c
}

func (c *CronScheduler) Wednesday() *CronScheduler {
	c.interval = 1
	c.weekday = 3
	c.unit = "weeks"
	return c
}

func (c *CronScheduler) Thursday() *CronScheduler {
	c.interval = 1
	c.weekday = 4
	c.unit = "weeks"
	return c
}

func (c *CronScheduler) Friday() *CronScheduler {
	c.interval = 1
	c.weekday = 5
	c.unit = "weeks"
	return c
}

func (c *CronScheduler) Saturday() *CronScheduler {
	c.interval = 1
	c.weekday = 6
	c.unit = "weeks"
	return c
}

func (c *CronScheduler) At(hour int, minute int, second int) *CronScheduler {
	c.hour = hour % 24
	c.minute = minute % 60
	c.second = second % 60
	return c
}

func (c *CronScheduler) Do(taskFun interface{}, params ...interface{}) system.IDisposable {
	c.task = core.NewTask(taskFun, params...)
	firstInMs := int64(0)
	interval := int64(0)
	switch c.unit {
	case "weeks":
		i := (7 - (int(time.Now().Weekday() - c.weekday))) % 7
		mock := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+int(i), c.hour, c.minute, c.second, 1000000, c.loc)
		if mock.Before(time.Now()) {
			mock = mock.AddDate(0, 0, 7)
		}
		firstInMs = int64(mock.Sub(time.Now())) / (1000 * 1000)
		interval = c.interval * 60 * 60 * 24 * 7 * 1000
		c.taskDisposer = cronSchedulerExecutor.fiber.ScheduleOnInterval(firstInMs, interval, c.task.Execute)
		return c
	case "days":
		if c.second < 0 || c.minute < 0 || c.hour < 0 {
			//c.second = time.Now().Second()
			//c.minute = time.Now().Minute()
			//c.hour = time.Now().Hour()
			firstInMs = int64(time.Now().AddDate(0, 0, 1).Sub(time.Now())) / (1000 * 1000)
		} else {
			mock := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), c.hour, c.minute, c.second, 1000000, c.loc)
			if c.interval > 1 {
				mock = mock.AddDate(0, 0, int(c.interval))
			}
			if time.Now().After(mock) {
				mock = mock.AddDate(0, 0, 1)
			}
			firstInMs = int64(mock.Sub(time.Now())) / (1000 * 1000)
		}
		interval = c.interval * 60 * 60 * 24 * 1000
		c.taskDisposer = cronSchedulerExecutor.fiber.ScheduleOnInterval(firstInMs, interval, c.task.Execute)
		return c
	case "hour":
		interval = c.interval * 60 * 60 * 1000
		firstInMs = interval
		if c.minute >= 0 {
			mock := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), c.minute, c.second, 1000000, c.loc)
			if mock.Before(time.Now()) {
				mock = mock.Add(time.Duration(60*60*1000) * time.Millisecond)
			}
			firstInMs = int64(mock.Sub(time.Now())) / (1000 * 1000)
		}
		c.taskDisposer = cronSchedulerExecutor.fiber.ScheduleOnInterval(firstInMs, interval, c.task.Execute)
		return c
	case "hours":
		interval = c.interval * 60 * 60 * 1000
		firstInMs = interval
		c.taskDisposer = cronSchedulerExecutor.fiber.ScheduleOnInterval(firstInMs, interval, c.task.Execute)
	case "minutes":
		interval = c.interval * 60 * 1000
		firstInMs = interval
		c.taskDisposer = cronSchedulerExecutor.fiber.ScheduleOnInterval(firstInMs, interval, c.task.Execute)
		return c
	case "seconds":
		interval = c.interval * 1000
		firstInMs = interval
		c.taskDisposer = cronSchedulerExecutor.fiber.ScheduleOnInterval(firstInMs, interval, c.task.Execute)
		return c
	}
	return nil
}
