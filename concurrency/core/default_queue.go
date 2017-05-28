package core

import (
	"sync"
)

type IQueue interface {
	Enqueue(t Task)
	DequeueAll() ([]Task, bool)
	Count() int
}

type defaultQueue struct {
	executor     IExecutor
	paddingTasks []Task
	toDoTasks    []Task
	lock         *sync.Mutex
	cond         *sync.Cond
}

func (d *defaultQueue) init() *defaultQueue {
	d.toDoTasks = []Task{}
	d.executor = NewDefaultExecutor()
	d.paddingTasks = []Task{}
	d.lock = new(sync.Mutex)
	d.cond = sync.NewCond(d.lock)
	return d
}

func NewDefaultQueue() *defaultQueue {
	return new(defaultQueue).init()
}

func (d *defaultQueue) Enqueue(task Task) {
	d.lock.Lock()
	d.paddingTasks = append(d.paddingTasks, task)
	d.lock.Unlock()
}

func (d *defaultQueue) DequeueAll() ([]Task, bool) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.toDoTasks, d.paddingTasks = d.paddingTasks, d.toDoTasks
	d.paddingTasks = d.paddingTasks[:0]
	return d.toDoTasks, len(d.toDoTasks) > 0
}

func (d *defaultQueue) Count() int {
	d.lock.Lock()
	defer d.lock.Unlock()
	return len(d.paddingTasks)
}
