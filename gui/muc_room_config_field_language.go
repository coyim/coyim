package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldLanguage struct {
	*roomConfigFormField
	value             *muc.RoomConfigFieldTextValue
	languageComponent *languageSelectorComponent

	languageList  gtki.ComboBoxText `gtk-widget:"room-language-combobox"`
	languageEntry gtki.Entry        `gtk-widget:"room-language-entry"`
}

func newRoomConfigFormFieldLanguage(fieldInfo roomConfigFieldTextInfo, value *muc.RoomConfigFieldTextValue) *roomConfigFormFieldLanguage {
	field := &roomConfigFormFieldLanguage{value: value}
	field.roomConfigFormField = newRoomConfigFormField(fieldInfo, "MUCRoomConfigFormFieldLanguage")

	panicOnDevError(field.builder.bindObjects(field))

	field.languageComponent = createLanguageSelectorComponent(field.languageEntry, field.languageList)
	field.languageComponent.setLanguage(value.Text())

	return field
}

// updateFieldValue MUST be called from the UI thread
func (f *roomConfigFormFieldLanguage) updateFieldValue() {
	f.value.SetText(f.languageComponent.currentLanguage())
}
