package cron

import (
    "fmt"
    "time"

    "github.com/jiansoft/chopper/concurrency/core"
    "github.com/jiansoft/chopper/concurrency/fiber"
    "github.com/jiansoft/chopper/math"
    "github.com/jiansoft/chopper/system"

    "log"
)

var schedulerExecutor = NewSchedulerExecutor()

type jobSchedulerExecutor struct {
    fiber *fiber.GoroutineMulti
}

func NewSchedulerExecutor() *jobSchedulerExecutor {
    return new(jobSchedulerExecutor).init()
}

func (c *jobSchedulerExecutor) init() *jobSchedulerExecutor {
    c.fiber = fiber.NewGoroutineMulti()
    c.fiber.Start()
    return c
}

func Every(interval int64) *Job {
    return schedulerExecutor.Every(interval)
}

func (c *jobSchedulerExecutor) Every(interval int64) *Job {
    return NewJob(interval, c.fiber)
}

type Job struct {
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
    nextRunTime  time.Time
    //period       int64
}

func (c *Job) init(intervel int64, fiber fiber.IFiber) *Job {
    c.hour = -1
    c.minute = -1
    c.second = -1
    c.fiber = fiber
    c.loc = time.Local
    c.interval = intervel
    c.identifyId = fmt.Sprintf("%p-%p", &c, &fiber)
    return c
}

func NewJob(intervel int64, fiber fiber.IFiber) *Job {
    return new(Job).init(intervel, fiber)
}

func (c *Job) Dispose() {
    c.taskDisposer.Dispose()
    c.fiber = nil
}

func (c Job) IdentifyId() string {
    return c.identifyId
}

func (c *Job) Days() *Job {
    c.unit = "days"
    return c
}

func (c *Job) Hours() *Job {
    c.unit = "hours"
    return c
}

func (c *Job) Minutes() *Job {
    c.unit = "minutes"
    return c
}

func (c *Job) Seconds() *Job {
    c.unit = "seconds"
    return c
}

func (c *Job) At(hour int, minute int, second int) *Job {
    c.hour = math.Abs(c.hour)
    c.minute = math.Abs(c.minute)
    c.second = math.Abs(c.second)

    if c.unit != "hours" {
        c.hour = hour % 24
    }

    c.minute = minute % 60
    c.second = second % 60
    return c
}

func (c *Job) Do(taskFun interface{}, params ...interface{}) system.IDisposable {
    c.task = core.NewTask(taskFun, params...)
    firstInMs := int64(0)
    now := time.Now()
    switch c.unit {
    case "weeks":
        i := (7 - (int(now.Weekday() - c.weekday))) % 7
        c.nextRunTime = time.Date(now.Year(), now.Month(), now.Day()+int(i), c.hour, c.minute, c.second, 0, c.loc)
        if c.nextRunTime.Before(now) {
            c.nextRunTime = c.nextRunTime.AddDate(0, 0, 7)
        }
        log.Printf("weeks nextRunTime:%v", c.nextRunTime)
        firstInMs = int64(c.nextRunTime.Sub(now)) / (1000000)
        //c.period = c.interval * 60 * 60 * 24 * 7 * 1000
        c.taskDisposer = c.fiber.Schedule(firstInMs, c.canDo)
        return c
    case "days":
        if c.second < 0 || c.minute < 0 || c.hour < 0 {
            c.nextRunTime = now.AddDate(0, 0, 1)
            c.second = c.nextRunTime.Second()
            c.minute = c.nextRunTime.Minute()
            c.hour = c.nextRunTime.Hour()
            firstInMs = int64(c.nextRunTime.Sub(now)) / (1000000)
        } else {
            c.nextRunTime = time.Date(now.Year(), now.Month(), now.Day(), c.hour, c.minute, c.second, 0, c.loc)
            c.nextRunTime = c.nextRunTime.AddDate(0, 0, int(c.interval-1))
            //if c.interval > 1 {
            // c.nextRunTime = c.nextRunTime.AddDate(0, 0, int(c.interval))
            //}
            if now.After(c.nextRunTime) {
                c.nextRunTime = c.nextRunTime.AddDate(0, 0, int(c.interval))
            }
            firstInMs = int64(c.nextRunTime.Sub(now) / 1000000)
        }
        log.Printf("days nextRunTime:%v", c.nextRunTime)
        //c.period = c.interval * 60 * 60 * 24 * 1000
        c.taskDisposer = c.fiber.Schedule(firstInMs, c.canDo)
        return c
    case "hours":
        //c.period = c.interval * 60 * 60 * 1000
        firstInMs = c.interval * 60 * 60 * 1000
        if c.minute < 0 {
            c.minute = now.Minute()
        }
        if c.second < 0 {
            c.second = now.Second()
        }
        c.nextRunTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), c.minute, c.second, 0, c.loc)
        c.nextRunTime.Add(time.Duration(c.interval-1) * time.Hour)
        if c.nextRunTime.Before(now) {
            c.nextRunTime = c.nextRunTime.Add(time.Duration(c.interval) * time.Hour /*.Duration(60*60*1000000)*/)
        }
        log.Printf("hours nextRunTime:%v", c.nextRunTime)
        firstInMs = int64(c.nextRunTime.Sub(now) / 1000000)
        //}
        c.taskDisposer = c.fiber.Schedule(firstInMs, c.canDo)
        return c
        //case "hours":
        //    c.period = c.interval * 60 * 60 * 1000
        //    firstInMs = c.period
        //    c.taskDisposer = c.fiber.ScheduleOnInterval(firstInMs, c.period, c.task.Execute)
    case "minutes":
        //c.period = c.interval * 60 * 1000
        firstInMs = c.interval * 60 * 1000
        c.taskDisposer = c.fiber.ScheduleOnInterval(firstInMs, firstInMs, c.task.Execute)
        return c
    case "seconds":
        //c.period = c.interval * 1000
        firstInMs = c.interval * 1000
        c.taskDisposer = c.fiber.ScheduleOnInterval(firstInMs, firstInMs, c.task.Execute)
        return c
    case "delay":
        c.taskDisposer = c.fiber.Schedule(c.interval, c.task.Execute)
        return c
    }
    return nil
}

func (c *Job) canDo() {
    now := time.Now()
    var adjustTime int64
    if now.After(c.nextRunTime) {
        c.fiber.EnqueueWithTask(c.task)
        switch c.unit {
        case "weeks":
            c.nextRunTime = c.nextRunTime.AddDate(0, 0, 7)
            adjustTime = int64(c.nextRunTime.Sub(now) / 1000000)
        case "days":
            c.nextRunTime = c.nextRunTime.AddDate(0, 0, int(c.interval))
            adjustTime = int64(c.nextRunTime.Sub(now) / 1000000)
        case "hours":
            c.nextRunTime = c.nextRunTime.Add(time.Hour)
            adjustTime = int64(c.nextRunTime.Sub(now) / 1000000)
        }
    } else {
        c.taskDisposer.Dispose()
        adjustTime = int64(c.nextRunTime.Sub(now)) / 1000000
        if adjustTime < 1 {
            adjustTime = 1
        }
        // c.taskDisposer = c.fiber.Schedule(adjustTime, c.canDo)
    }
    c.taskDisposer = c.fiber.Schedule(adjustTime, c.canDo)
}
