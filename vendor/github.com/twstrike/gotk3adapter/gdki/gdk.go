package gdki

type Gdk interface {
	EventButtonFrom(Event) EventButton
	EventKeyFrom(Event) EventKey
	PixbufLoaderNew() (PixbufLoader, error)
	PixbufLoaderNewWithType(string) (PixbufLoader, error)
	ScreenGetDefault() (Screen, error)
	WorkspaceControlSupported() bool
	PixbufGetFormats() []PixbufFormat
}

func AssertGdk(_ Gdk) {}
