package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type toolButton struct {
	*bin
	internal *gtk.ToolButton
}

func WrapToolButtonSimple(v *gtk.ToolButton) gtki.ToolButton {
	if v == nil {
		return nil
	}
	return &toolButton{WrapBinSimple(&v.Bin).(*bin), v}
}

func WrapToolButton(v *gtk.ToolButton, e error) (gtki.ToolButton, error) {
	return WrapToolButtonSimple(v), e
}

func UnwrapToolButton(v gtki.ToolButton) *gtk.ToolButton {
	if v == nil {
		return nil
	}
	return v.(*toolButton).internal
}

func (v *toolButton) Add(v1 gtki.Widget) {
	v.internal.Add(UnwrapWidget(v1))
}
