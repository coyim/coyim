package definitions

import (
	"fmt"
	"sync"
)

var lock sync.RWMutex
var definitions = make(map[string]fmt.Stringer)

// Get returns the XML description of a UI definition and whether it was found
func Get(uiName string) (fmt.Stringer, bool) {
	lock.RLock()
	defer lock.RUnlock()

	def, ok := definitions[uiName]
	return def, ok
}

func add(uiName string, def fmt.Stringer) {
	lock.Lock()
	defer lock.Unlock()

	definitions[uiName] = def
}
