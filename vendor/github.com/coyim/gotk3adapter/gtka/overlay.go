package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type overlay struct {
	*bin
	internal *gtk.Overlay
}

type asOverlay interface {
	toOverlay() *overlay
}

func (v *overlay) toOverlay() *overlay {
	return v
}

func WrapOverlaySimple(v *gtk.Overlay) gtki.Overlay {
	if v == nil {
		return nil
	}
	return &overlay{WrapBinSimple(&v.Bin).(*bin), v}
}

func WrapOverlay(v *gtk.Overlay, e error) (gtki.Overlay, error) {
	return WrapOverlaySimple(v), e
}

func UnwrapOverlay(v gtki.Overlay) *gtk.Overlay {
	if v == nil {
		return nil
	}
	return v.(asOverlay).toOverlay().internal
}
