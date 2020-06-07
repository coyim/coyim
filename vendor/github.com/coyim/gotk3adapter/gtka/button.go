package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type button struct {
	*bin
	internal *gtk.Button
}

func WrapButtonSimple(v *gtk.Button) gtki.Button {
	if v == nil {
		return nil
	}
	return &button{WrapBinSimple(&v.Bin).(*bin), v}
}

func WrapButton(v *gtk.Button, e error) (gtki.Button, error) {
	return WrapButtonSimple(v), e
}

func UnwrapButton(v gtki.Button) *gtk.Button {
	if v == nil {
		return nil
	}
	return v.(*button).internal
}

func (v *button) SetImage(v1 gtki.Widget) {
	v.internal.SetImage(UnwrapWidget(v1))
}
