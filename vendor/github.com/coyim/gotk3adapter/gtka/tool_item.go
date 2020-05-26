package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type toolItem struct {
	*bin
	internal *gtk.ToolItem
}

func wrapToolItemSimple(v *gtk.ToolItem) *toolItem {
	if v == nil {
		return nil
	}
	return &toolItem{wrapBinSimple(&v.Bin), v}
}

func wrapToolItem(v *gtk.ToolItem, e error) (*toolItem, error) {
	return wrapToolItemSimple(v), e
}

func unwrapToolItem(v gtki.ToolItem) *gtk.ToolItem {
	if v == nil {
		return nil
	}
	return v.(*toolItem).internal
}

func (v *toolItem) Add(v1 gtki.Widget) {
	v.internal.Add(unwrapWidget(v1))
}
