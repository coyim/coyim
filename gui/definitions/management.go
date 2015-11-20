package definitions

import "sync"

type UIDefinition interface {
	String() string
}

var definitions = struct {
	m map[string]UIDefinition
	sync.RWMutex
}{
	m: make(map[string]UIDefinition),
}

func Get(uiName string) (UIDefinition, bool) {
	definitions.RLock()
	defer definitions.RUnlock()

	def, ok := definitions.m[uiName]
	return def, ok
}

func add(uiName string, def UIDefinition) {
	definitions.Lock()
	defer definitions.Unlock()

	definitions.m[uiName] = def
}
