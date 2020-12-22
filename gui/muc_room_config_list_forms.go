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

type mucRoomConfigListKickedForm struct {
	form        gtki.Box   `gtk-widget:"room-config-list-kicked-form"`
	jidEntry    gtki.Entry `gtk-widget:"room-config-list-kicked-jid"`
	reasonEntry gtki.Entry `gtk-widget:"room-config-list-kicked-reason"`

	onFormFieldValueChanges func()
}

func newMUCRoomConfigListKickedForm(onFormFieldValueChanges func()) mucRoomConfigListForm {
	kf := &mucRoomConfigListKickedForm{
		onFormFieldValueChanges: onFormFieldValueChanges,
	}

	kf.initBuilder()

	return kf
}

func (kf *mucRoomConfigListKickedForm) initBuilder() {
	builder := newBuilder("MUCRoomConfigListKickedForm")
	panicOnDevError(builder.bindObjects(kf))

	builder.ConnectSignals(map[string]interface{}{
		"on_jid_changed":    kf.onValueChanged,
		"on_reason_changed": kf.onValueChanged,
	})
}

func (kf *mucRoomConfigListKickedForm) onValueChanged() {
	if kf.onFormFieldValueChanges != nil {
		kf.onFormFieldValueChanges()
	}
}

func (kf *mucRoomConfigListKickedForm) roomConfigListForm() gtki.Widget {
	return kf.form
}

func (kf *mucRoomConfigListKickedForm) isValid() (bool, error) {
	return true, nil
}

func (kf *mucRoomConfigListKickedForm) jid() string {
	return getEntryText(kf.jidEntry)
}

func (kf *mucRoomConfigListKickedForm) reason() string {
	return getEntryText(kf.reasonEntry)
}

func (kf *mucRoomConfigListKickedForm) getValues() []string {
	return []string{kf.jid(), kf.reason()}
}

func (kf *mucRoomConfigListKickedForm) friendlyErrorMessage(err error) string {
	switch err {
	default:
		return i18n.Local("Invalid form")
	}
}

type mucRoomConfigListMembersForm struct {
	form          gtki.Box   `gtk-widget:"room-config-list-members-form"`
	jidEntry      gtki.Entry `gtk-widget:"room-config-list-members-jid"`
	nicknameEntry gtki.Entry `gtk-widget:"room-config-list-members-nickname"`
	roleEntry     gtki.Entry `gtk-widget:"room-config-list-members-role"`

	onFormFieldValueChanges func()
}

func newMUCRoomConfigListMembersForm(onFormFieldValueChanges func()) mucRoomConfigListForm {
	mf := &mucRoomConfigListMembersForm{
		onFormFieldValueChanges: onFormFieldValueChanges,
	}

	mf.initBuilder()

	return mf
}

func (mf *mucRoomConfigListMembersForm) initBuilder() {
	builder := newBuilder("MUCRoomConfigListMembersForm")
	panicOnDevError(builder.bindObjects(mf))

	builder.ConnectSignals(map[string]interface{}{
		"on_jid_changed":      mf.onValueChanged,
		"on_nickname_changed": mf.onValueChanged,
		"on_role_changed":     mf.onValueChanged,
	})
}

func (mf *mucRoomConfigListMembersForm) onValueChanged() {
	if mf.onFormFieldValueChanges != nil {
		mf.onFormFieldValueChanges()
	}
}

func (mf *mucRoomConfigListMembersForm) roomConfigListForm() gtki.Widget {
	return mf.form
}

func (mf *mucRoomConfigListMembersForm) isValid() (bool, error) {
	return true, nil
}

func (mf *mucRoomConfigListMembersForm) jid() string {
	return getEntryText(mf.jidEntry)
}

func (mf *mucRoomConfigListMembersForm) nickname() string {
	return getEntryText(mf.nicknameEntry)
}

func (mf *mucRoomConfigListMembersForm) role() string {
	return getEntryText(mf.roleEntry)
}

func (mf *mucRoomConfigListMembersForm) getValues() []string {
	return []string{mf.jid(), mf.nickname(), mf.role()}
}

func (mf *mucRoomConfigListMembersForm) friendlyErrorMessage(err error) string {
	switch err {
	default:
		return i18n.Local("Invalid form")
	}
}

type mucRoomConfigListOwnersForm struct {
	form     gtki.Box   `gtk-widget:"room-config-list-owners-form"`
	jidEntry gtki.Entry `gtk-widget:"room-config-list-owners-jid"`

	onFormFieldValueChanges func()
}

func newMUCRoomConfigListOwnersForm(onFormFieldValueChanges func()) mucRoomConfigListForm {
	of := &mucRoomConfigListOwnersForm{
		onFormFieldValueChanges: onFormFieldValueChanges,
	}

	of.initBuilder()

	return of
}

func (of *mucRoomConfigListOwnersForm) initBuilder() {
	builder := newBuilder("MUCRoomConfigListOwnersForm")
	panicOnDevError(builder.bindObjects(of))

	builder.ConnectSignals(map[string]interface{}{
		"on_jid_changed": of.onValueChanged,
	})
}

func (of *mucRoomConfigListOwnersForm) onValueChanged() {
	if of.onFormFieldValueChanges != nil {
		of.onFormFieldValueChanges()
	}
}

func (of *mucRoomConfigListOwnersForm) roomConfigListForm() gtki.Widget {
	return of.form
}

func (of *mucRoomConfigListOwnersForm) isValid() (bool, error) {
	return true, nil
}

func (of *mucRoomConfigListOwnersForm) jid() string {
	return getEntryText(of.jidEntry)
}

func (of *mucRoomConfigListOwnersForm) getValues() []string {
	return []string{of.jid()}
}

func (of *mucRoomConfigListOwnersForm) friendlyErrorMessage(err error) string {
	switch err {
	default:
		return i18n.Local("Invalid form")
	}
}

type mucRoomConfigListAdminsForm struct {
	form                    gtki.Box   `gtk-widget:"room-config-list-admins-form"`
	jidEntry                gtki.Entry `gtk-widget:"room-config-list-admins-jid"`
	onFormFieldValueChanges func()
}

func newMUCRoomConfigListAdminsForm(onFormFieldValueChanges func()) mucRoomConfigListForm {
	af := &mucRoomConfigListAdminsForm{
		onFormFieldValueChanges: onFormFieldValueChanges,
	}

	af.initBuilder()

	return af
}

func (af *mucRoomConfigListAdminsForm) initBuilder() {
	builder := newBuilder("MUCRoomConfigListAdminsForm")
	panicOnDevError(builder.bindObjects(af))

	builder.ConnectSignals(map[string]interface{}{
		"on_jid_changed": af.onValueChanged,
	})
}

func (af *mucRoomConfigListAdminsForm) onValueChanged() {
	if af.onFormFieldValueChanges != nil {
		af.onFormFieldValueChanges()
	}
}

func (af *mucRoomConfigListAdminsForm) roomConfigListForm() gtki.Widget {
	return af.form
}

func (af *mucRoomConfigListAdminsForm) isValid() (bool, error) {
	return true, nil
}

func (af *mucRoomConfigListAdminsForm) jid() string {
	return getEntryText(af.jidEntry)
}

func (af *mucRoomConfigListAdminsForm) getValues() []string {
	return []string{af.jid()}
}

func (af *mucRoomConfigListAdminsForm) friendlyErrorMessage(err error) string {
	switch err {
	default:
		return i18n.Local("Invalid form")
	}
}
