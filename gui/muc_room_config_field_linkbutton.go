package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldLinkButton struct {
	*roomConfigFormField
	pageID         int
	setCurrentPage func(indexPage int)

	linkButton gtki.LinkButton `gtk-widget:"room-config-link-button-field"`
}

func newRoomConfigFormFieldLinkButton(pageID int, setCurrentPage func(indexPage int)) hasRoomConfigFormField {
	field := &roomConfigFormFieldLinkButton{
		pageID:         pageID,
		setCurrentPage: setCurrentPage,
	}
	field.initBuilder()
	field.initDefaults()

	return field
}

func (f *roomConfigFormFieldLinkButton) initBuilder() {
	builder := newBuilder("MUCRoomConfigFormFieldLinkButton")
	panicOnDevError(builder.bindObjects(f))
	builder.ConnectSignals(map[string]interface{}{
		"go_to_page": f.goToPage,
	})
}

func (f *roomConfigFormFieldLinkButton) initDefaults() {
	f.linkButton.SetLabel(configPageDisplayTitle(f.pageID))
}

func (f *roomConfigFormFieldLinkButton) fieldWidget() gtki.Widget {
	return f.linkButton
}

func (f *roomConfigFormFieldLinkButton) goToPage() {
	f.setCurrentPage(f.pageID)
}
