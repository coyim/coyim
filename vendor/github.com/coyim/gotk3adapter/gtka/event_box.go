package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
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

func WrapEventBoxSimple(v *gtk.EventBox) gtki.EventBox {
	if v == nil {
		return nil
	}
	return &eventBox{WrapBinSimple(&v.Bin).(*bin), v}
}

func WrapEventBox(v *gtk.EventBox, e error) (gtki.EventBox, error) {
	return WrapEventBoxSimple(v), e
}

func UnwrapEventBox(v gtki.EventBox) *gtk.EventBox {
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
