package gui

import (
	"errors"
	"sync"

	"github.com/coyim/coyim/i18n"

	"github.com/coyim/gotk3adapter/gtki"
)

var (
	errRoomConfigListFormNotFilled = errors.New("not filled form")
)

type roomConfigListForm struct {
	formView gtki.Box   `gtk-widget:"room-config-list-form"`
	jidEntry gtki.Entry `gtk-widget:"room-config-list-jid"`

	doAfterFieldChanged []func() // Each function will be called from the UI thread
	fieldChandedLock    sync.Mutex

	doAfterFieldActivate []func() // Each function will be called from the UI thread
	fieldActivateLock    sync.Mutex
}

func newRoomConfigListForm(onFieldChanged, onFieldActivate func()) *roomConfigListForm {
	f := &roomConfigListForm{}

	if onFieldChanged != nil {
		f.onFieldChanged(onFieldChanged)
	}

	if onFieldActivate != nil {
		f.onFieldActivate(onFieldActivate)
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
	for _, fn := range f.doAfterFieldChanged {
		fn()
	}
}

// afterFieldActivate MUST be called from the UI thread
func (f *roomConfigListForm) afterFieldActivate() {
	for _, fn := range f.doAfterFieldActivate {
		fn()
	}
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

func (f *roomConfigListForm) onFieldChanged(fn func()) {
	f.fieldChandedLock.Lock()
	defer f.fieldChandedLock.Unlock()

	f.doAfterFieldChanged = append(f.doAfterFieldChanged, fn)
}

func (f *roomConfigListForm) onFieldActivate(fn func()) {
	f.fieldActivateLock.Lock()
	defer f.fieldActivateLock.Unlock()

	f.doAfterFieldActivate = append(f.doAfterFieldActivate, fn)
}

func (f *roomConfigListForm) friendlyErrorMessage(err error) string {
	switch err {
	case errInvalidMemberIdentifier:
		return i18n.Local("The account address is not valid.")
	case errRoomConfigListFormNotFilled:
		return i18n.Local("Please, fill in the form fields.")
	default:
		return i18n.Local("Invalid form values.")
	}
}
