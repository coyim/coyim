package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type modelButton struct {
	*button
	internal *gtk.ModelButton
}

type asModelButton interface {
	toModelButton() *modelButton
}

func (v *modelButton) toModelButton() *modelButton {
	return v
}

func WrapModelButtonSimple(v *gtk.ModelButton) gtki.ModelButton {
	if v == nil {
		return nil
	}
	return &modelButton{WrapButtonSimple(&v.Button).(*button), v}
}

func WrapModelButton(v *gtk.ModelButton, e error) (gtki.ModelButton, error) {
	return WrapModelButtonSimple(v), e
}

func UnwrapModelButton(v gtki.ModelButton) *gtk.ModelButton {
	if v == nil {
		return nil
	}
	return v.(asModelButton).toModelButton().internal
}
