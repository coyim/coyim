package gdka

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/twstrike/gotk3adapter/gdki"
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

func (*RealGdk) PixbufLoaderNewWithType(t string) (gdki.PixbufLoader, error) {
	return wrapPixbufLoader(gdk.PixbufLoaderNewWithType(t))
}

func (*RealGdk) PixbufGetFormats() []gdki.PixbufFormat {
	formats := gdk.PixbufGetFormats()
	if formats == nil {
		return nil
	}

	ret := make([]gdki.PixbufFormat, len(formats))
	for i, f := range formats {
		ret[i] = f
	}

	return ret
}

func (*RealGdk) ScreenGetDefault() (gdki.Screen, error) {
	return wrapScreen(gdk.ScreenGetDefault())
}

func (*RealGdk) WorkspaceControlSupported() bool {
	return gdk.WorkspaceControlSupported()
}
