package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type buttonBox struct {
	*box
	internal *gtk.ButtonBox
}

type asButtonBox interface {
	toButtonBox() *buttonBox
}

func (v *buttonBox) toButtonBox() *buttonBox {
	return v
}

func WrapButtonBoxSimple(v *gtk.ButtonBox) gtki.ButtonBox {
	if v == nil {
		return nil
	}
	return &buttonBox{WrapBoxSimple(&v.Box).(*box), v}
}

func WrapButtonBox(v *gtk.ButtonBox, e error) (gtki.ButtonBox, error) {
	return WrapButtonBoxSimple(v), e
}

func UnwrapButtonBox(v gtki.ButtonBox) *gtk.ButtonBox {
	if v == nil {
		return nil
	}
	return v.(asButtonBox).toButtonBox().internal
}
