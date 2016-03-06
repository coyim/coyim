package gdki

type Gdk interface {
	EventButtonFrom(Event) EventButton
	PixbufLoaderNew() (PixbufLoader, error)
	ScreenGetDefault() (Screen, error)
	WorkspaceControlSupported() bool
}

func AssertGdk(_ Gdk) {}
