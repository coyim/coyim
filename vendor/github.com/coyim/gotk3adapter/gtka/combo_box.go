package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type comboBox struct {
	*bin
	internal *gtk.ComboBox
}

func wrapComboBoxSimple(v *gtk.ComboBox) *comboBox {
	if v == nil {
		return nil
	}
	return &comboBox{wrapBinSimple(&v.Bin), v}
}

func wrapComboBox(v *gtk.ComboBox, e error) (*comboBox, error) {
	return wrapComboBoxSimple(v), e
}

func unwrapComboBox(v gtki.ComboBox) *gtk.ComboBox {
	if v == nil {
		return nil
	}
	return v.(*comboBox).internal
}

func (v *comboBox) GetActiveIter() (gtki.TreeIter, error) {
	return wrapTreeIter(v.internal.GetActiveIter())
}

func (v *comboBox) GetActiveID() string {
	return v.internal.GetActiveID()
}

func (v *comboBox) GetActive() int {
	return v.internal.GetActive()
}

func (v *comboBox) SetActive(v1 int) {
	v.internal.SetActive(v1)
}

func (v *comboBox) SetModel(v1 gtki.TreeModel) {
	v.internal.SetModel(unwrapTreeModel(v1))
}

func (v *comboBox) AddAttribute(v1 gtki.CellRenderer, v2 string, v3 int) {
	v.internal.AddAttribute(unwrapCellRenderer(v1), v2, v3)
}

func (v *comboBox) PackStart(v1 gtki.CellRenderer, v2 bool) {
	v.internal.PackStart(unwrapCellRenderer(v1), v2)
}
