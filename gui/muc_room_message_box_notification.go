package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomMessageBoxNotification struct {
	view    gtki.InfoBar `gtk-widget:"notification-bar"`
	message gtki.Label   `gtk-widget:"message"`
}

func newRoomMessageBoxNotification() *roomMessageBoxNotification {
	v := &roomMessageBoxNotification{}

	builder := newBuilder("MUCRoomMessageBoxNotification")
	panicOnDevError(builder.bindObjects(v))

	return v
}

// widget implements widget interface
func (v *roomMessageBoxNotification) widget() gtki.Widget {
	return v.view
}

// updateMessage MUST be called from the UI thread
func (v *roomMessageBoxNotification) updateMessage(m string) {
	v.message.SetText(m)
}
