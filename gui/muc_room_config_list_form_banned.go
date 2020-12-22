package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomConfigListBannedForm struct {
	*roomConfigListForm

	reasonEntry gtki.Entry `gtk-widget:"room-config-list-banned-reason"`
}

func newMUCRoomConfigListBannedForm(onFieldChanged func()) mucRoomConfigListForm {
	f := &mucRoomConfigListBannedForm{}
	f.roomConfigListForm = newRoomConfigListForm("MUCRoomConfigListFormBanned", f, onFieldChanged)
	return f
}

func (f *mucRoomConfigListBannedForm) reason() string {
	return getEntryText(f.reasonEntry)
}

func (f *mucRoomConfigListBannedForm) getValues() []string {
	return []string{f.jid(), f.reason()}
}
