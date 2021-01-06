package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type passwordComponent struct {
	box          gtki.Box    `gtk-widget:"password-box"`
	entry        gtki.Entry  `gtk-widget:"password-entry"`
	toggleButton gtki.Button `gtk-widget:"password-toggle-button"`
}

func (u *gtkUI) createPasswordComponent() *passwordComponent {
	pc := &passwordComponent{}

	builder := newBuilder("Password")
	panicOnDevError(builder.bindObjects(pc))

	builder.ConnectSignals(map[string]interface{}{
		"on_password_toggle": pc.onPasswordToggle,
	})

	return pc
}

func (pc *passwordComponent) setPlaceholder(p string) {
	pc.entry.SetPlaceholderText(p)
}

func (pc *passwordComponent) setPassword(p string) {
	setEntryText(pc.entry, p)
}

func (pc *passwordComponent) currentPassword() string {
	return getEntryText(pc.entry)
}

func (pc *passwordComponent) onPasswordToggle() {
	visible := pc.entry.GetVisibility()
	pc.entry.SetVisibility(!visible)
	pc.updateToggleLabel(!visible)
}

func (pc *passwordComponent) updateToggleLabel(v bool) {
	l := i18n.Local("Show")
	if v {
		l = i18n.Local("Hide")
	}
	pc.toggleButton.SetProperty("label", l)
}

func (pc *passwordComponent) widget() gtki.Widget {
	return pc.box
}
