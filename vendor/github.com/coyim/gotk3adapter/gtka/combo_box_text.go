package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type comboBoxText struct {
	*comboBox
	internal *gtk.ComboBoxText
}

func wrapComboBoxTextSimple(v *gtk.ComboBoxText) *comboBoxText {
	if v == nil {
		return nil
	}
	return &comboBoxText{wrapComboBoxSimple(&v.ComboBox), v}
}

func wrapComboBoxText(v *gtk.ComboBoxText, e error) (*comboBoxText, error) {
	return wrapComboBoxTextSimple(v), e
}

func unwrapComboBoxText(v gtki.ComboBoxText) *gtk.ComboBoxText {
	if v == nil {
		return nil
	}
	return v.(*comboBoxText).internal
}

func (v *comboBoxText) AppendText(v1 string) {
	v.internal.AppendText(v1)
}

func (v *comboBoxText) GetActiveText() string {
	return v.internal.GetActiveText()
}
