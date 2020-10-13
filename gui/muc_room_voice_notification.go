package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomVoiceNotification struct {
	view gtki.InfoBar `gtk-widget:"notification-bar"`
}

func newRoomVoiceNotification() *roomVoiceNotification {
	v := &roomVoiceNotification{}

	builder := newBuilder("MUCRoomVoiceNotification")
	panicOnDevError(builder.bindObjects(v))

	return v
}

// widget implements widget interface
func (v *roomVoiceNotification) widget() gtki.Widget {
	return v.view
}
