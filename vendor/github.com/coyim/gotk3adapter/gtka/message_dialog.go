package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type messageDialog struct {
	*dialog
	internal *gtk.MessageDialog
}

func wrapMessageDialogSimple(v *gtk.MessageDialog) *messageDialog {
	if v == nil {
		return nil
	}
	return &messageDialog{wrapDialogSimple(&v.Dialog), v}
}

func wrapMessageDialog(v *gtk.MessageDialog, e error) (*messageDialog, error) {
	return wrapMessageDialogSimple(v), e
}

func unwrapMessageDialog(v gtki.MessageDialog) *gtk.MessageDialog {
	if v == nil {
		return nil
	}
	return v.(*messageDialog).internal
}
