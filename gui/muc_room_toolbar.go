package gui

import "github.com/coyim/gotk3adapter/gtki"

type roomViewToolbar struct {
	view                    gtki.Box    `gtk-widget:"roomToolbar"`
	roomNameLabel           gtki.Label  `gtk-widget:"roomNameLabel"`
	roomDescriptionLabel    gtki.Label  `gtk-widget:"roomDescriptionLabel"`
	togglePanelButton       gtki.Button `gtk-widget:"togglePanelButton"`
	closeConversationButton gtki.Button `gtk-widget:"closeConversationButton"`
}

func newRoomViewToolbar() *roomViewToolbar {
	t := &roomViewToolbar{}

	builder := newBuilder("MUCRoomToolbar")
	panicOnDevError(builder.bindObjects(t))

	return t
}
