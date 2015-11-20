package definitions

import "sync"

var definitions = struct {
	m map[string]UI
	sync.RWMutex
}{
	m: make(map[string]UI),
}

// UI represents a bundled UI description for GTK builder
type UI interface {
	String() string
}

// Get returns the XML description of a UI definition and whether it was found
func Get(uiName string) (UI, bool) {
	definitions.RLock()
	defer definitions.RUnlock()

	def, ok := definitions.m[uiName]
	return def, ok
}

func add(uiName string, def UI) {
	definitions.Lock()
	defer definitions.Unlock()

	definitions.m[uiName] = def
}
