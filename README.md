# Chopper

chopper is used for rapid development of TCP application and some useful library in Go. It is inspired by retlang.

### Features

* Fiber
* Cron this is a task scheduling package which lets your functions runs periodically at pre-determined interval and that using a human-friendly syntax.  It is inspired by [schedule](<https://github.com/dbader/schedule>).
  
this package is my first Golang program, just for fun and practice.

Usage
================

### Install

~~~
go get github.com/jiansoft/chopper
~~~

### Quick Start

```
import (
    "time"
    
    "github.com/jiansoft/chopper/concurrency/cron"
    "github.com/jiansoft/chopper/concurrency/fiber"
)

//Make a fiber
var runCronFiber = fiber.NewGoroutineMulti()
runCronFiber.Start()

//Make a cron that will be executed after 2000 ms
a := cron.Delay(2000).Do(runCron, "a Delay 2000")
 
//Make a cron that will be executed everyday at 15:30:04(HH:mm:ss)
b := cron.Every(1).Days().At(15, 30, 4).Do(runCron, "Will be cancel")

//After one second, a & b crons well be canceled by fiber 
runCronFiber.Schedule(1000, func() {
    log.Infof("a & b are dispose")
    b.Dispose()
    a.Dispose()
})

//Every friday do once at 11:50:00(HH:mm:ss).
cron.EveryFriday().At(11, 50, 0).Do(runCron, "Friday")

//Every day do once at 11:50:00(HH:mm:ss).
cron.Every(1).Days().At(11, 50, 0).Do(runCron, "Days")

//Every hour do once at N:52:20(HH:mm:ss).
cron.EveryHour().At(0, 52, 20).Do(runCron, "EveryHour")

//Every N hours do once.
cron.Every(12).Hours().Do(runCron, "Hours")

//Every N minutes do once.
cron.Every(30).Minutes().Do(runCron, "Minutes")

//Every N seconds do once.
cron.Every(100).Seconds().Do(runCron, "Seconds")

// Use a new cron executor
newCronDelay := cron.NewCronDelayExecutor()
newCronDelay.Delay(2000).Do(runCron, "newDelayCron Delay 2000 in Ms")

// Use a new cron executor
newCronScheduler := cron.NewCronSchedulerExecutor()
newCronScheduler.Every(1).Hours().Do(runCron, "newCronScheduler Hours")
newCronScheduler.Every(1).Minutes().Do(runCron, "newCronScheduler Minutes")

func runCron(s string) {
    log.Infof("I am %s CronTest %v", s, time.Now())
}

```
## License

Copyright (c) 2017

Released under the MIT license:

- [www.opensource.org/licenses/MIT](http://www.opensource.org/licenses/MIT)