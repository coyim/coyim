package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type dialog struct {
	*window
	internal *gtk.Dialog
}

func wrapDialogSimple(v *gtk.Dialog) *dialog {
	if v == nil {
		return nil
	}
	return &dialog{wrapWindowSimple(&v.Window), v}
}

func wrapDialog(v *gtk.Dialog, e error) (*dialog, error) {
	return wrapDialogSimple(v), e
}

func unwrapDialog(v gtki.Dialog) *gtk.Dialog {
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
