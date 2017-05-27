package core

import (
	"sync"

	"github.com/jiansoft/chopper/collections"
	"github.com/jiansoft/chopper/system"
)

type Disposer struct {
	disposables collections.ConcurrentMap
	lock        *sync.Mutex
}

func (d *Disposer) Init() *Disposer {
	d.disposables = collections.NewConcurrentMap()
	d.lock = new(sync.Mutex)
	return d
}

func NewDisposer() *Disposer {
	return new(Disposer).Init()
}

func (d *Disposer) Add(disposable system.IDisposable) {
	d.disposables.Set(disposable.IdentifyId(), disposable)

}

func (d *Disposer) Remove(disposable system.IDisposable) {
	d.disposables.Remove(disposable.IdentifyId())
}

func (d *Disposer) Count() int {
	return d.disposables.Count()
}

func (d *Disposer) Dispose() {
	d.lock.Lock()
	defer d.lock.Unlock()
	for _, key := range d.disposables.Keys() {
		if tmp, ok := d.disposables.Pop(key); ok {
			tmp.(system.IDisposable).Dispose()
		}
	}
}
