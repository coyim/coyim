package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type toolItem struct {
	*bin
	internal *gtk.ToolItem
}

func WrapToolItemSimple(v *gtk.ToolItem) gtki.ToolItem {
	if v == nil {
		return nil
	}
	return &toolItem{WrapBinSimple(&v.Bin).(*bin), v}
}

func WrapToolItem(v *gtk.ToolItem, e error) (gtki.ToolItem, error) {
	return WrapToolItemSimple(v), e
}

func UnwrapToolItem(v gtki.ToolItem) *gtk.ToolItem {
	if v == nil {
		return nil
	}
	return v.(*toolItem).internal
}

func (v *toolItem) Add(v1 gtki.Widget) {
	v.internal.Add(UnwrapWidget(v1))
}
