package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldLanguage struct {
	form              *muc.RoomConfigForm
	languageComponent *languageSelectorComponent

	view          gtki.Box          `gtk-widget:"room-language-box"`
	languageList  gtki.ComboBoxText `gtk-widget:"room-language-combobox"`
	languageEntry gtki.Entry        `gtk-widget:"room-language-entry"`
}

func (c *mucRoomConfigComponent) newRoomConfigFormFieldLanguage() *roomConfigFormFieldLanguage {
	field := &roomConfigFormFieldLanguage{
		form: c.form,
	}

	builder := newBuilder("MUCRoomConfigFormFieldLanguage")
	panicOnDevError(builder.bindObjects(field))

	field.languageComponent = c.u.createLanguageSelectorComponent(field.languageEntry, field.languageList)
	field.languageComponent.setLanguage(field.form.Language)

	return field
}

// collectFieldValue MUST be called from the UI thread
func (f *roomConfigFormFieldLanguage) collectFieldValue() {
	f.form.Language = f.languageComponent.currentLanguage()
}

// fieldWidget implements the hasRoomConfigFormField interface
func (f *roomConfigFormFieldLanguage) fieldWidget() gtki.Widget {
	return f.view
}
