package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomConfigListForm interface {
	roomConfigListForm() gtki.Widget
	jid() string
	isValid() (bool, error)
	friendlyErrorMessage(error) string
	getValues() []string
}

type mucRoomConfigListMembersForm struct {
	form          gtki.Box   `gtk-widget:"room-config-list-members-form"`
	jidEntry      gtki.Entry `gtk-widget:"room-config-list-members-jid"`
	nicknameEntry gtki.Entry `gtk-widget:"room-config-list-members-nickname"`
	roleEntry     gtki.Entry `gtk-widget:"room-config-list-members-role"`

	onAnyValueChanges func()
}

func newMUCRoomConfigListMembersForm(onAnyValueChanges func()) *mucRoomConfigListMembersForm {
	f := &mucRoomConfigListMembersForm{
		onAnyValueChanges: onAnyValueChanges,
	}

	f.initBuilder()

	return f
}

func (f *mucRoomConfigListMembersForm) initBuilder() {
	builder := newBuilder("MUCRoomConfigListMembersForm")
	panicOnDevError(builder.bindObjects(f))

	builder.ConnectSignals(map[string]interface{}{
		"on_jid_changed":      f.onValueChanged,
		"on_nickname_changed": f.onValueChanged,
		"on_role_changed":     f.onValueChanged,
	})
}

func (f *mucRoomConfigListMembersForm) onValueChanged() {
	if f.onAnyValueChanges != nil {
		f.onAnyValueChanges()
	}
}

func (f *mucRoomConfigListMembersForm) roomConfigListForm() gtki.Widget {
	return f.form
}

func (f *mucRoomConfigListMembersForm) isValid() (bool, error) {
	return true, nil
}

func (f *mucRoomConfigListMembersForm) jid() string {
	return getEntryText(f.jidEntry)
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

func (f *mucRoomConfigListMembersForm) friendlyErrorMessage(err error) string {
	switch err {
	default:
		return i18n.Local("Invalid form")
	}
}

func getEntryText(e gtki.Entry) string {
	t, _ := e.GetText()
	return t
}
