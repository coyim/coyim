package gui

// OSHooks represents different hooks that will be called where OS specific functionality can be added
type OSHooks interface {
	// AfterInit will be called after GTK has been initialized
	AfterInit()
}

// NoHooks implements the OSHooks interface, doing nothing with the hooks
type NoHooks struct{}

// AfterInit implements the OSHooks interface
func (*NoHooks) AfterInit() {}
