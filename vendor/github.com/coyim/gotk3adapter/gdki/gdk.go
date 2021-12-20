package gdki

type Gdk interface {
	EventButtonFrom(Event) EventButton
	EventKeyFrom(Event) EventKey
	PixbufLoaderNew() (PixbufLoader, error)
	ScreenGetDefault() (Screen, error)
	WorkspaceControlSupported() bool
	NewRGBA(...float64) Rgba
}

func AssertGdk(_ Gdk) {}
