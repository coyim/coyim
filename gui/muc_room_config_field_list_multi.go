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

func newRoomConfigFieldListMulti(f *muc.RoomConfigFormField, value *muc.RoomConfigFieldListMultiValue) hasRoomConfigFormField {
	field := &roomConfigFieldListMulti{value: value}
	field.roomConfigFormField = newRoomConfigFormField(f, "MUCRoomConfigFormFieldListMulti")

	panicOnDevError(field.builder.bindObjects(field))

	field.builder.ConnectSignals(map[string]interface{}{
		"on_check_changed": func(_ gtki.CellRenderer, path string) {
			field.onCheckChanged(path)
		},
	})

	field.model, _ = g.gtk.ListStoreNew(glibi.TYPE_STRING, glibi.TYPE_STRING, glibi.TYPE_BOOLEAN)
	field.list.SetModel(field.model)

	for _, o := range value.Options() {
		iter := field.model.Append()
		field.model.SetValue(iter, configFieldListMultiValueIndex, o.Value)
		field.model.SetValue(iter, configFieldListMultiTextIndex, configOptionToFriendlyMessage(o.Value, o.Label))
		field.model.SetValue(iter, configFieldListMultiActivableIndex, value.IsSelected(o))
	}

	return field

}

func (f *roomConfigFieldListMulti) onCheckChanged(path string) {
	iter, _ := f.model.GetIterFromString(path)
	mv, _ := f.model.GetValue(iter, configFieldListMultiActivableIndex)
	gv, _ := mv.GoValue()
	if active, ok := gv.(bool); ok {
		f.model.SetValue(iter, configFieldListMultiActivableIndex, !active)
	}
}

// collectFieldValue MUST be called from the UI thread
func (f *roomConfigFieldListMulti) collectFieldValue() {
	selected := []string{}

	itr, ok := f.model.GetIterFirst()
	for ok {
		mv, _ := f.model.GetValue(itr, configFieldListMultiActivableIndex)
		gv, _ := mv.GoValue()
		if active, ok := gv.(bool); ok && active {
			k, _ := f.model.GetValue(itr, configFieldListMultiValueIndex)
			kv, _ := k.GetString()
			selected = append(selected, kv)
		}
		ok = f.model.IterNext(itr)
	}
	f.value.SetSelected(selected)
}
