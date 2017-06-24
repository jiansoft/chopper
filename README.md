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

var runCronFiber = fiber.NewGoroutineMulti()
runCronFiber.Start()

//Make a cron that will be executed everyday
b := cron.Every(1).Days().AtTime(15, 30, 4).Do(runCron, "Will be cancel")

//After one second, this cron well be canceled  
runCronFiber.Schedule(1000, func() {    
    log.Infof("Dispose()")
    b.Dispose()
})
    
cron.Every(1).Friday().At(11, 50, 0).Do(runCron, "Friday")
cron.Every(1).Days().At(11, 50, 0).Do(runCron, "Days")
cron.Every(1).Hour().At(0, 50, 0).Do(runCron, "Hour")
cron.Every(1).Hours().Do(runCron, "Hours")
cron.Every(1).Minutes().Do(runCron, "Minutes")

func runCron(s string) {
    log.Infof("I am %s CronTest %v", s, time.Now())
}

```
## License

Copyright (c) 2017

Released under the MIT license:

- [www.opensource.org/licenses/MIT](http://www.opensource.org/licenses/MIT)