package gui

import (
	"errors"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomConfigListForm interface {
	jid() string
	isValid() (bool, error)
	isFilled() bool
	getValues() []string
	friendlyErrorMessage(error) string
	getFormView() gtki.Widget
}

var (
	errRoomConfigListFormInvalidJid = errors.New("invalid jid")
)

type roomConfigListForm struct {
	form     gtki.Box   `gtk-widget:"room-config-list-form"`
	jidEntry gtki.Entry `gtk-widget:"room-config-list-jid"`
}

func newRoomConfigListForm(ui string, ef interface{}, onFieldChanged, onFieldActivate func()) *roomConfigListForm {
	f := &roomConfigListForm{}

	builder := newBuilder(ui)
	panicOnDevError(builder.bindObjects(f))
	if ef != nil {
		panicOnDevError(builder.bindObjects(ef))
	}

	builder.ConnectSignals(map[string]interface{}{
		"on_field_changed": func() {
			if onFieldChanged != nil {
				onFieldChanged()
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

func (f *roomConfigListForm) jid() string {
	return getEntryText(f.jidEntry)
}

func (*roomConfigListForm) isValid() (bool, error) {
	return true, nil
}

func (f *roomConfigListForm) isFilled() bool {
	return f.jid() != ""
}

func (f *roomConfigListForm) getValues() []string {
	return []string{f.jid()}
}

func (f *roomConfigListForm) friendlyErrorMessage(err error) string {
	return i18n.Local("Invalid form values")
}

func (f *roomConfigListForm) getFormView() gtki.Widget {
	return f.form
}
