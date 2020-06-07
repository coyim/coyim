package gdka

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/gotk3/gotk3/gdk"
)

type RealGdk struct{}

var Real = &RealGdk{}

func (*RealGdk) EventButtonFrom(ev gdki.Event) gdki.EventButton {
	return WrapEventAsEventButton(eventCast(ev))
}

func (*RealGdk) EventKeyFrom(ev gdki.Event) gdki.EventKey {
	return WrapEventAsEventKey(eventCast(ev))
}

func (*RealGdk) PixbufLoaderNew() (gdki.PixbufLoader, error) {
	return WrapPixbufLoader(gdk.PixbufLoaderNew())
}

func (*RealGdk) ScreenGetDefault() (gdki.Screen, error) {
	return WrapScreen(gdk.ScreenGetDefault())
}

func (*RealGdk) WorkspaceControlSupported() bool {
	return gdk.WorkspaceControlSupported()
}
