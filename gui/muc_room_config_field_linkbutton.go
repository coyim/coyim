package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldLinkButton struct {
	*roomConfigFormField
	pageID         mucRoomConfigPageID
	setCurrentPage func(pageID mucRoomConfigPageID)

	linkButton gtki.LinkButton `gtk-widget:"room-config-link-button-field"`
}

func newRoomConfigFormFieldLinkButton(pageID mucRoomConfigPageID, setCurrentPage func(pageID mucRoomConfigPageID)) hasRoomConfigFormField {
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

// focusWidget implements the hasRoomConfigFormField interface
func (f *roomConfigFormFieldLinkButton) focusWidget() focusable {
	return f.linkButton
}
