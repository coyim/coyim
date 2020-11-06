package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type menuButton struct {
	*bin
	internal *gtk.MenuButton
}

type asMenuButton interface {
	toMenuButton() *menuButton
}

func (v *menuButton) toMenuButton() *menuButton {
	return v
}

func WrapMenuButtonSimple(v *gtk.MenuButton) gtki.MenuButton {
	if v == nil {
		return nil
	}
	return &menuButton{WrapBinSimple(&v.Bin).(*bin), v}
}

func WrapMenuButton(v *gtk.MenuButton, e error) (gtki.MenuButton, error) {
	return WrapMenuButtonSimple(v), e
}

func UnwrapMenuButton(v gtki.MenuButton) *gtk.MenuButton {
	if v == nil {
		return nil
	}
	return v.(asMenuButton).toMenuButton().internal
}

func (*RealGtk) MenuButtonNew() (gtki.MenuButton, error) {
	return WrapMenuButton(gtk.MenuButtonNew())
}

func (v *menuButton) SetPopover(v2 gtki.Popover) {
	v.internal.SetPopover(UnwrapPopover(v2))
}
