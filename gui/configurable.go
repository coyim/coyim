package gui

import (
	"sync"

	"github.com/coyim/coyim/config"
)

type configurable struct {
	internalConfig     *config.ApplicationConfig
	internalConfigLock sync.Mutex
	haveConfigEntries  *callbacksSet
}

func (c *configurable) config() *config.ApplicationConfig {
	c.internalConfigLock.Lock()
	defer c.internalConfigLock.Unlock()

	return c.internalConfig
}

func (c *configurable) setConfig(conf *config.ApplicationConfig) {
	c.internalConfigLock.Lock()
	c.internalConfig = conf
	c.internalConfigLock.Unlock()

	c.haveConfig()
}

func (c *configurable) whenHaveConfig(f func()) {
	if c.config() != nil {
		f()
		return
	}
	c.haveConfigEntries.add(f)
}

func (c *configurable) haveConfig() {
	cs := c.haveConfigEntries
	c.haveConfigEntries = nil
	cs.invokeAll()
}
