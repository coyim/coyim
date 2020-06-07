package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type aboutDialog struct {
	*dialog
	internal *gtk.AboutDialog
}

func WrapAboutDialogSimple(v *gtk.AboutDialog) gtki.AboutDialog {
	if v == nil {
		return nil
	}
	return &aboutDialog{WrapDialogSimple(&v.Dialog).(*dialog), v}
}

func WrapAboutDialog(v *gtk.AboutDialog, e error) (gtki.AboutDialog, error) {
	return WrapAboutDialogSimple(v), e
}

func UnwrapAboutDialog(v gtki.AboutDialog) *gtk.AboutDialog {
	if v == nil {
		return nil
	}
	return v.(*aboutDialog).internal
}

func (v *aboutDialog) SetAuthors(v1 []string) {
	v.internal.SetAuthors(v1)
}

func (v *aboutDialog) SetProgramName(v1 string) {
	v.internal.SetProgramName(v1)
}

func (v *aboutDialog) SetVersion(v1 string) {
	v.internal.SetVersion(v1)
}

func (v *aboutDialog) SetLicense(v1 string) {
	v.internal.SetLicense(v1)
}

func (v *aboutDialog) SetWrapLicense(v1 bool) {
	v.internal.SetWrapLicense(v1)
}
