package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type aboutDialog struct {
	*dialog
	internal *gtk.AboutDialog
}

func wrapAboutDialogSimple(v *gtk.AboutDialog) *aboutDialog {
	if v == nil {
		return nil
	}
	return &aboutDialog{wrapDialogSimple(&v.Dialog), v}
}

func wrapAboutDialog(v *gtk.AboutDialog, e error) (*aboutDialog, error) {
	return wrapAboutDialogSimple(v), e
}

func unwrapAboutDialog(v gtki.AboutDialog) *gtk.AboutDialog {
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
