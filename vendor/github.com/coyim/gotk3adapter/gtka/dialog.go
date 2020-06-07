package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type dialog struct {
	*window
	internal *gtk.Dialog
}

func WrapDialogSimple(v *gtk.Dialog) gtki.Dialog {
	if v == nil {
		return nil
	}
	return &dialog{WrapWindowSimple(&v.Window).(*window), v}
}

func WrapDialog(v *gtk.Dialog, e error) (gtki.Dialog, error) {
	return WrapDialogSimple(v), e
}

func UnwrapDialog(v gtki.Dialog) *gtk.Dialog {
	if v == nil {
		return nil
	}
	return v.(*dialog).internal
}

func (v *dialog) Run() int {
	return int(v.internal.Run())
}

func (v *dialog) SetDefaultResponse(v1 gtki.ResponseType) {
	v.internal.SetDefaultResponse(gtk.ResponseType(v1))
}
