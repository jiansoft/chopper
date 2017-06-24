package cron

import (
    "fmt"
    "time"

    "github.com/jiansoft/chopper/concurrency/core"
    "github.com/jiansoft/chopper/concurrency/fiber"
    "github.com/jiansoft/chopper/system"
)

var cronScheduler = NewCronScheduler()

type CronScheduler struct {
    fiber *fiber.GoroutineMulti
}

func (c *CronScheduler) init() *CronScheduler {
    c.fiber = fiber.NewGoroutineMulti()
    c.fiber.Start()
    return c
}

func Every(interval int64) *Cron {
    return cronScheduler.Every(interval)
}

func NewCronScheduler() *CronScheduler {
    return new(CronScheduler).init()
}

func (c *CronScheduler) Every(interval int64) *Cron {
    return NewCron(interval, c.fiber)
}

type Cron struct {
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

func (c *Cron) init(intervel int64, fiber fiber.IFiber) *Cron {
    c.hour = -1
    c.minute = -1
    c.second = -1
    c.fiber = fiber
    c.loc = time.Local
    c.interval = intervel
    c.identifyId = fmt.Sprintf("%p-%p", &c, &fiber)
    return c
}

func NewCron(intervel int64, fiber fiber.IFiber) *Cron {
    return new(Cron).init(intervel, fiber)
}

func (c *Cron) Dispose() {
    c.taskDisposer.Dispose()
    c.fiber = nil
}

func (c Cron) IdentifyId() string {
    return c.identifyId
}

func (c *Cron) Days() *Cron {
    c.unit = "days"
    return c
}

func (c *Cron) Hour() *Cron {
    c.interval = 1
    c.unit = "hour"
    return c
}

func (c *Cron) Hours() *Cron {
    c.unit = "hours"
    return c
}

func (c *Cron) Minutes() *Cron {
    c.unit = "minutes"
    return c
}

func (c *Cron) Seconds() *Cron {
    c.unit = "seconds"
    return c
}

func (c *Cron) Sunday() *Cron {
    c.interval = 1
    c.weekday = 0
    c.unit = "weeks"
    return c
}

func (c *Cron) Monday() *Cron {
    c.interval = 1
    c.weekday = 1
    c.unit = "weeks"
    return c
}

func (c *Cron) Tuesday() *Cron {
    c.interval = 1
    c.weekday = 2
    c.unit = "weeks"
    return c
}

func (c *Cron) Wednesday() *Cron {
    c.interval = 1
    c.weekday = 3
    c.unit = "weeks"
    return c
}

func (c *Cron) Thursday() *Cron {
    c.interval = 1
    c.weekday = 4
    c.unit = "weeks"
    return c
}

func (c *Cron) Friday() *Cron {
    c.interval = 1
    c.weekday = 5
    c.unit = "weeks"
    return c
}

func (c *Cron) Saturday() *Cron {
    c.interval = 1
    c.weekday = 6
    c.unit = "weeks"
    return c
}

func (c *Cron) At(hour int, minute int, second int) *Cron {
    c.hour = hour % 24
    c.minute = minute % 60
    c.second = second % 60
    return c
}

func (c *Cron) Do(taskFun interface{}, params ...interface{}) system.IDisposable {
    c.task = core.NewTask(taskFun, params...)
    return c.schedule()
}

func (c *Cron) schedule() system.IDisposable {
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
        c.taskDisposer = cronScheduler.fiber.ScheduleOnInterval(firstInMs, interval, c.task.Execute)
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
        c.taskDisposer = cronScheduler.fiber.ScheduleOnInterval(firstInMs, interval, c.task.Execute)
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
        c.taskDisposer = cronScheduler.fiber.ScheduleOnInterval(firstInMs, interval, c.task.Execute)
        return c
    case "hours":
        interval = c.interval * 60 * 60 * 1000
        firstInMs = interval
        c.taskDisposer = cronScheduler.fiber.ScheduleOnInterval(firstInMs, interval, c.task.Execute)
    case "minutes":
        interval = c.interval * 60 * 1000
        firstInMs = interval
        c.taskDisposer = cronScheduler.fiber.ScheduleOnInterval(firstInMs, interval, c.task.Execute)
        return c
    case "seconds":
        interval = c.interval * 1000
        firstInMs = interval
        c.taskDisposer = cronScheduler.fiber.ScheduleOnInterval(firstInMs, interval, c.task.Execute)
        return c
    }
    return nil
}
