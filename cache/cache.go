package cache

import (
	"sync"
	"time"
)

// New returns a new simple cache
func New() Simple {
	return NewWithExpiry()
}

// NewWithExpiry returns a new cache with expiry
func NewWithExpiry() WithExpiry {
	return &cache{data: make(map[string]*entry)}
}

type entry struct {
	key   string
	value interface{}

	cancelExpiryCheck chan<- bool
}

// expects lock to be held
func (e *entry) killPotentialExpiry() {
	if e.cancelExpiryCheck != nil {
		e.cancelExpiryCheck <- true
		e.cancelExpiryCheck = nil
	}
}

// no expectation of lock
func (c *cache) startExpiryCheck(e *entry, lifetime time.Duration, cancel <-chan bool) {
	select {
	case <-time.After(lifetime):
		c.Lock()
		defer c.Unlock()
		e.cancelExpiryCheck = nil
		delete(c.data, e.key)
	case <-cancel:
		// Do nothing, just fall through here
	}
}

// expects lock to be held
func (c *cache) setExpiry(e *entry, lifetime time.Duration) {
	cancel := make(chan bool)
	e.cancelExpiryCheck = cancel
	go c.startExpiryCheck(e, lifetime, cancel)
}

type cache struct {
	sync.RWMutex
	data map[string]*entry
}

func (c *cache) Get(key string) (result interface{}, found bool) {
	c.RLock()
	defer c.RUnlock()
	if e, ok := c.data[key]; ok {
		return e.value, true
	}
	return nil, false
}

func (c *cache) GetOrCompute(key string, creator func(key string) interface{}) (result interface{}, found bool) {
	c.Lock()
	defer c.Unlock()
	if e, ok := c.data[key]; ok {
		return e.value, true
	}
	res := creator(key)
	c.data[key] = &entry{key: key, value: res}
	return res, false
}

func (c *cache) Put(key string, value interface{}) bool {
	c.Lock()
	defer c.Unlock()
	if e, ok := c.data[key]; ok {
		e.killPotentialExpiry()
		e.value = value
		return true
	}
	c.data[key] = &entry{key: key, value: value}
	return false
}

func (c *cache) PutIfAbsent(key string, creator func(key string) interface{}) bool {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.data[key]; !ok {
		c.data[key] = &entry{key: key, value: creator(key)}
		return true
	}
	return false
}

func (c *cache) Has(key string) bool {
	c.RLock()
	defer c.RUnlock()
	_, ok := c.data[key]
	return ok
}

func (c *cache) Remove(key string) bool {
	c.Lock()
	defer c.Unlock()
	if e, ok := c.data[key]; ok {
		e.killPotentialExpiry()
		delete(c.data, key)
		return true
	}
	return false
}

func (c *cache) Clear() {
	c.Lock()
	defer c.Unlock()
	for _, e := range c.data {
		e.killPotentialExpiry()
	}
	c.data = make(map[string]*entry)
}

func (c *cache) PutTimed(key string, lifetime time.Duration, value interface{}) bool {
	c.Lock()
	defer c.Unlock()
	if e, ok := c.data[key]; ok {
		e.killPotentialExpiry()
		e.value = value
		c.setExpiry(e, lifetime)
		return true
	}
	e := &entry{key: key, value: value}
	c.setExpiry(e, lifetime)
	c.data[key] = e
	return false
}

func (c *cache) PutTimedIfAbsent(key string, lifetime time.Duration, creator func(key string) interface{}) bool {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.data[key]; !ok {
		e := &entry{key: key, value: creator(key)}
		c.setExpiry(e, lifetime)
		c.data[key] = e
		return true
	}
	return false
}

func (c *cache) GetOrComputeTimed(key string, lifetime time.Duration, creator func(key string) interface{}) (result interface{}, found bool) {
	c.Lock()
	defer c.Unlock()
	if e, ok := c.data[key]; ok {
		return e.value, true
	}
	e := &entry{key: key, value: creator(key)}
	c.setExpiry(e, lifetime)
	c.data[key] = e
	return e.value, false
}
