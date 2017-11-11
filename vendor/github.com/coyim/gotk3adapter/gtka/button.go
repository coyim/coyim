package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type button struct {
	*bin
	internal *gtk.Button
}

func wrapButtonSimple(v *gtk.Button) *button {
	if v == nil {
		return nil
	}
	return &button{wrapBinSimple(&v.Bin), v}
}

func wrapButton(v *gtk.Button, e error) (*button, error) {
	return wrapButtonSimple(v), e
}

func unwrapButton(v gtki.Button) *gtk.Button {
	if v == nil {
		return nil
	}
	return v.(*button).internal
}

func (v *button) SetImage(v1 gtki.Widget) {
	v.internal.SetImage(unwrapWidget(v1))
}
