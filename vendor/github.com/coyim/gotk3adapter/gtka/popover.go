package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type popover struct {
	*bin
	internal *gtk.Popover
}

type asPopover interface {
	toPopover() *popover
}

func (v *popover) toPopover() *popover {
	return v
}

func WrapPopoverSimple(v *gtk.Popover) gtki.Popover {
	if v == nil {
		return nil
	}
	return &popover{WrapBinSimple(&v.Bin).(*bin), v}
}

func WrapPopover(v *gtk.Popover, e error) (gtki.Popover, error) {
	return WrapPopoverSimple(v), e
}

func UnwrapPopover(v gtki.Popover) *gtk.Popover {
	if v == nil {
		return nil
	}
	return v.(asPopover).toPopover().internal
}

func (*RealGtk) PopoverNew(w gtki.Widget) (gtki.Popover, error) {
	return WrapPopover(gtk.PopoverNew(UnwrapWidget(w)))
}
