package gtka

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
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
