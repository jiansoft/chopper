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

cron.Every(1).Days().AtTime(16, 50, 0).Do(func(s string){log.Infof("I am %s now:%v", s, time.Now())}, "jIAn")
cron.Every(10).Seconds().Do(func(s string){log.Infof("I am %s now:%v", s, time.Now())}, "jIAn")

```
## License

Copyright (c) 2017

Released under the MIT license:

- [www.opensource.org/licenses/MIT](http://www.opensource.org/licenses/MIT)