package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldLinkButton struct {
	*roomConfigFormField
	pageID   string
	goToPage func(indexPage int)

	linkButton gtki.LinkButton `gtk-widget:"room-config-link-button-field"`
}

func newRoomConfigFormFieldLinkButton(pageID string, goToPage func(indexPage int)) hasRoomConfigFormField {
	field := &roomConfigFormFieldLinkButton{
		pageID:   pageID,
		goToPage: goToPage,
	}
	field.initBuilder()
	field.initDefaults()

	return field
}

func (f *roomConfigFormFieldLinkButton) initBuilder() {
	builder := newBuilder("MUCRoomConfigFormFieldLinkButton")
	panicOnDevError(builder.bindObjects(f))
	builder.ConnectSignals(map[string]interface{}{
		"go_to_page": func() {
			f.goToPage(0)
		},
	})
}

func (f *roomConfigFormFieldLinkButton) initDefaults() {
	f.linkButton.SetLabel(configPageDisplayTitle(f.pageID))
}

func (f *roomConfigFormFieldLinkButton) fieldWidget() gtki.Widget {
	return f.linkButton
}
