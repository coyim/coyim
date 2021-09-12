package gui

import (
	"sync"

	"github.com/coyim/coyim/config"
)

type mainConfiguration struct {
	internalConfig     *config.ApplicationConfig
	internalConfigLock sync.Mutex
	haveConfigEntries  *callbacksSet

	keySupplier config.KeySupplier
}

func (c *mainConfiguration) config() *config.ApplicationConfig {
	c.internalConfigLock.Lock()
	defer c.internalConfigLock.Unlock()

	return c.internalConfig
}

func (c *mainConfiguration) setConfig(conf *config.ApplicationConfig) {
	c.internalConfigLock.Lock()
	c.internalConfig = conf
	c.internalConfigLock.Unlock()

	c.haveConfig()
}

func (c *mainConfiguration) whenHaveConfig(f func()) {
	if c.config() != nil {
		f()
		return
	}
	c.haveConfigEntries.add(f)
}

func (c *mainConfiguration) haveConfig() {
	cs := c.haveConfigEntries
	c.haveConfigEntries = nil
	cs.invokeAll()
}
