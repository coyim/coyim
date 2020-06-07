package gdka

import "github.com/gotk3/gotk3/gdk"
import "github.com/coyim/gotk3adapter/gdki"

func WrapEventAsEventKey(v *event) gdki.EventKey {
	wrapped, _ := WrapEventKey(&gdk.EventKey{v.Event}, nil)
	return wrapped
}
