package muc

import "sync"

type roomObservers struct {
	sync.RWMutex
	observers []func(MUC)
}

func newRoomObservers() *roomObservers {
	return &roomObservers{}
}

func (o *roomObservers) subscribe(f func(MUC)) {
	o.Lock()
	defer o.Unlock()

	o.observers = append(o.observers, f)
}

func (o *roomObservers) publishEvent(ev MUC) {
	o.RLock()
	observers := append([]func(MUC){}, o.observers...)
	o.RUnlock()

	for _, f := range observers {
		f(ev)
	}
}
