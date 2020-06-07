package gdka

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/gotk3/gotk3/gdk"
)

type screen struct {
	*gdk.Screen
}

func WrapScreenSimple(v *gdk.Screen) gdki.Screen {
	if v == nil {
		return nil
	}
	return &screen{v}
}

func WrapScreen(v *gdk.Screen, e error) (gdki.Screen, error) {
	return WrapScreenSimple(v), e
}

func UnwrapScreen(v gdki.Screen) *gdk.Screen {
	if v == nil {
		return nil
	}
	return v.(*screen).Screen
}
