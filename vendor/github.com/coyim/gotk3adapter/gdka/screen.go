package gdka

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/gotk3/gotk3/gdk"
)

type screen struct {
	*gdk.Screen
}

func wrapScreenSimple(v *gdk.Screen) *screen {
	if v == nil {
		return nil
	}
	return &screen{v}
}

func wrapScreen(v *gdk.Screen, e error) (*screen, error) {
	return wrapScreenSimple(v), e
}

func UnwrapScreen(v gdki.Screen) *gdk.Screen {
	if v == nil {
		return nil
	}
	return v.(*screen).Screen
}
