package gui

import (
	"errors"

	"github.com/coyim/coyim/xmpp/jid"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

var (
	errRoomConfigListFormInvalidJid = errors.New("invalid jid")
	errRoomConfigListFormNotFilled  = errors.New("not filled form")
)

type roomConfigListForm struct {
	form                gtki.Box   `gtk-widget:"room-config-list-form"`
	jidEntry            gtki.Entry `gtk-widget:"room-config-list-jid"`
	doAfterFieldChanged []func()
}

func newRoomConfigListForm(onFieldChanged, onFieldActivate func()) *roomConfigListForm {
	f := &roomConfigListForm{}

	builder := newBuilder("MUCRoomConfigListAddForm")
	panicOnDevError(builder.bindObjects(f))

	if onFieldChanged != nil {
		f.onFieldChanged(onFieldChanged)
	}

	builder.ConnectSignals(map[string]interface{}{
		"on_field_changed": func() {
			for _, fn := range f.doAfterFieldChanged {
				fn()
			}
		},
		"on_field_activate": func() {
			if onFieldActivate != nil {
				onFieldActivate()
			}
		},
	})

	return f
}

func (f *roomConfigListForm) onFieldChanged(fn func()) {
	f.doAfterFieldChanged = append(f.doAfterFieldChanged, fn)
}

func (f *roomConfigListForm) jid() string {
	return getEntryText(f.jidEntry)
}

func (f *roomConfigListForm) setJid(v string) {
	setEntryText(f.jidEntry, v)
}

func (f *roomConfigListForm) setters() []func(string) {
	return []func(string){
		f.setJid,
	}
}

func (f *roomConfigListForm) isValid() (bool, error) {
	if f.isFilled() {
		if jid.Parse(f.jid()).Valid() {
			return true, nil
		}
		return false, errRoomConfigListFormInvalidJid
	}
	return false, errRoomConfigListFormNotFilled
}

func (f *roomConfigListForm) isFilled() bool {
	return f.jid() != ""
}

func (f *roomConfigListForm) focus() {
	f.jidEntry.GrabFocus()
}

func (f *roomConfigListForm) setValue(jid string) {
	for _, setter := range f.setters() {
		setter(jid)
	}
}

func (f *roomConfigListForm) reset() {
	for _, setter := range f.setters() {
		setter("")
	}
}

func (f *roomConfigListForm) friendlyErrorMessage(err error) string {
	return i18n.Local("Invalid form values")
}

func (f *roomConfigListForm) getFormView() gtki.Widget {
	return f.form
}
