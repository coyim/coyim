package gdki

type Gdk interface {
	EventButtonFrom(Event) EventButton
	EventKeyFrom(Event) EventKey
	PixbufLoaderNew() (PixbufLoader, error)
	ScreenGetDefault() (Screen, error)
	WorkspaceControlSupported() bool
}

func AssertGdk(_ Gdk) {}
