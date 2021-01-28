package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type passwordConfirmationComponent struct {
	box          gtki.Box   `gtk-widget:"password-confirmation-box"`
	entry        gtki.Entry `gtk-widget:"password-entry"`
	confirmEntry gtki.Entry `gtk-widget:"password-confirmation-entry"`
}

func (u *gtkUI) createPasswordConfirmationComponent() *passwordConfirmationComponent {
	pc := &passwordConfirmationComponent{}

	builder := newBuilder("PasswordConfirmation")
	panicOnDevError(builder.bindObjects(pc))

	return pc
}

func (pc *passwordConfirmationComponent) setPassword(p string) {
	setEntryText(pc.entry, p)
}

func (pc *passwordConfirmationComponent) passwordsMatch() bool {
	return getEntryText(pc.entry) == getEntryText(pc.confirmEntry)
}

func (pc *passwordConfirmationComponent) currentPassword() string {
	return getEntryText(pc.entry)
}

func (pc *passwordConfirmationComponent) focus() {
	pc.entry.GrabFocus()
}

func (pc *passwordConfirmationComponent) focusConfirm() {
	pc.confirmEntry.GrabFocus()
}

func (pc *passwordConfirmationComponent) contentBox() gtki.Widget {
	return pc.box
}
