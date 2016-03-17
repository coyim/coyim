package gdka

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/gotk3/gotk3/gdk"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gdki"
)

type RealGdk struct{}

var Real = &RealGdk{}

func (*RealGdk) EventButtonFrom(ev gdki.Event) gdki.EventButton {
	return wrapEventAsEventButton(eventCast(ev))
}

func (*RealGdk) EventKeyFrom(ev gdki.Event) gdki.EventKey {
	return wrapEventAsEventKey(eventCast(ev))
}

func (*RealGdk) PixbufLoaderNew() (gdki.PixbufLoader, error) {
	return wrapPixbufLoader(gdk.PixbufLoaderNew())
}

func (*RealGdk) ScreenGetDefault() (gdki.Screen, error) {
	return wrapScreen(gdk.ScreenGetDefault())
}

func (*RealGdk) WorkspaceControlSupported() bool {
	return gdk.WorkspaceControlSupported()
}
