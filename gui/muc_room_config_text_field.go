package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormTextField struct {
	*roomConfigFormField

	entry gtki.Entry `gtk-widget:"room-config-text-field-entry"`
}

func newRoomConfigFormTextField(field *muc.RoomConfigFormField) *roomConfigFormTextField {
	f := &roomConfigFormTextField{}
	f.roomConfigFormField = newRoomConfigFormField(field, "MUCRoomConfigFormTextField")

	panicOnDevError(f.builder.bindObjects(f))

	return f
}
