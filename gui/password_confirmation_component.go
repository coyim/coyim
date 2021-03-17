package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type passwordConfirmationComponent struct {
	box                gtki.Box    `gtk-widget:"password-confirmation-box"`
	entry              gtki.Entry  `gtk-widget:"password-entry"`
	confirmEntry       gtki.Entry  `gtk-widget:"password-confirmation-entry"`
	showPasswordButton gtki.Button `gtk-widget:"password-show-button"`
}

func (u *gtkUI) createPasswordConfirmationComponent() *passwordConfirmationComponent {
	pc := &passwordConfirmationComponent{}

	pc.initBuilder()
	pc.initDefaults()

	return pc
}

func (pc *passwordConfirmationComponent) initBuilder() {
	builder := newBuilder("PasswordConfirmation")
	panicOnDevError(builder.bindObjects(pc))
	builder.ConnectSignals(map[string]interface{}{
		"on_show_password_clicked": pc.onShowPasswordClicked,
		"on_password_change":       pc.onPasswordChange,
	})
}

func (pc *passwordConfirmationComponent) initDefaults() {
	pc.confirmEntry.SetSensitive(!pc.entry.GetVisibility())
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

// onPasswordChange MUST be called from the UI thread
func (pc *passwordConfirmationComponent) onPasswordChange() {
	if pc.entry.GetVisibility() {
		pc.confirmEntry.SetText(getEntryText(pc.entry))
	}
}

// onShowPasswordClicked MUST be called from the UI thread
func (pc *passwordConfirmationComponent) onShowPasswordClicked() {
	visible := pc.entry.GetVisibility()
	if !visible {
		pc.confirmEntry.SetText(getEntryText(pc.entry))
	}
	pc.confirmEntry.SetVisibility(!visible)
	pc.confirmEntry.SetSensitive(visible)
	pc.entry.SetVisibility(!visible)
	pc.updateShowPasswordLabel(!visible)
}

// updateShowPasswordLabel MUST be called from the UI thread
func (pc *passwordConfirmationComponent) updateShowPasswordLabel(v bool) {
	if v {
		pc.showPasswordButton.SetLabel(i18n.Local("Hide"))
		return
	}
	pc.showPasswordButton.SetLabel(i18n.Local("Show"))
}
