package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type comboBoxText struct {
	*comboBox
	internal *gtk.ComboBoxText
}

func WrapComboBoxTextSimple(v *gtk.ComboBoxText) gtki.ComboBoxText {
	if v == nil {
		return nil
	}
	return &comboBoxText{WrapComboBoxSimple(&v.ComboBox).(*comboBox), v}
}

func WrapComboBoxText(v *gtk.ComboBoxText, e error) (gtki.ComboBoxText, error) {
	return WrapComboBoxTextSimple(v), e
}

func UnwrapComboBoxText(v gtki.ComboBoxText) *gtk.ComboBoxText {
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

func (v *comboBoxText) RemoveAll() {
	return v.internal.RemoveAll()
}
