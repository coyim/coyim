package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomConfigListMembersForm struct {
	*roomConfigListForm

	nicknameEntry gtki.Entry `gtk-widget:"room-config-list-members-nickname"`
	roleEntry     gtki.Entry `gtk-widget:"room-config-list-members-role"`
}

func newMUCRoomConfigListMembersForm(onFieldChanged func()) mucRoomConfigListForm {
	f := &mucRoomConfigListMembersForm{}
	f.roomConfigListForm = newRoomConfigListForm("MUCRoomConfigListFormOwners", f, onFieldChanged)
	return f
}

func (f *mucRoomConfigListMembersForm) nickname() string {
	return getEntryText(f.nicknameEntry)
}

func (f *mucRoomConfigListMembersForm) role() string {
	return getEntryText(f.roleEntry)
}

func (f *mucRoomConfigListMembersForm) getValues() []string {
	return []string{f.jid(), f.nickname(), f.role()}
}
