package muc

import (
	"sync"

	"github.com/coyim/coyim/session/events"
)

type roomObservers struct {
	sync.RWMutex
	observers []func(events.MUC)
}

func newRoomObservers() *roomObservers {
	return &roomObservers{}
}

func (o *roomObservers) subscribe(f func(events.MUC)) {
	o.Lock()
	defer o.Unlock()

	o.observers = append(o.observers, f)
}

func (o *roomObservers) publishEvent(ev events.MUC) {
	o.RLock()
	observers := append([]func(events.MUC){}, o.observers...)
	o.RUnlock()

	for _, f := range observers {
		f(ev)
	}
}
