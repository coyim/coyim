package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type messageDialog struct {
	*dialog
	internal *gtk.MessageDialog
}

func WrapMessageDialogSimple(v *gtk.MessageDialog) gtki.MessageDialog {
	if v == nil {
		return nil
	}
	return &messageDialog{WrapDialogSimple(&v.Dialog).(*dialog), v}
}

func WrapMessageDialog(v *gtk.MessageDialog, e error) (gtki.MessageDialog, error) {
	return WrapMessageDialogSimple(v), e
}

func UnwrapMessageDialog(v gtki.MessageDialog) *gtk.MessageDialog {
	if v == nil {
		return nil
	}
	return v.(*messageDialog).internal
}
