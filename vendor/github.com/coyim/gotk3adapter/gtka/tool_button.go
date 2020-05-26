package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type toolButton struct {
	*bin
	internal *gtk.ToolButton
}

func wrapToolButtonSimple(v *gtk.ToolButton) *toolButton {
	if v == nil {
		return nil
	}
	return &toolButton{wrapBinSimple(&v.Bin), v}
}

func wrapToolButton(v *gtk.ToolButton, e error) (*toolButton, error) {
	return wrapToolButtonSimple(v), e
}

func unwrapToolButton(v gtki.ToolButton) *gtk.ToolButton {
	if v == nil {
		return nil
	}
	return v.(*toolButton).internal
}

func (v *toolButton) Add(v1 gtki.Widget) {
	v.internal.Add(unwrapWidget(v1))
}
