package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type eventBox struct {
	*bin
	internal *gtk.EventBox
}

type asEventBox interface {
	toEventBox() *eventBox
}

func (v *eventBox) toEventBox() *eventBox {
	return v
}

func wrapEventBoxSimple(v *gtk.EventBox) *eventBox {
	if v == nil {
		return nil
	}
	return &eventBox{wrapBinSimple(&v.Bin), v}
}

func wrapEventBox(v *gtk.EventBox, e error) (*eventBox, error) {
	return wrapEventBoxSimple(v), e
}

func unwrapEventBox(v gtki.EventBox) *gtk.EventBox {
	if v == nil {
		return nil
	}
	return v.(asEventBox).toEventBox().internal
}

func (v *eventBox) SetAboveChild(v1 bool) {
	v.internal.SetAboveChild(v1)
}

func (v *eventBox) GetAboveChild() bool {
	return v.internal.GetAboveChild()
}

func (v *eventBox) SetVisibleWindow(v1 bool) {
	v.internal.SetVisibleWindow(v1)
}

func (v *eventBox) GetVisibleWindow() bool {
	return v.internal.GetVisibleWindow()
}
