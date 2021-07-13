package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

const (
	configFieldListMultiValueIndex int = iota
	configFieldListMultiTextIndex
	configFieldListMultiActivableIndex
)

type roomConfigFieldListMulti struct {
	*roomConfigFormField
	value *muc.RoomConfigFieldListMultiValue
	model gtki.ListStore

	list gtki.TreeView `gtk-widget:"room-config-field-treeview"`
}

func newRoomConfigFieldListMulti(ft muc.RoomConfigFieldType, fieldInfo roomConfigFieldTextInfo, value *muc.RoomConfigFieldListMultiValue) hasRoomConfigFormField {
	field := &roomConfigFieldListMulti{value: value}
	field.roomConfigFormField = newRoomConfigFormField(ft, fieldInfo, "MUCRoomConfigFormFieldListMulti")

	field.initBuilder()
	field.initModel()

	return field
}

func (f *roomConfigFieldListMulti) initBuilder() {
	panicOnDevError(f.builder.bindObjects(f))
	f.builder.ConnectSignals(map[string]interface{}{
		"on_check_changed": func(_ gtki.CellRenderer, path string) {
			f.onCheckChanged(path)
		},
	})
}

func (f *roomConfigFieldListMulti) initModel() {
	f.model, _ = g.gtk.ListStoreNew(
		// column key
		glibi.TYPE_STRING,
		// column label
		glibi.TYPE_STRING,
		// column selected
		glibi.TYPE_BOOLEAN,
	)
	f.list.SetModel(f.model)

	for _, o := range f.value.Options() {
		iter := f.model.Append()
		f.model.SetValue(iter, configFieldListMultiValueIndex, o.Value)
		f.model.SetValue(iter, configFieldListMultiTextIndex, configOptionToFriendlyMessage(o.Value, o.Label))
		f.model.SetValue(iter, configFieldListMultiActivableIndex, f.value.IsSelected(o.Value))
	}
}

// onCheckChanged MUST be called from the UI thread
func (f *roomConfigFieldListMulti) onCheckChanged(path string) {
	iter, _ := f.model.GetIterFromString(path)
	f.model.SetValue(iter, configFieldListMultiActivableIndex, !f.isIterSelected(iter))
}

// updateFieldValue MUST be called from the UI thread
func (f *roomConfigFieldListMulti) updateFieldValue() {
	selected := []string{}

	iter, ok := f.model.GetIterFirst()
	for ok {
		if f.isIterSelected(iter) {
			selected = append(selected, f.getIterColumnKeyValue(iter))
		}
		ok = f.model.IterNext(iter)
	}

	f.value.SetSelected(selected)
}

func (f *roomConfigFieldListMulti) getIterColumnValue(iter gtki.TreeIter, columnIndex int) interface{} {
	mv, _ := f.model.GetValue(iter, columnIndex)
	gv, _ := mv.GoValue()
	return gv
}

func (f *roomConfigFieldListMulti) getIterColumnKeyValue(iter gtki.TreeIter) string {
	return f.getIterColumnValue(iter, configFieldListMultiValueIndex).(string)
}

func (f *roomConfigFieldListMulti) isIterSelected(iter gtki.TreeIter) bool {
	gv := f.getIterColumnValue(iter, configFieldListMultiActivableIndex)
	if selected, ok := gv.(bool); ok {
		return selected
	}
	return false
}
