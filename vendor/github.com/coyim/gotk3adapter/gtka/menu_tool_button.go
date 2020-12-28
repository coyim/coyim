package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3extra"
)

type menuToolButton struct {
	*toolButton
	internal *gotk3extra.MenuToolButton
}

func WrapMenuToolButtonSimple(v *gotk3extra.MenuToolButton) gtki.MenuToolButton {
	if v == nil {
		return nil
	}
	return &menuToolButton{WrapToolButtonSimple(&v.ToolButton).(*toolButton), v}
}

func WrapMenuToolButton(v *gotk3extra.MenuToolButton, e error) (gtki.MenuToolButton, error) {
	return WrapMenuToolButtonSimple(v), e
}

func UnwrapMenuToolButton(v gtki.MenuToolButton) *gotk3extra.MenuToolButton {
	if v == nil {
		return nil
	}
	return v.(*menuToolButton).internal
}
