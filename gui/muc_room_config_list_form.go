package gui

import (
	"errors"

	"github.com/coyim/coyim/i18n"

	"github.com/coyim/gotk3adapter/gtki"
)

var (
	errRoomConfigListFormNotFilled = errors.New("not filled form")
)

type roomConfigListForm struct {
	formView gtki.Box   `gtk-widget:"room-config-list-form"`
	jidEntry gtki.Entry `gtk-widget:"room-config-list-jid"`

	doAfterFieldChanged  *callbacksSet // Each function will be called from the UI thread
	doAfterFieldActivate *callbacksSet // Each function will be called from the UI thread
}

func newRoomConfigListForm(onFieldChanged, onFieldActivate func()) *roomConfigListForm {
	f := &roomConfigListForm{
		doAfterFieldChanged:  newCallbacksSet(),
		doAfterFieldActivate: newCallbacksSet(),
	}

	if onFieldChanged != nil {
		f.doAfterFieldChanged.add(onFieldChanged)
	}

	if onFieldActivate != nil {
		f.doAfterFieldActivate.add(onFieldActivate)
	}

	f.initBuilder()

	return f
}

func (f *roomConfigListForm) initBuilder() {
	builder := newBuilder("MUCRoomConfigListAddForm")
	panicOnDevError(builder.bindObjects(f))

	builder.ConnectSignals(map[string]interface{}{
		"on_field_changed":  f.afterFieldChanged,
		"on_field_activate": f.afterFieldActivate,
	})
}

// afterFieldChanged MUST be called from the UI thread
func (f *roomConfigListForm) afterFieldChanged() {
	f.doAfterFieldChanged.invokeAll()
}

// afterFieldActivate MUST be called from the UI thread
func (f *roomConfigListForm) afterFieldActivate() {
	f.doAfterFieldActivate.invokeAll()
}

// jid MUST be called from the UI thread
func (f *roomConfigListForm) jid() string {
	return getEntryText(f.jidEntry)
}

// setJid MUST be called from the UI thread
func (f *roomConfigListForm) setJid(v string) {
	setEntryText(f.jidEntry, v)
}

// isFilled MUST be called from the UI thread
func (f *roomConfigListForm) isFilled() bool {
	return f.jid() != ""
}

// resetAndFocusJidEntry MUST be called from the UI thread
func (f *roomConfigListForm) resetAndFocusJidEntry() {
	f.setJid("")
	f.focusJidEntry()
}

// focusJidEntry MUST be called from the UI thread
func (f *roomConfigListForm) focusJidEntry() {
	f.jidEntry.GrabFocus()
}

func (f *roomConfigListForm) friendlyErrorMessage(err error) string {
	switch err {
	case errInvalidMemberIdentifier:
		return i18n.Local("The account address is not valid.")
	case errRoomConfigListFormNotFilled:
		return i18n.Local("Please, fill in the form fields.")
	}

	return i18n.Local("Invalid form values.")
}
